package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Error for an incorrect email or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// Error for an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
