package version

import (
	"strconv"
	"time"

	"github.com/google/logger"
)

var (
	App *Application

	buildDate = "0000000000"

	commitLong  = ""
	commitShort = ""
	commitDate  = "0000000000"

	branch = ""
)

type Application struct {
	BuildDate time.Time

	CommitDate  time.Time
	CommitLong  string
	CommitShort string

	BranchName string
}

func init() {
	App = &Application{
		BuildDate:   parseTime(buildDate),
		CommitDate:  parseTime(commitDate),
		CommitLong:  commitLong,
		CommitShort: commitShort,
		BranchName:  branch,
	}
}

func parseTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logger.Error(err)
		return time.Time{}
	}

	return time.Unix(i, 0)
}
