package auth

import "github.com/google/uuid"

const (
	CookieKey = "REF"
)

type refCookieValue struct {
	State uuid.UUID `json:"state"`
	Ref   string    `json:"ref"`
}
