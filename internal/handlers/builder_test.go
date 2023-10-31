package handlers

import (
	"context"
	"errors"
	"race-weekend-bot/internal/logger"
	"race-weekend-bot/internal/racingapi"
	"race-weekend-bot/internal/racingapi/f1"
	"race-weekend-bot/internal/racingapi/motogp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF1ApiResponse(t *testing.T) {
	ctx := context.Background()

	series := []racingapi.Series{
		f1.F1API{URL: "https://ergast.com/api/f1/current/next.json"},
		motogp.MotoGPApi{URL: "https://to.do"},
	}

	service := Service{
		series: series,
		log:    logger.NewLogger("test"),
	}

	s, err := service.BuildAnnounceText(ctx)

	if errors.Is(err, racingapi.ErrNoRaceThisWeekend) {
		assert.Empty(t, s)
	}
}
