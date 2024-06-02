package calendar

import (
	"github.com/google/uuid"
	"time"
)

type ScheduleAddParam struct {
	UserId    uuid.UUID
	Title     string
	StartTime time.Time
	EndTime   time.Time
}

type ScheduleListParam struct {
	UserId    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	PageToken *string
	Limit     int
}
