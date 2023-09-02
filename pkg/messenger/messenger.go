package messenger

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/syfaro/telegram-bot-api"
)

type Bot struct {
	BotApi     *tgbotapi.BotAPI
	UpdateChan tgbotapi.UpdatesChannel
}

func NewBot() (Bot, error) {
	// bot initialization
	if os.Getenv("BOT_TOKEN") == "" {
		log.Panic("CRITICAL No bot token provided")
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Println(err)
		return Bot{}, fmt.Errorf("CRITICAL Can't initialize bot")
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	return Bot{
		BotApi:     bot,
		UpdateChan: updates,
	}, nil
}

func (bot *Bot) SendAnnounce(announceText string, userList []int64) error {
	for _, user := range userList {
		msg := tgbotapi.NewMessage(user, announceText)
		_, err := bot.BotApi.Send(msg)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("ERROR Sending message: %s", err.Error())
		}
		log.Println("New announce sent:")
		log.Println(announceText)
	}
	return nil
}
