package search

import (
	"context"
	"errors"
)

var global Client

func checkGlobal() error {
	if global == nil {
		return errors.New("global client is not initialized")
	}

	return nil
}

func InitGlobal(cli Client) {
	global = cli
}

func Query(ctx context.Context, queryString string) (*QueryResult, error) {
	err := checkGlobal()
	if err != nil {
		return nil, err
	}

	return global.Query(ctx, queryString)
}
