package location

import "github.com/google/uuid"

type resLocation struct {
	Id       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Address  string    `json:"address"`
	Category string    `json:"category"`
	Position [2]int    `json:"position"`
}
