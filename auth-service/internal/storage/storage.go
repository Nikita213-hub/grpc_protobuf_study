package storage

import "errors"

var (
	ErrTokenNotFound      = errors.New("token not found")
	ErrTokenAlreadyExists = errors.New("toke already exists")
)
