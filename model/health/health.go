package health

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/core/search"
	"github.com/tarkov-database/website/model/item"
)

var latencyThreshold time.Duration

func init() {
	if env := os.Getenv("UNHEALTHY_LATENCY"); len(env) > 0 {
		d, err := time.ParseDuration(env)
		if err != nil {
			log.Printf("Unhealthy latency value is not valid: %s\n", err)
			os.Exit(2)
		}
		latencyThreshold = d
	} else {
		latencyThreshold = 300 * time.Millisecond
	}
}

// Status represents the status code of a service
type Status int

const (
	// OK status if all checks were successful
	OK Status = iota

	// Warning status if non-critical issues are discovered
	Warning

	// Failure status when critical problems are discovered
	Failure
)

// Health represents the object of the health root endpoint
type Health struct {
	OK      bool     `json:"ok"`
	Service *Service `json:"service"`
}

// Service holds all services with their respective status
type Service struct {
	API    Status `json:"rest-api"`
	Search Status `json:"search"`
}

// GetHealth performs a self-check and returns the result
func GetHealth() *Health {
	svc := &Service{}

	health := &Health{
		OK:      true,
		Service: svc,
	}

	svc.API = getAPIStatus()
	if svc.API != OK {
		health.OK = false
	}

	svc.Search = getSearchStatus()
	if svc.Search != OK {
		health.OK = false
	}

	return health
}

func getAPIStatus() Status {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := &api.Options{
		Filter: map[string]string{
			"skipKinds": "1",
		},
	}

	start := time.Now()

	if err := api.GET(ctx, "/item", opts, &item.Index{}); err != nil {
		logger.Errorf("Error while checking API connection: %s", err)
		return Failure
	}

	latency := time.Since(start)
	if latency > latencyThreshold {
		logger.Warningf("API latency exceeds threshold with %s", latency)
		return Warning
	}

	return OK
}

func getSearchStatus() Status {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := &search.Query{
		Query: "colt 1911",
	}

	opts := &search.Options{
		Limit: 1,
	}

	start := time.Now()

	if _, err := search.Search(ctx, query, opts); err != nil {
		logger.Errorf("Error while checking search server connection: %s", err)
		return Failure
	}

	latency := time.Since(start)
	if latency > latencyThreshold {
		logger.Warningf("Search server latency exceeds threshold with %s", latency)
		return Warning
	}

	return OK
}
