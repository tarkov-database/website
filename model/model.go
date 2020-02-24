package model

import (
	"net/http"
)

type EntityType int

const (
	Item EntityType = iota
	Location
	Feature
)

var entityTypeString = [...]string{
	"item",
	"location",
	"feature",
}

func (et EntityType) String() string {
	return entityTypeString[et]
}

// Response describes a JSON status response
type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
}

// NewResponse creates a new JSON status response based on parameters
func NewResponse(msg string, code int) *Response {
	return &Response{
		Status:     http.StatusText(code),
		Message:    msg,
		StatusCode: code,
	}
}
