package appointment

import (
	"github.com/google/uuid"
	"time"
)

type AddParam struct {
	UserId       uuid.UUID
	Title        string
	Description  string
	LocationId   uuid.NullUUID
	CategoryList []string
	Deadline     time.Time
}

type ListParam struct {
	UserId    uuid.UUID
	Status    string
	PageToken *string
	Limit     int
}

type EditParam struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	Title        *string
	Description  *string
	LocationId   *uuid.UUID
	CategoryList []string
}

type JoinSchedule struct {
	Title     string
	StartTime time.Time
	EndTime   time.Time
}

type JoinNonmemberParam struct {
	Id           uuid.UUID
	Username     string
	JoinSchedule []JoinSchedule
}
