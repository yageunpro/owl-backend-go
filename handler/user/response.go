package user

import "github.com/google/uuid"

type resUserInfo struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}
