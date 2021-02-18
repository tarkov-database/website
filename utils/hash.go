package utils

import (
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/logger"
	"github.com/zeebo/blake3"
)

type FileSums map[string]string

type FileSum struct {
	Path, Sum string
}

type SumOptions struct {
	Extensions []string
	BasePath   string
	Recursive  bool
}

func SumDir(dir string, opts *SumOptions) FileSums {
	ch := make(chan FileSum, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go hashFiles(dir, opts, ch, wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	sums := make(FileSums)
	for m := range ch {
		k := m.Path
		if opts.BasePath != "" {
			k, _ = filepath.Rel(opts.BasePath, m.Path)
		}

		sums[k] = m.Sum
	}

	return sums
}

func SumFile(path string, basePath ...string) (FileSum, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return FileSum{}, err
	}

	if len(basePath) > 0 {
		path, err = filepath.Rel(basePath[0], path)
		if err != nil {
			return FileSum{}, err
		}
	}

	h := blake3.New()
	h.Write(d)

	return FileSum{path, hex.EncodeToString(h.Sum(nil))}, nil
}

func hashFiles(p string, opts *SumOptions, ch chan FileSum, wg *sync.WaitGroup) {
	index, err := os.ReadDir(p)
	if err != nil {
		logger.Fatal(err)
	}

	for _, e := range index {
		path := filepath.Join(p, e.Name())
		if e.IsDir() && opts.Recursive {
			wg.Add(1)
			go hashFiles(path, opts, ch, wg)
		} else {
			if len(opts.Extensions) == 0 || containsStr(opts.Extensions, filepath.Ext(e.Name())) {
				wg.Add(1)
				go hashFile(path, ch, wg)
			}
		}
	}

	wg.Done()
}

func hashFile(p string, ch chan FileSum, wg *sync.WaitGroup) {
	file, err := SumFile(p)
	if err != nil {
		log.Fatal(err)
	}

	ch <- file

	wg.Done()
}

func containsStr(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}

	return false
}
