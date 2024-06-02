package jwt

import "time"

type ContextKey string

const AuthnCtxKey ContextKey = "authn"
const (
	AccessValidDuration  = 10 * time.Minute
	RefreshValidDuration = 7 * 24 * time.Hour
	CookieKey            = "AUTHN"
	Issuer               = "yageun.pro"
	Subject              = "owl-authn"
	RefreshAud           = "refresh"
	AccessAud            = "access"
)
