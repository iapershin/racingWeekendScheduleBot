package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"race-weekend-bot/internal/boto"
	"race-weekend-bot/internal/config"
	"race-weekend-bot/internal/handlers"
	"race-weekend-bot/internal/logger"
	"race-weekend-bot/internal/racingapi"
	"race-weekend-bot/internal/racingapi/f1"
	"race-weekend-bot/internal/racingapi/motogp"
	"race-weekend-bot/internal/storage/postgres"
	"race-weekend-bot/internal/users"
)

type Flags struct {
	Config string
}

func main() {

	// parse flags
	flags := Flags{}
	parseFlags(&flags)

	// load config
	cfg := config.MustLoad(flags.Config)

	//load logger
	if err := logger.Init(logger.HandlerOptions{
		Level:  cfg.App.Logger.Level,
		Format: logger.Format(cfg.App.Logger.Format),
	}); err != nil {
		log.Fatalf("unable to initialize logger: %s", err.Error())
	}

	// load logger
	slog.Info("race weekend schedule bot starting...")

	// set parent context
	ctx := context.Background()

	// init bot
	bot, err := boto.NewBot(boto.BotConfig{
		BotToken:  cfg.Bot.BotToken,
		Timeout:   cfg.Bot.Timeout,
		DebugMode: cfg.Bot.DebugMode,
	})
	if err != nil {
		log.Fatal("can't init bot: %w", err)
	}

	slog.Info(fmt.Sprintf("authorized on account: %s", bot.BotApi.Self.UserName))

	// init database
	db, err := postgres.New(ctx, postgres.Config{
		Host:     cfg.Postgress.Host,
		Port:     cfg.Postgress.Port,
		User:     cfg.Postgress.User,
		Password: cfg.Postgress.Password,
		Database: cfg.Postgress.Database,
	})
	if err != nil {
		log.Fatal("unable to connect to postgres: %w", err)
	}

	// init handlers
	series := []racingapi.Series{
		f1.F1API{URL: cfg.Api.F1},
		motogp.MotoGPApi{URL: cfg.Api.Motogp},
	}
	botHandlers := handlers.NewBotHandlers(bot, users.NewPostgresRepository(db), series)

	// schedule handler
	botHandlers.RunAnnounceScheduler(ctx)

	// chat handler
	for update := range bot.UpdateChan {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		switch update.Message.Text {
		case "/subscribe":
			botHandlers.SubscribeHandler(ctx, &update)
		case "/unsubscribe":
			botHandlers.UnsubscribeHandler(ctx, &update)
		default:
			botHandlers.DefaultHandler(ctx, &update)
		}
	}
}

func parseFlags(f *Flags) {
	flag.StringVar(&f.Config, "config", "defaultValue", "path to config file")
	flag.Parse()
}
