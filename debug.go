// +build DEBUG

package main

import (
	"github.com/tarkov-database/website/bundler"
	"github.com/tarkov-database/website/version"

	"github.com/google/logger"
)

const (
	sourceDir = "static/src/scripts"
	outDir    = "static/dist/resources/js"
)

func init() {
	opts := &bundler.BuildOptions{
		Sourcemap: true,
	}

	events, err := bundler.Watch(sourceDir, outDir, opts)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		for event := range events {
			if err := event.Error; err != nil {
				logger.Error(err)
			}

			if err := version.RefreshSumOf(event.Filename); err != nil {
				logger.Error(err)
			}
		}
	}()
}
