module github.com/tarkov-database/website

go 1.20

replace github.com/tarkov-database/website/bundler => ./bundler

require (
	github.com/google/logger v1.1.1
	github.com/gorilla/websocket v1.5.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/tarkov-database/website/bundler v0.0.0-20220912123526-213d2c143c6c
	github.com/zeebo/blake3 v0.2.3
	golang.org/x/net v0.8.0
	golang.org/x/text v0.8.0
)

require (
	github.com/evanw/esbuild v0.17.13 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	golang.org/x/sys v0.6.0 // indirect
)
