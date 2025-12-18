package models

import "errors"

var (
	ErrInvalidIdentity = errors.New("invalid screen identity")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCredentialsExpired = errors.New("credentials expired")
	ErrCredentialsNotFound = errors.New("credentials not found")
)

