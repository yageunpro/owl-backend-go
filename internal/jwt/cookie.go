package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type CookieData struct {
	AccessToken  string `json:"a"`
	RefreshToken string `json:"r"`
}

func ToCookie(accessToken string, refreshToken string) (*http.Cookie, error) {
	data := CookieData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal cookie data"), err)
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

func FromCookie(cookie *http.Cookie) (*CookieData, error) {
	if cookie == nil {
		return nil, errors.New("cookie is nil")
	}

	raw, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, errors.Join(errors.New("failed to decode cookie data"), err)
	}

	data := new(CookieData)
	err = json.Unmarshal(raw, data)
	if err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal cookie data"), err)
	}

	return data, nil
}
