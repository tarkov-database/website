package model

import (
	"errors"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrInternalError = errors.New("server or network error")
)
