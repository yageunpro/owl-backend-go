package appointment

import (
	"github.com/google/uuid"
	"time"
)

type reqAdd struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	LocationId   *uuid.UUID `json:"location_id"`
	CategoryList []string   `json:"categoryList"`
	Deadline     time.Time  `json:"deadline"`
}

type reqList struct {
	Type      string  `query:"type"`
	PageToken *string `query:"page_token"`
	Limit     *int    `query:"limit"`
}

type reqInfo struct {
	Id uuid.UUID `param:"id"`
}

type reqEdit struct {
	Id           uuid.UUID  `param:"id"`
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	LocationId   *uuid.UUID `json:"location_id"`
	CategoryList []string   `json:"categoryList"`
}

type reqDelete struct {
	Id uuid.UUID `param:"id"`
}

type reqShare struct {
	Id uuid.UUID `param:"id"`
}

type reqJoin struct {
	Id uuid.UUID `param:"id"`
}

type reqJoinNonmember struct {
	Id           uuid.UUID `param:"id"`
	Username     string    `json:"username"`
	ScheduleList []struct {
		Title     string    `json:"title"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	} `json:"scheduleList"`
}

type reqRecommendTime struct {
	Id uuid.UUID `param:"id"`
}

type reqConfirm struct {
	Id   uuid.UUID `param:"id"`
	Time time.Time `json:"time"`
}
