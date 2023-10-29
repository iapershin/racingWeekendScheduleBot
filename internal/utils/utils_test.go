package utils_test

import (
	"race-weekend-bot/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsDateInCurrentWeek(t *testing.T) {

	now := time.Now().Format("2006-01-02T15:04:05Z")
	bad := "2023-10-22" + "T" + "20:00:00Z"

	tc := map[string]bool{
		now: true,
		bad: false,
	}

	for k, v := range tc {
		b, err := utils.IsDateOnCurrentWeek(k)
		assert.NoError(t, err, "error")
		assert.Equal(t, v, b)
	}
}
