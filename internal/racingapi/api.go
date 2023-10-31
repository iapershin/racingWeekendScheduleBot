package racingapi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

var (
	ErrNoRaceThisWeekend = errors.New("no race this weekend")
)

const (
	SERIES_f1     = "f1"
	SERIES_motogp = "motogp"
)

type Event struct {
	Date string
	Time string
}

type RaceWeekendSchedule struct {
	RaceType       string
	RaceName       string
	Race           Event
	FirstPractice  Event
	SecondPractice Event
	ThirdPractice  Event
	Qualifying     Event
	Sprint         Event
}

type ApiResponseError struct {
	StatusCode int
	Error      error
}

type Series interface {
	GetData(ctx context.Context, logger Logger) (RaceWeekendSchedule, error)
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	With(args ...interface{}) *slog.Logger
}

func DataCollector(ctx context.Context, series []Series, logger Logger) (map[string]RaceWeekendSchedule, error) {
	var (
		eventsMap = make(map[string]RaceWeekendSchedule)
		mu        sync.Mutex
		wg        sync.WaitGroup
	)
	source := "racingapi.datacollector"
	log := logger.With("collector", source)

	for _, s := range series {
		wg.Add(1)
		go func(s Series) {
			defer wg.Done()
			data, err := s.GetData(ctx, logger)
			if err != nil {
				switch {
				case errors.Is(err, ErrNoRaceThisWeekend):
					log.Warn(fmt.Sprintf("%s no races found", s))
				default:
					log.Error("Error fetching data source: %w", err)
				}
				return
			}
			if data.RaceName != "" {
				mu.Lock()
				defer mu.Unlock()
				eventsMap[data.RaceType] = data
			}
		}(s)
	}

	wg.Wait()

	if len(eventsMap) == 0 {
		return eventsMap, ErrNoRaceThisWeekend
	}

	return eventsMap, nil
}
