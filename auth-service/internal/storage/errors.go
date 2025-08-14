package storage

import "errors"

var (
	ErrKeyNotFound           = errors.New("key not found")
	ErrRedisConnectionFailed = errors.New("redis connection failed")
)
