package appointment

import (
	"github.com/google/uuid"
	"time"
)

type resParticipant struct {
	UserId   uuid.UUID
	Username string
}

type resLocation struct {
	Id       uuid.UUID
	Title    string
	Address  string
	Category string
	Position [2]int
}

type resInfo struct {
	Id           uuid.UUID
	OrganizerId  uuid.UUID
	Status       string
	Title        string
	Description  string
	Deadline     time.Time
	ConfirmTime  *time.Time
	Location     *resLocation
	CategoryList []string
	Participants []resParticipant
}

type absInfo struct {
	Id          uuid.UUID
	OrganizerId uuid.UUID
	Title       string
	Location    *resLocation
	Status      string
	ConfirmTime *time.Time
	HeadCount   int
}

type resInfoList struct {
	InfoList  []absInfo
	NextToken string
}
