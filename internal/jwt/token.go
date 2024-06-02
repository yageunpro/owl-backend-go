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

func ValidateToken(tokenString string) (uuid.UUID, error) {
	tokenClaims := new(claims)

	tok, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		iss, err := token.Claims.GetIssuer()
		if err != nil {
			return nil, errors.Join(errors.New("could not get issuer"), err)
		}
		if iss != Issuer {
			return nil, errors.New("invalid issuer")
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			return nil, errors.Join(errors.New("could not get subject"), err)
		}
		if sub != Subject {
			return nil, errors.New("invalid subject")
		}

		aud, err := token.Claims.GetAudience()
		if err != nil {
			return nil, errors.Join(errors.New("could not get audience"), err)
		}

		switch aud[0] {
		case AccessAud:
			return getSecretKey(access)
		case RefreshAud:
			return getSecretKey(refresh)
		}
		return nil, errors.New("invalid audience")
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.Nil, ErrInvalidToken
		}
		return uuid.Nil, errors.Join(errors.New("could not parse token"), err)
	}

	if tok.Valid {
		return tokenClaims.UserId, nil
	}

	return uuid.Nil, ErrInvalidToken
}
