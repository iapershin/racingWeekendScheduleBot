package motogp

import (
	"context"
	"race-weekend-bot/internal/racingapi"
)

type MotoGPApi struct {
	URL string
}

func (a MotoGPApi) GetData(ctx context.Context) (racingapi.RaceWeekendSchedule, error) {
	// TODO: to implemet
	return racingapi.RaceWeekendSchedule{}, racingapi.ErrNoRaceThisWeekend
}
