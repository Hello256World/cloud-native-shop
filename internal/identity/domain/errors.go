package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrPhoneTaken         = errors.New("phone already registered")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")

	ErrRefreshTokenInvalid = errors.New("refresh token invalid or expired")

	ErrAddressNotFound = errors.New("address not found")
)
