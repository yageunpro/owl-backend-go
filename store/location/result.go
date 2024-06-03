package location

import (
	"github.com/google/uuid"
	"time"
)

type resGetQuery struct {
	QueryId   uuid.UUID
	QueryTime time.Time
	Locations []Location
}
