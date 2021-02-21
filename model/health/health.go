package health

import (
	core "github.com/tarkov-database/website/core/health"
)

// Health represents the object of the health root endpoint
type Health struct {
	OK      bool     `json:"ok"`
	Service *Service `json:"service"`
}

// Service holds all services with their respective status
type Service struct {
	API    core.Status `json:"rest-api"`
	Search core.Status `json:"search"`
}

// GetHealth performs a self-check and returns the result
func GetHealth() *Health {
	svc := &Service{}

	h := &Health{
		OK:      true,
		Service: svc,
	}

	svc.API = core.APIStatus()
	if svc.API != core.OK {
		h.OK = false
	}

	svc.Search = core.SearchStatus()
	if svc.Search != core.OK {
		h.OK = false
	}

	return h
}
