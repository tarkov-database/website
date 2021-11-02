//go:build DEBUG

package main

import (
	"github.com/tarkov-database/website/bundler"
	"github.com/tarkov-database/website/version"

	"github.com/google/logger"
)

const (
	sourceDir = "static/src"
	outDir    = "static/dist/resources"
)

func init() {
	watchScripts()
	watchStyles()
}

func watchScripts() {
	events, err := bundler.Watch(sourceDir+"/scripts", outDir+"/js", &bundler.BuildOptions{
		Sourcemap: true,
	})
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

func watchStyles() {
	events, err := bundler.Watch(sourceDir+"/styles", outDir+"/css", nil)
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
