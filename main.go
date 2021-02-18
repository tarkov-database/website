package main

import (
	"fmt"
	"io"

	"github.com/tarkov-database/website/core/server"
	"github.com/tarkov-database/website/version"

	"github.com/google/logger"
)

func main() {
	fmt.Printf("Starting up Tarkov Database Frontend (Commit: %s-%s Build Date: %s)\n\n",
		version.App.CommitShort, version.App.BranchName, version.App.BuildDate)

	defLog := logger.Init("default", true, false, io.Discard)
	defer defLog.Close()

	server.Start()
}
