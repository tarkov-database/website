package version

import (
	"errors"

	"github.com/tarkov-database/website/utils"
)

const staticDir = "static/dist"

var staticSums utils.FileSums

func init() {
	RefreshSumAll()
}

func SumOf(path string) (string, error) {
	sum, ok := staticSums[path]
	if !ok {
		return "", errors.New("path not found")
	}

	return sum, nil
}

func RefreshSumOf(path string) error {
	sum, err := utils.SumFile(path, staticDir)
	if err != nil {
		return err
	}
	if _, ok := staticSums[sum.Path]; !ok {
		return errors.New("file does not exist")
	}

	staticSums[sum.Path] = sum.Sum

	return nil
}

func RefreshSumAll() {
	staticSums = utils.SumDir(staticDir, &utils.SumOptions{
		BasePath:  staticDir,
		Recursive: true,
	})
}
