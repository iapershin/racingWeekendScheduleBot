package racingapi_test

import (
	"context"
	"fmt"
	"race-weekend-bot/internal/logger"
	"race-weekend-bot/internal/racingapi/f1"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF1ApiResponse(t *testing.T) {
	ctx := context.Background()
	api := f1.F1API{URL: "https://ergast.com/api/f1/current/next.json"}
	data, err := api.GetData(ctx, logger.NewLogger("test"))
	assert.NoError(t, err, "error")
	fmt.Println(data)
}
