package bundler

import (
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
	"github.com/google/logger"
)

var defaultOptions = api.BuildOptions{
	EntryPoints:       nil,
	Bundle:            true,
	MinifyIdentifiers: true,
	MinifyWhitespace:  true,
	MinifySyntax:      true,
	Sourcemap:         api.SourceMapLinked,
	Target:            api.ESNext,
	Platform:          api.PlatformBrowser,
	Format:            api.FormatESModule,
	MainFields:        []string{"module", "browser", "main"},
	Incremental:       false,
	Write:             true,
}

type BuildOptions struct {
	Sourcemap bool
}

func Build(source string, out string, opts *BuildOptions) error {
	modules, err := getModules(source, out)
	if err != nil {
		return err
	}

	options := defaultOptions
	if opts != nil && !opts.Sourcemap {
		options.Sourcemap = api.SourceMapNone
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(modules))

	for _, m := range modules {
		go func(m module) {
			newBuild(m, options)
			wg.Done()
		}(m)
	}

	wg.Wait()

	return nil
}

func Watch(source string, out string, opts *BuildOptions) (chan BuildEvent, error) {
	modules, err := getModules(source, out)
	if err != nil {
		return nil, err
	}

	options := defaultOptions
	options.Incremental = true
	if opts != nil && !opts.Sourcemap {
		options.Sourcemap = api.SourceMapNone
	}

	events := make(chan BuildEvent, 1)
	builders := make(map[string]*builder)

	wg := &sync.WaitGroup{}
	mutex := sync.Mutex{}

	wg.Add(len(modules))

	for _, m := range modules {
		go func(m module) {
			dir := filepath.Dir(m.entryPoint)
			builder := newBuilder(newBuild(m, options), events)
			mutex.Lock()
			builders[dir] = builder
			mutex.Unlock()
			wg.Done()
		}(m)
	}

	wg.Wait()

	if err = watchChanges(builders); err != nil {
		return nil, err
	}

	return events, nil
}

type module struct {
	entryPoint string
	outFile    string
}

func getModules(dir string, out string) ([]module, error) {
	index, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	modules := make([]module, 0)
	for _, e := range index {
		if e.IsDir() {
			path := filepath.Join(dir, e.Name())
			file := filepath.Join(path, "index.ts")
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				modules = append(modules, module{
					entryPoint: file,
					outFile: filepath.Join(
						out,
						e.Name()+".js",
					),
				})
			}
		}
	}

	return modules, nil
}

func newBuild(mod module, opts api.BuildOptions) *api.BuildResult {
	opts.EntryPoints = []string{mod.entryPoint}
	opts.Outfile = mod.outFile

	result := api.Build(opts)
	logMessages(&result)

	return &result
}

func logMessages(result *api.BuildResult) {
	for _, warn := range result.Warnings {
		logger.Warningf("%s %v:%v %s",
			warn.Location.File,
			warn.Location.Line,
			warn.Location.Column,
			warn.Text,
		)
	}
	for _, err := range result.Errors {
		logger.Errorf("%s %v:%v %s",
			err.Location.File,
			err.Location.Line,
			err.Location.Column,
			err.Text,
		)
	}
}

type BuildEvent struct {
	Filename string
	Error    error
}

type builder struct {
	delay  time.Duration
	build  *api.BuildResult
	ticker *time.Ticker
	events chan BuildEvent
	close  chan bool
}

func (b *builder) Close() {
	close(b.close)
}

func (b *builder) Rebuild() {
	if b.ticker != nil {
		b.ticker.Reset(b.delay)
		return
	}

	b.ticker = time.NewTicker(b.delay)

	go func() {
		defer b.ticker.Stop()

		for {
			select {
			case _, ok := <-b.ticker.C:
				if !ok {
					return
				}

				start := time.Now()

				res := b.build.Rebuild()

				elapsed := time.Since(start)

				logMessages(&res)

				for _, out := range res.OutputFiles {
					cwd, err := os.Getwd()
					if err != nil {
						b.events <- BuildEvent{Error: err}
						return
					}

					file, err := filepath.Rel(cwd, out.Path)
					if err != nil {
						b.events <- BuildEvent{Error: err}
						return
					}

					logger.Infof("File %s built in %s", file, elapsed)

					b.events <- BuildEvent{Filename: file}
				}

				b.build = &res

				b.ticker.Stop()
			case <-b.close:
				return
			}
		}
	}()
}

func newBuilder(result *api.BuildResult, events chan BuildEvent) *builder {
	return &builder{
		delay:  300 * time.Millisecond,
		build:  result,
		events: events,
		close:  make(chan bool, 1),
	}
}

func watchChanges(builders map[string]*builder) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				builder, ok := builders[filepath.Dir(event.Name)]
				if !ok {
					continue
				}

				if event.Op == fsnotify.Write {
					builder.Rebuild()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error(err)
			case <-sig:
				for _, b := range builders {
					b.Close()
				}
				return
			}
		}
	}()

	for dir := range builders {
		err := watcher.Add(dir)
		if err != nil {
			logger.Fatal(err)
		}
	}

	return nil
}
