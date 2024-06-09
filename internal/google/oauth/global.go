package oauth

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"time"
)

var global OAuth

func checkGlobal() error {
	if global == nil {
		return errors.New("global oauth is not initialized")
	}

	return nil
}

func InitGlobal(auth OAuth) {
	global = auth
}

func AuthCodeURL(state string, isForce bool) (string, error) {
	err := checkGlobal()
	if err != nil {
		return "", err
	}

	return global.AuthCodeURL(state, isForce), nil
}

func Token(ctx context.Context, code string) (*oauth2.Token, error) {
	err := checkGlobal()
	if err != nil {
		return nil, err
	}

	return global.Token(ctx, code)
}

func IsAllowSync(scope []string) (bool, error) {
	err := checkGlobal()
	if err != nil {
		return false, err
	}

	return global.IsAllowSync(scope), nil
}

func TokenSource(ctx context.Context, token *oauth2.Token) (oauth2.TokenSource, error) {
	err := checkGlobal()
	if err != nil {
		return nil, err
	}

	return global.Config().TokenSource(ctx, token), nil
}

func ToToken(accessToken, refreshToken string, exp time.Time) *oauth2.Token {
	tok := oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "",
		RefreshToken: refreshToken,
		Expiry:       exp.UTC(),
	}
	return &tok
}
