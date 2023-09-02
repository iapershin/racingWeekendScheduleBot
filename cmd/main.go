package main

import (
	"main/handlers"
	"main/pkg/messenger"

	_ "github.com/lib/pq"
)

func main() {

	bot, err := messenger.NewBot()
	if err != nil {
		panic(err)
	}

	// chat handler
	for update := range bot.UpdateChan {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		switch update.Message.Text {
		case "/subscribe":
			handlers.SubscribeHandler(bot.BotApi, &update)
		case "/unsubscribe":
			handlers.UnsubscribeHandler(bot.BotApi, &update)
		default:
			handlers.DefaultHandler(bot.BotApi, &update)
		}
	}
}
