package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrPasswordRequired = errors.New("models: password required")
)
