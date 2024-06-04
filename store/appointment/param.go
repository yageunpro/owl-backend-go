package appointment

import (
	"github.com/google/uuid"
	"time"
)

type CreateAppointmentParam struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	Title        string
	Description  string
	LocationId   uuid.NullUUID
	CategoryList []string
	Deadline     time.Time
}

type UpdateAppointmentParam struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	Title        *string
	Description  *string
	LocationId   uuid.NullUUID
	CategoryList []string
}

type Schedule struct {
	Title     string
	StartTime time.Time
	EndTime   time.Time
}

type AddParticipantNonMemberParam struct {
	Id           uuid.UUID
	Username     string
	ScheduleList []Schedule
}
