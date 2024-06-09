package auth

import (
	"github.com/google/uuid"
	"net/http"
)

type resGoogleLogin struct {
	RedirectURL string
	Cookie      *http.Cookie
}

type resGoogleCallback struct {
	UserId      uuid.UUID
	RedirectURL string
	Cookie      *http.Cookie
}

type resToken struct {
	AccessToken  string
	RefreshToken string
}
