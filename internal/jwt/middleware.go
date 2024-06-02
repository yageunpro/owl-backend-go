package jwt

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	headerLogic := func(c echo.Context) error {
		accessToken := c.Request().Header.Get("access_token")
		refreshToken := c.Request().Header.Get("refresh_token")

		if accessToken == "" {
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), AuthnCtxKey, uuid.Nil)))
			return next(c)
		}

		userId, err := ValidateToken(accessToken)
		if errors.Is(err, ErrInvalidToken) {
			if refreshToken == "" {
				return echo.ErrUnauthorized
			}

			userId, err = ValidateToken(refreshToken)
			if err != nil {
				return echo.ErrUnauthorized
			}

			newToken, err := NewAccessToken(userId)
			if err != nil {
				return err
			}
			c.Response().Header().Set("access_token", newToken)
		} else if err != nil {
			return echo.ErrUnauthorized
		}

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), AuthnCtxKey, userId)))
		return next(c)
	}

	return func(c echo.Context) error {
		cookie, err := c.Request().Cookie(CookieKey)
		if errors.Is(err, http.ErrNoCookie) {
			return headerLogic(c)
		} else if err != nil {
			return err
		}

		tok, err := FromCookie(cookie)

		userId, err := ValidateToken(tok.AccessToken)
		if errors.Is(err, ErrInvalidToken) {
			userId, err = ValidateToken(tok.RefreshToken)
			if err != nil {
				return echo.ErrUnauthorized
			}

			newToken, err := NewAccessToken(userId)
			if err != nil {
				return err
			}

			newCookie, err := ToCookie(newToken, tok.RefreshToken)
			if err != nil {
				return err
			}
			c.SetCookie(newCookie)
		} else if err != nil {
			return echo.ErrUnauthorized
		}

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), AuthnCtxKey, userId)))
		return next(c)
	}
}

func GetUserID(ctx context.Context) uuid.UUID {
	userId, ok := ctx.Value(AuthnCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return userId
}
