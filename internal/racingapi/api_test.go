package racingapi_test

import (
	"context"
	mock "race-weekend-bot/internal/racingapi/_mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF1ApiResponse(t *testing.T) {
	ctx := context.Background()
	api := mock.MockAPI{}
	_, err := api.GetData(ctx, nil)

	assert.NoError(t, err)
}
