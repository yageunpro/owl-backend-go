package oauth

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
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
