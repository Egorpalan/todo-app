package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	targetDate, err := time.Parse(dateLayout, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %s", date)
	}

	if repeat == "" {
		return "", nil
	}

	repeatParts := strings.Fields(repeat)
	if len(repeatParts) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}

	repeatType := repeatParts[0]

	switch repeatType {
	case "y": 
		return handleYearlyRepeat(now, targetDate)
	case "d": 
		return handleDailyRepeat(now, targetDate, repeatParts)
	default:
		return "", fmt.Errorf("unsupported repeat type: %s", repeatType)
	}
}

func handleYearlyRepeat(now, date time.Time) (string, error) {
	if date.Before(now) {
		for date.Before(now) {
			date = date.AddDate(1, 0, 0)
		}
	} else {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format(dateLayout), nil
}

func handleDailyRepeat(now, date time.Time, repeatParts []string) (string, error) {
	if len(repeatParts) < 2 {
		return "", fmt.Errorf("invalid daily repeat format")
	}

	days, err := strconv.Atoi(repeatParts[1])
	if err != nil || days <= 0 || days > 400 {
		return "", fmt.Errorf("invalid number of days: %s", repeatParts[1])
	}

	nextDate := date
	if nextDate.Before(now) || nextDate.Equal(now) {
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
	} else if nextDate.After(now) {
		nextDate = nextDate.AddDate(0, 0, days)
	}

	return nextDate.Format(dateLayout), nil
}
