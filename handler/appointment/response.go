package appointment

import (
	"github.com/google/uuid"
	"time"
)

type resLocation struct {
	Id       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Address  string    `json:"address"`
	Category string    `json:"category"`
	Position [2]int    `json:"position"`
}

type resParticipant struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type resInfo struct {
	Id              uuid.UUID        `json:"id"`
	OrganizerId     uuid.UUID        `json:"organizer_id"`
	Title           string           `json:"title"`
	Location        *resLocation     `json:"location"`
	Status          string           `json:"status"`
	ConfirmTime     *time.Time       `json:"confirmTime"`
	Description     string           `json:"description"`
	CategoryList    []string         `json:"categoryList"`
	ParticipantList []resParticipant `json:"participantList"`
	Deadline        time.Time        `json:"deadline"`
}

type absInfo struct {
	Id          uuid.UUID    `json:"id"`
	OrganizerId uuid.UUID    `json:"organizer_id"`
	Title       string       `json:"title"`
	Location    *resLocation `json:"location"`
	Status      string       `json:"status"`
	ConfirmTime *time.Time   `json:"confirmTime"`
	HeadCount   int          `json:"headCount"`
}

type resInfoList struct {
	Data      []absInfo `json:"data"`
	NextToken string    `json:"nextToken"`
}
