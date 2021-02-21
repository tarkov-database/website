package health

import (
	"context"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/core/search"
	"github.com/tarkov-database/website/model/item"
)

var (
	apiStatus    = Failure
	searchStatus = Failure
)

// APIStatus returns the database status of the last health check
func APIStatus() Status {
	return apiStatus
}

// SearchStatus returns the database status of the last health check
func SearchStatus() Status {
	return searchStatus
}

func scheduler(t *time.Ticker, c chan os.Signal) {
	for {
		select {
		case <-t.C:
			updateStatus()
		case <-c:
			t.Stop()
			return
		}
	}
}

func updateStatus() {
	apiStatus = getAPIStatus()
	searchStatus = getSearchStatus()
}

func getAPIStatus() Status {
	timeout := 30 * time.Second
	if cfg.updateInterval < timeout {
		timeout = cfg.updateInterval
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	if latency > cfg.latencyThreshold {
		logger.Warningf("API latency exceeds threshold with %s", latency)
		return Warning
	}

	return OK
}

func getSearchStatus() Status {
	timeout := 30 * time.Second
	if cfg.updateInterval < timeout {
		timeout = cfg.updateInterval
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	if latency > cfg.latencyThreshold {
		logger.Warningf("Search server latency exceeds threshold with %s", latency)
		return Warning
	}

	return OK
}
