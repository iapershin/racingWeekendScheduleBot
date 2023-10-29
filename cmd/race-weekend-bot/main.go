package main

import (
	"context"
	"flag"
	"os"
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

	//parse flags
	flags := Flags{}
	parseFlags(&flags)

	//load config
	cfg := config.MustLoad(flags.Config)

	//load logger
	log := logger.NewLogger(cfg.Env)
	log.Info("race weekend schedule bot starting...")

	//set parent context
	ctx := context.Background()

	//init bot
	bot, err := boto.NewBot(boto.BotConfig{
		BotToken:  cfg.Bot.BotToken,
		Timeout:   cfg.Bot.Timeout,
		DebugMode: cfg.Bot.DebugMode,
	})
	if err != nil {
		log.Error("can't init bot: %w", err)
		os.Exit(1)
	}

	//init database
	db, err := postgres.New(ctx, postgres.Config{
		Host:     cfg.Postgress.Host,
		Port:     cfg.Postgress.Port,
		User:     cfg.Postgress.User,
		Password: cfg.Postgress.Password,
		Database: cfg.Postgress.Database,
	})
	if err != nil {
		log.Error("unable to connect to postgres: %v", err)
		os.Exit(1)
	}
	//define series
	series := []racingapi.Series{
		f1.F1API{URL: cfg.Api.F1},
		motogp.MotoGPApi{URL: cfg.Api.Motogp},
	}

	// init handlers
	botHandlers := handlers.NewBotHandlers(bot,
		users.NewUserRepository(db),
		log,
		series)

	//schedule handler
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
