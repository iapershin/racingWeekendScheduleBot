package datacollector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"race-weekend-bot/internal/racingapi"
	"sync"
)

var (
	ErrNoRaceThisWeekend = errors.New("no race this weekend")
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

const (
	SERIES_f1     = "f1"
	SERIES_motogp = "motogp"
)

type ApiResponseError struct {
	StatusCode int
	Error      error
}

type DataCollectorOutput struct {
	sync.Mutex
	Events map[string]racingapi.RaceWeekendSchedule
}

func NewDataCollectorOutput(l int) *DataCollectorOutput {
	return &DataCollectorOutput{
		Events: make(map[string]racingapi.RaceWeekendSchedule, l),
	}
}

func (rm *DataCollectorOutput) Store(key string, value racingapi.RaceWeekendSchedule) {
	rm.Lock()
	rm.Events[key] = value
	rm.Unlock()
}

type Series interface {
	GetData(ctx context.Context) (racingapi.RaceWeekendSchedule, error)
}

func CollectData(ctx context.Context, series []racingapi.Series) (map[string]racingapi.RaceWeekendSchedule, error) {
	var wg sync.WaitGroup
	source := "racingapi.datacollector"
	log := slog.With("collector", source)

	dataCollectorOut := NewDataCollectorOutput(len(series))

	for _, s := range series {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := s.GetData(ctx)
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
				dataCollectorOut.Store(data.RaceType, data)
			}
		}()
	}

	wg.Wait()

	if len(dataCollectorOut.Events) == 0 {
		return nil, ErrNoRaceThisWeekend
	}

	return dataCollectorOut.Events, nil
}
