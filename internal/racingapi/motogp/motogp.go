package motogp

import (
	"context"
	"race-weekend-bot/internal/racingapi"
)

type MotoGPApi struct {
	URL string
}

func (a MotoGPApi) GetData(ctx context.Context, logger racingapi.Logger) (racingapi.RaceWeekendSchedule, error) {
	// TODO: to implemet
	return racingapi.RaceWeekendSchedule{}, racingapi.NoRaceThisWeekend
}
