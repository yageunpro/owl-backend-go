package appointment

import (
	"time"
)

type TimeSlot struct {
	StartTime       time.Time
	EndTime         time.Time
	OverlapCount    int
	PreferenceScore int
}

type TimePeriod struct {
	StartTime time.Time
	EndTime   time.Time
}

func doTimePeriodsOverlap(period1, period2 TimePeriod) bool {
	if !period1.StartTime.Before(period2.EndTime) || !period2.StartTime.Before(period1.EndTime) {
		return false
	}

	if !period1.StartTime.After(period2.StartTime) && !period1.EndTime.Before(period2.EndTime) {
		return true
	}

	if !period2.StartTime.After(period1.StartTime) && !period2.EndTime.Before(period1.EndTime) {
		return true
	}

	return true
}
