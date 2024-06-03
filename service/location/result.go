package location

import "github.com/google/uuid"

type Location struct {
	Id       uuid.UUID
	Title    string
	Address  string
	Category string
	MapX     int
	MapY     int
}

type resQuery struct {
	Locations []Location
}
