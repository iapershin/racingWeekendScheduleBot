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
	f1_template = `

🏎️FORMULA 1
%s
🏁Race: %s %s
🏴Qualifying: %s %s
🏴Practice1: %s %s
🏴Practice2: %s %s
🏴Practice3: %s %s
`
	f1_template_sprint = `

🏎️FORMULA 1
%s
🏁Race: %s %s
🏁Sprint: %s %s
🏴Qualifying: %s %s
🏴Practice1: %s %s
🏴Practice2: %s %s
🏴Practice3: %s %s
`
	motogp_template = `
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
	var final_string string
	data, err := racingapi.DataCollector(ctx, s.series, log)
	if err != nil {
		return final_string, err
	}

	final_string = template

	if d, ok := data[racingapi.SERIES_f1]; ok {
		if d.Sprint.Date != "" {
			final_string += eventWithSprint(d, f1_template_sprint)
		} else {
			final_string += eventGeneral(d, f1_template)
		}
	}

	if d, ok := data[racingapi.SERIES_motogp]; ok {
		if d.Sprint.Date != "" {
			final_string += eventWithSprint(d, f1_template_sprint)
		} else {
			final_string += eventGeneral(d, f1_template)
		}
	}
	return final_string, err
}

func eventWithSprint(d racingapi.RaceWeekendSchedule, t string) string {
	return fmt.Sprintf(f1_template_sprint,
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
	return fmt.Sprintf(f1_template,
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
