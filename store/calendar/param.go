package calendar

import (
	"github.com/google/uuid"
	"time"
)

type CreateScheduleParam struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Title     string
	StartTime time.Time
	EndTime   time.Time
}

type FindScheduleParam struct {
	UserId    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	Offset    int
	Limit     int
}
