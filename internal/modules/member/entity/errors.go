package entity

import "errors"

var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrNameTooShort       = errors.New("name too short")
)
