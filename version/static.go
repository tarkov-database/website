package version

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/google/logger"
	"github.com/zeebo/blake3"
)

const rootDir = "./static"

type FileSums map[string]string

var StaticSums FileSums

func init() {
	ch := make(chan File, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go hashFiles(rootDir, ch, wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	StaticSums = make(FileSums)
	prefix := rootDir + "/"
	for m := range ch {
		k := strings.TrimPrefix(m.Path, prefix)
		StaticSums[k] = m.Sum
	}
}

type File struct {
	Path, Sum string
}

var hasher = blake3.New()

func hashFile(p string, ch chan File, wg *sync.WaitGroup) {
	d, err := ioutil.ReadFile(p)
	if err != nil {
		logger.Fatal(err)
	}

	sum := hasher.Sum(d)
	defer hasher.Reset()

	ch <- File{p, hex.EncodeToString(sum)}

	wg.Done()
}

func hashFiles(p string, ch chan File, wg *sync.WaitGroup) {
	index, err := ioutil.ReadDir(p)
	if err != nil {
		logger.Fatal(err)
	}

	wg.Add(len(index))

	for _, e := range index {
		path := fmt.Sprintf("%s/%s", p, e.Name())
		if e.IsDir() {
			go hashFiles(path, ch, wg)
		} else {
			go hashFile(path, ch, wg)
		}
	}

	wg.Done()
}
