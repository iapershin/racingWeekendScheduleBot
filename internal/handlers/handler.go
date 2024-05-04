package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"race-weekend-bot/internal/racingapi"
	"race-weekend-bot/internal/users"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	subMsg        = "You are subscribed!"
	subAlreadyMsg = "You are already subscribed"
	unsubMsg      = "You are unsubscribed"

	defaultMsg = `Hi! I'm Racing Weekend Schedule Bot
To subscribe on weekly updates type /subscribe.
To unsubscribe from updates type /unsubscribe.
Type whatever you want to see this message.`
)

type Bot interface {
	NewMessageToUser(user int64, message string) error
}

type Repository interface {
	AddUserToDB(ctx context.Context, chatID int64) error
	DeleteUserFromDB(ctx context.Context, chatID int64) error
	GetUsersList(ctx context.Context) ([]int64, error)
	CheckIfUserExists(ctx context.Context, chatID int64) error
}

type Service struct {
	bot    Bot
	repo   Repository
	series []racingapi.Series
}

func NewBotHandlers(bot Bot, repo Repository, series []racingapi.Series) *Service {
	return &Service{
		bot:    bot,
		repo:   repo,
		series: []racingapi.Series{},
	}
}

func (s Service) SubscribeHandler(ctx context.Context, update *tgbotapi.Update) error {
	source := "handlers.subscribe"
	log := slog.With("handler", source)

	log.Info(fmt.Sprintf("%d %s", update.Message.From.ID, update.Message.Text))

	user := update.Message.From.ID

	var msg string
	err := s.repo.CheckIfUserExists(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserExists):
			msg = subAlreadyMsg
		case errors.Is(err, users.ErrUserNotExists):
			msg = subMsg
			err = s.repo.AddUserToDB(ctx, user)
			if err != nil {
				log.Error(fmt.Sprintf("add user [%d] to db error: %s", update.Message.Chat.ID, err.Error()))
				return err
			}
			log.Info(fmt.Sprintf("new user [%d] subscribed", update.Message.Chat.ID))
		default:
			log.Error("check if user exists failed: %w", err)
			return err
		}
	}

	err = s.bot.NewMessageToUser(user, msg)
	if err != nil {
		log.Error("bot send message error: %w", err)
		return fmt.Errorf("can't send message to [%d]", user)
	}

	return nil
}

func (s Service) UnsubscribeHandler(ctx context.Context, update *tgbotapi.Update) error {
	source := "handlers.unsubscribe"
	log := slog.With("handler", source)

	log.Info(fmt.Sprintf("%d %s", update.Message.From.ID, update.Message.Text))

	user := update.Message.From.ID

	err := s.repo.DeleteUserFromDB(ctx, user)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("User [%d] unsubscribed", update.Message.From.ID))

	err = s.bot.NewMessageToUser(user, unsubMsg)
	if err != nil {
		log.Error("bot send message error: %w", err)
		return fmt.Errorf("can't send message to [%d]", user)
	}

	return nil
}

func (s Service) DefaultHandler(ctx context.Context, update *tgbotapi.Update) error {
	source := "handlers.unsubscribe"
	log := slog.With("handler", source)
	user := update.Message.Chat.ID
	err := s.bot.NewMessageToUser(user, defaultMsg)
	if err != nil {
		log.Error("bot send message error: %w", err)
		return fmt.Errorf("can't send message to [%d]", user)
	}
	return nil
}

func (s Service) RunAnnounceScheduler(ctx context.Context) {
	source := "handlers.announce"
	log := slog.With("handler", source)
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Friday().At("18:00").Do(func() {
		log.Info("It's time to send announce...")
		announceText, err := s.BuildAnnounceText(ctx)
		if err != nil {
			log.Error("failed to build announce text: %w", err)
		}
		userList, err := s.repo.GetUsersList(ctx)
		if err != nil {
			log.Error("failed to get user list: %w", err)
		}
		s.SendAnnounce(announceText, userList)
	})
	scheduler.StartAsync()
}

func (s Service) RunAnnounceSchedulerForce(ctx context.Context) {
	// for testing only
	source := "handlers.announce"
	log := slog.With("handler", source)
	log.Info("It's time to send announce...")
	announceText, err := s.BuildAnnounceText(ctx)
	if err != nil {
		log.Error("failed to build announce text: %w", err)
	}
	userList, err := s.repo.GetUsersList(ctx)
	if err != nil {
		log.Error("failed to get user list: %w", err)
	}
	s.SendAnnounce(announceText, userList)
}

func (s Service) SendAnnounce(announceText string, userList []int64) {
	source := "handlers.announce.sender"
	log := slog.With("handler", source)
	subsTotal := len(userList)
	affected := []int64{}

	userCh := make(chan int64, len(userList))

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	go func() {
		for _, user := range userList {
			userCh <- user
		}
		close(userCh)
	}()

	maxGR := 50 // to config??
	var errCount int32
	for i := 0; i < maxGR; i++ {
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			for user := range userCh {
				err := s.bot.NewMessageToUser(user, text)
				if err != nil {
					atomic.AddInt32(&errCount, 1)
					mu.Lock()
					affected = append(affected, user)
					mu.Unlock()
				}
			}
		}(announceText)
	}
	wg.Wait()

	rate := calcRate(errCount, subsTotal)

	log.Info(fmt.Sprintf("Sending announce is complete. Success rate is: %s", rate))
	if len(affected) > 0 {
		log.Warn(fmt.Sprintf("Affected users: %v", affected))
	}
}

func calcRate(count int32, subsTotal int) string {
	pct := (float64(count) / float64(subsTotal)) * 100
	rate := float64(100) - pct
	return fmt.Sprint(math.Round(rate)) + "%"
}
