package main

import (
	"flag"
	"io"

	"github.com/tarkov-database/website/bundler"

	"github.com/google/logger"
)

func main() {
	source := flag.String("source", "", "Source directory")
	out := flag.String("out", "", "Out directory")
	sourcemap := flag.Bool("sourcemap", false, "Write sourcemap")

	flag.Parse()

	defLog := logger.Init("default", true, false, io.Discard)
	defer defLog.Close()

	if *source == "" {
		logger.Fatal("Source directory is missing")
	}
	if *out == "" {
		logger.Fatal("Out directory is missing")
	}

	opts := &bundler.BuildOptions{
		Sourcemap: *sourcemap,
	}

	err := bundler.Build(*source, *out, opts)
	if err != nil {
		logger.Fatal(err)
	}
}
