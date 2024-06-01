package auth

import "github.com/google/uuid"

type resGetOAuthUser struct {
	UserId uuid.UUID
}

type resGetDevUser struct {
	UserId       uuid.UUID
	PasswordHash string
}
