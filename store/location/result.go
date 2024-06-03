package location

import (
	"github.com/google/uuid"
	"time"
)

type Location struct {
	Id       uuid.UUID
	Title    string
	Address  string
	Category string
	MapX     int
	MapY     int
}

type resGetQuery struct {
	QueryId   uuid.UUID
	QueryTime time.Time
	Locations []Location
}
