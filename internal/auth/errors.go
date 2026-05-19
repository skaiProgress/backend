package auth

import "errors"

var (
	// ErrInvalidCredentials is returned when email or password is wrong.
	ErrInvalidCredentials = errors.New("invalid login credentials")
	// ErrUserBanned is returned when the account is banned.
	ErrUserBanned = errors.New("user is banned")
	// ErrUnauthorized is returned when the token is missing or invalid.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden is returned when the user lacks permission.
	ErrForbidden = errors.New("forbidden")
	// ErrWrongPassword is returned when the current password does not match.
	ErrWrongPassword = errors.New("wrong password")
)
