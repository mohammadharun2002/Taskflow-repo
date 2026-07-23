package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNameRequired       = errors.New("name is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrPasswordTooShort   = errors.New("password must contain at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
