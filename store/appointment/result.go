package appointment

import (
	"github.com/google/uuid"
	"time"
)

type resParticipant struct {
	UserId   uuid.UUID
	Username string
}

type resAppointment struct {
	Id           uuid.UUID
	OrganizerId  uuid.UUID
	Status       Status
	Title        string
	Description  string
	Deadline     time.Time
	ConfirmTime  *time.Time
	LocationId   uuid.NullUUID
	CategoryList []string
	Participants []resParticipant
}

type resAppointmentList struct {
	Id          uuid.UUID
	OrganizerId uuid.UUID
	Title       string
	LocationId  uuid.NullUUID
	Status      Status
	ConfirmTime *time.Time
	HeadCount   int
}
