package auth

import "errors"

var (
	ErrIncorrectCode   = errors.New("verification code is incorrect")
	ErrCodeExpired     = errors.New("verification code has expired")
	ErrSessionExpired  = errors.New("session has expired")
	ErrSessionNotFound = errors.New("session not found")

	ErrCodeSaveFailed      = errors.New("failed to save verification code")
	ErrSessionCreateFailed = errors.New("failed to create session")
	ErrCodeGenFailed       = errors.New("failed to generate verification code")
)
