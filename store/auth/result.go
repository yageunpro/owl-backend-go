package auth

import (
	"github.com/google/uuid"
	"time"
)

type resGetOAuthUser struct {
	UserId uuid.UUID
}

type resGetOAuthToken struct {
	UserId       uuid.UUID
	OpenId       string
	AccessToken  string
	RefreshToken string
	ExpireTime   time.Time
}

type resGetDevUser struct {
	UserId       uuid.UUID
	PasswordHash string
}
