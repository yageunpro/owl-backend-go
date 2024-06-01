package auth

import "net/http"

type GoogleCallbackParam struct {
	Cookie *http.Cookie
	State  string
	Code   string
	Scope  []string
}
