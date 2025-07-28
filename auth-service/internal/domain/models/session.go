package models

type Session struct {
	ID        string
	Email     string
	ExpiresAt int64
}
