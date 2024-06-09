package calendar

import (
	"github.com/google/uuid"
	"time"
)

type resSchedule struct {
	Id        uuid.UUID
	Title     string
	StartTime time.Time
	EndTime   time.Time
}

type resScheduleTime struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}
