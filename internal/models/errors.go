package models

import "errors"

var (
	// ErrNoRecord is used if no matching snippet record is found.
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials is used if a user tries to log in with an invalid email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail is used if a user tries to sign up with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
