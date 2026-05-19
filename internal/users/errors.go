package users

import "errors"

var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrInvalidInput     = errors.New("invalid input")
	ErrEmailExists      = errors.New("email already registered")
	ErrCannotDeleteSelf = errors.New("cannot delete own account")
)
