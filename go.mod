module github.com/tarkov-database/website

go 1.21

replace github.com/tarkov-database/website/bundler => ./bundler

require (
	github.com/goccy/go-json v0.10.2
	github.com/google/logger v1.1.1
	github.com/gorilla/websocket v1.5.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/tarkov-database/website/bundler v0.0.0-20231024161109-e647b533174f
	github.com/zeebo/blake3 v0.2.3
	golang.org/x/net v0.17.0
	golang.org/x/text v0.13.0
)

require (
	github.com/evanw/esbuild v0.19.5 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	golang.org/x/sys v0.13.0 // indirect
)
