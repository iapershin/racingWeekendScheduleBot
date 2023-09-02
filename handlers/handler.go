package handlers

import (
	"fmt"
	"log"
	"main/pkg/db"

	tgbotapi "github.com/syfaro/telegram-bot-api"
)

const (
	SUB_MSG         = "You are subscribed!"
	SUB_ALREADY_MSG = "You are already subscribed"
	UNSUB_MSG       = "You are unsubscribed"

	DEFAULT_MSG = `Hi! I'm Racing Weekend Schedule Bot
	To subscribe on weekly updates type /subscribe.
	To unsubscribe from updates type /unsubscribe.
	Type whatever you want to see this message.`
)

func SubscribeHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	log.Printf("[%s] %s", update.Message.ForwardFrom.UserName, update.Message.Text)

	userExists, err := db.CheckIfUserExists(update.Message.Chat.ID)
	if err != nil {
		return err
	}

	var msg tgbotapi.MessageConfig

	if userExists {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, SUB_ALREADY_MSG)
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, SUB_MSG)
		err = db.AddUserToDB(update.Message.Chat.ID)
		if err != nil {
			return err
		}
		log.Printf("New user %d subscribed", update.Message.Chat.ID)
	}

	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("ERROR Can't send message")
	}

	return nil
}

func UnsubscribeHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	err := db.DeleteUserFromDB(update.Message.Chat.ID)
	if err != nil {
		return err
	}
	log.Printf("User %d unsubscribed", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, UNSUB_MSG)
	_, err = bot.Send(msg)

	if err != nil {
		return fmt.Errorf("ERROR Can't send message")
	}
	return nil
}

func DefaultHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, DEFAULT_MSG)
	_, err := bot.Send(msg)
	if err != nil {
		return fmt.Errorf("ERROR Can't send message")
	}
	return nil
}
