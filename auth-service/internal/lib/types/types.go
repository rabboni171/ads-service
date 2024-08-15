package types

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrAppNotFound        = errors.New("app not found")
	ErrEmailRequired      = errors.New("email is required")
	ErrPassRequired       = errors.New("password is required")
	ErrAppIdRequired       = errors.New("app_id is required")
)
