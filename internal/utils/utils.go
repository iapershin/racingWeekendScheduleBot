package utils

import "time"

func IsDateOnCurrentWeek(dateStr string) (bool, error) {
	date, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		return false, err
	}
	currentWeekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	nextWeekStart := currentWeekStart.AddDate(0, 0, 7)
	return date.After(currentWeekStart) && date.Before(nextWeekStart), nil
}
