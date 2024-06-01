package db

import "errors"

var ErrDBOpen = errors.New("fail open database")
var ErrDBPing = errors.New("fail ping database")
