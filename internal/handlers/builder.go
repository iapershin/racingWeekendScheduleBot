package handlers

import (
	"context"
	"fmt"
	"race-weekend-bot/internal/racingapi"
	"time"
)

var template = `
The racing weekend is coming! Don't miss anything
`

var (
	f1Template = `

🏎️FORMULA 1
%s
🏁Race: %s %s
🏴Qualifying: %s %s
🏴Practice1: %s %s
🏴Practice2: %s %s
🏴Practice3: %s %s
`
	f1TemplateSprint = `

🏎️FORMULA 1
%s
🏁Race: %s %s
🏁Sprint: %s %s
🏴Qualifying: %s %s
🏴Practice1: %s %s
🏴Practice2: %s %s
🏴Practice3: %s %s
`
	motogpTemplate = `
🏍️
`

	motogpTemplateSprint = `
🏍️
`
)

type Series interface {
	DataCollector(ctx context.Context) (map[string]racingapi.RaceWeekendSchedule, error)
}

func (s Service) BuildAnnounceText(ctx context.Context) (string, error) {
	source := "announce.builder.datacollector"
	log := s.log.With("handler", source)
	log.Info("build announce started")
	var finalString string
	data, err := racingapi.DataCollector(ctx, s.series, log)
	if err != nil {
		return finalString, err
	}

	finalString = template

	if d, ok := data[racingapi.SERIES_f1]; ok {
		if d.Sprint.Date != "" {
			finalString += eventWithSprint(d, f1TemplateSprint)
		} else {
			finalString += eventGeneral(d, f1Template)
		}
	}

	if d, ok := data[racingapi.SERIES_motogp]; ok {
		if d.Sprint.Date != "" {
			finalString += eventWithSprint(d, motogpTemplate)
		} else {
			finalString += eventGeneral(d, motogpTemplateSprint)
		}
	}
	return finalString, err
}

func eventWithSprint(d racingapi.RaceWeekendSchedule, t string) string {
	return fmt.Sprintf(t,
		d.RaceName,
		formatDate(d.Race.Date),
		formatTime(d.Race.Time),
		formatDate(d.Sprint.Date),
		formatTime(d.Sprint.Time),
		formatDate(d.Qualifying.Date),
		formatTime(d.Qualifying.Time),
		formatDate(d.FirstPractice.Date),
		formatTime(d.FirstPractice.Time),
		formatDate(d.SecondPractice.Date),
		formatTime(d.SecondPractice.Time),
		formatDate(d.ThirdPractice.Date),
		formatTime(d.ThirdPractice.Time),
	)
}

func eventGeneral(d racingapi.RaceWeekendSchedule, t string) string {
	return fmt.Sprintf(t,
		d.RaceName,
		formatDate(d.Race.Date),
		formatTime(d.Race.Time),
		formatDate(d.Qualifying.Date),
		formatTime(d.Qualifying.Time),
		formatDate(d.FirstPractice.Date),
		formatTime(d.FirstPractice.Time),
		formatDate(d.SecondPractice.Date),
		formatTime(d.SecondPractice.Time),
		formatDate(d.ThirdPractice.Date),
		formatTime(d.ThirdPractice.Time),
	)
}

func formatDate(str string) string {
	date, err := time.Parse("2006-01-02", str)
	if err != nil {
		return ""
	}
	return date.Weekday().String()
}

func formatTime(str string) string {
	t, err := time.Parse("15:04:05Z", str)
	if err != nil {
		return ""
	}
	return t.Format("15:04")
}
