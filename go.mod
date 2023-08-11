module github.com/tarkov-database/website

go 1.21

replace github.com/tarkov-database/website/bundler => ./bundler

require (
	github.com/goccy/go-json v0.10.2
	github.com/google/logger v1.1.1
	github.com/gorilla/websocket v1.5.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/tarkov-database/website/bundler v0.0.0-20230701115656-2ea42cd57839
	github.com/zeebo/blake3 v0.2.3
	golang.org/x/net v0.14.0
	golang.org/x/text v0.12.0
)

require (
	github.com/evanw/esbuild v0.19.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	golang.org/x/sys v0.11.0 // indirect
)
