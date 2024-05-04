package racingapi

import (
	"context"
	"errors"
)

var (
	ErrNoRaceThisWeekend = errors.New("no race this weekend")
)

const (
	SERIES_f1     = "f1"
	SERIES_motogp = "motogp"
)

type Series interface {
	GetData(ctx context.Context) (RaceWeekendSchedule, error)
}

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
