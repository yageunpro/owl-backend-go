package auth

import (
	"github.com/google/uuid"
	"time"
)

type CreateDevUserParam struct {
	UserId       uuid.UUID
	Email        string
	PasswordHash string
}

type CreateOAuthUserParam struct {
	UserId       uuid.UUID
	Email        string
	UserName     string
	OpenId       string
	AccessToken  string
	RefreshToken *string
	AllowSync    bool
	ValidUntil   time.Time
}
