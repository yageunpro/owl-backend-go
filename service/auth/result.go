package auth

import "net/http"

type resGoogleLogin struct {
	RedirectURL string
	Cookie      *http.Cookie
}

type resGoogleCallback struct {
	RedirectURL string
	Cookie      *http.Cookie
}

type resToken struct {
	AccessToken  string
	RefreshToken string
}
