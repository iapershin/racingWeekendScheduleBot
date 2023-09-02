package messenger

import (
	"time"

	"github.com/go-co-op/gocron"
)

func (bot *Bot) RunAnnounceScheduler(announceText string, userList []int64) {
	// scheduler block
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Friday().At("18:00").Do(func() {
		bot.SendAnnounce(announceText, userList)
	})
	scheduler.StartAsync()
}
