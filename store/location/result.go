package location

import "github.com/google/uuid"

type resGetQuery struct {
	QueryId   uuid.UUID
	Locations []Location
}
