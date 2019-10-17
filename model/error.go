package model

import (
	"errors"
)

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidKind     = errors.New("kind is not valid")
	ErrInvalidCategory = errors.New("kind is not valid")
	ErrInternalError   = errors.New("server or network error")
)
