package model

import (
	"net/http"
)

type EntityType int

const (
	TypeItem EntityType = iota
	TypeLocation
	TypeFeature
	TypeMixed
)

var entityTypeString = [...]string{
	"item",
	"location",
	"feature",
	"mixed",
}

func (et EntityType) String() string {
	return entityTypeString[et]
}

type Filter interface {
	GetAll() map[string][]string
	Get(string) []string
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
