package user

import "github.com/google/uuid"

type resInfo struct {
	Id       uuid.UUID
	Username string
	Email    string
}
