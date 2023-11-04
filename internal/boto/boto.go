package boto

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	BotApi     *tgbotapi.BotAPI
	UpdateChan tgbotapi.UpdatesChannel
}

type BotConfig struct {
	BotToken  string `yaml:"token"`
	Timeout   int    `yaml:"timeout" env-default:"60"`
	DebugMode bool   `yaml:"debug" env-default:"false"`
}

func NewBot(config BotConfig) (*Bot, error) {
	// bot initialization
	token, err := findTokenConf(config.BotToken)
	if err != nil {
		return nil, err
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = config.DebugMode

	u := tgbotapi.NewUpdate(0)
	u.Timeout = config.Timeout
	updates := bot.GetUpdatesChan(u)
	return &Bot{
		BotApi:     bot,
		UpdateChan: updates,
	}, nil
}

func (bot *Bot) NewMessageToUser(user int64, message string) error {
	msg := tgbotapi.NewMessage(user, message)
	_, err := bot.BotApi.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %s", err.Error())
	}
	return nil
}

func findTokenConf(token string) (string, error) {
	if token == "" {
		if t := os.Getenv("BOT_TOKEN"); t != "" {
			return t, nil
		} else {
			return "", fmt.Errorf("no bot token provided")
		}
	} else {
		return token, nil
	}
}
