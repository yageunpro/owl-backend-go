package calendar

import (
	"github.com/google/uuid"
	"time"
)

type resSchedule struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type resScheduleList struct {
	Data      []resSchedule `json:"data"`
	NextToken string        `json:"nextToken"`
}
