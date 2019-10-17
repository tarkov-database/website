package version

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/blake2b"

	"github.com/google/logger"
)

const rootDir = "./static/resources/"

type FileSums = map[string]string

var StaticSums = FileSums{}

func init() {
	index, err := ioutil.ReadDir(rootDir)
	if err != nil {
		logger.Fatal(err)
	}

	ch := make(chan FileSums)
	for _, e := range index {
		path := fmt.Sprintf("%s%s", rootDir, e.Name())
		if e.IsDir() {
			go hashFiles(path, ch)
		} else {
			StaticSums[e.Name()] = hashFile(path)
		}
	}

	i := 0
	for m := range ch {
		for k, v := range m {
			k = strings.TrimPrefix(k, rootDir)
			StaticSums[k] = v
		}

		i++
		if i == len(index)-1 {
			close(ch)
		}
	}
}

func hashFile(p string) string {
	d, err := ioutil.ReadFile(p)
	if err != nil {
		logger.Fatal(err)
	}
	sum := blake2b.Sum256(d)

	return hex.EncodeToString(sum[:])
}

func hashFiles(p string, ch chan FileSums) {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		logger.Fatal(err)
	}

	hashes := map[string]string{}
	for _, file := range files {
		if !file.IsDir() {
			filePath := fmt.Sprintf("%s/%s", p, file.Name())
			hashes[filePath] = hashFile(filePath)
		}
	}

	ch <- hashes
}
