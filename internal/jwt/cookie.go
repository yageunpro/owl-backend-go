package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type cookieData struct {
	AccessToken  string `json:"a"`
	RefreshToken string `json:"r"`
}

func ToCookie(accessToken string, refreshToken string) (*http.Cookie, error) {
	data := cookieData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal cookieData"), err)
	}

	cookie := http.Cookie{
		Name:     CookieKey,
		Value:    base64.URLEncoding.EncodeToString(raw),
		Path:     "/",
		Expires:  time.Now().Add(RefreshValidDuration),
		HttpOnly: true,
	}

	return &cookie, nil
}
