package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"time"
)

type claims struct {
	jwt.RegisteredClaims
	UserId uuid.UUID `json:"uid"`
}

func NewAccessToken(userId uuid.UUID) (string, error) {
	key, err := getSecretKey(access)
	if err != nil {
		return "", err
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   Subject,
			Audience:  []string{AccessAud},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(AccessValidDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		UserId: userId,
	}).SignedString(key)

	if err != nil {
		return "", errors.Join(errors.New("could not sign access token"), err)
	}
	return token, nil
}

func NewRefreshToken(userId uuid.UUID) (string, error) {
	key, err := getSecretKey(refresh)
	if err != nil {
		return "", err
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   Subject,
			Audience:  []string{RefreshAud},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(RefreshValidDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		UserId: userId,
	}).SignedString(key)

	if err != nil {
		return "", errors.Join(errors.New("could not sign refresh token"), err)
	}
	return token, nil
}
