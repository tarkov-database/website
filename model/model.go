package model

import (
	"net/http"
)

// Response describes a status response
type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
}

// NewResponse creates a new status response based on parameters
func NewResponse(msg string, code int) *Response {
	return &Response{
		Status:     http.StatusText(code),
		Message:    msg,
		StatusCode: code,
	}
}
