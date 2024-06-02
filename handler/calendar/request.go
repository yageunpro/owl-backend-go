package calendar

import (
	"github.com/google/uuid"
	"time"
)

type reqScheduleAdd struct {
	Title     string    `json:"title"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type reqScheduleInfo struct {
	Id uuid.UUID `param:"id"`
}

type reqScheduleDelete struct {
	Id uuid.UUID `param:"id"`
}

type reqScheduleList struct {
	StartTime time.Time `query:"start"`
	EndTime   time.Time `query:"end"`
	PageToken *string   `query:"page_token"`
	Limit     int       `query:"limit"`
}
