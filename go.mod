module github.com/tarkov-database/website

go 1.19

replace github.com/tarkov-database/website/bundler => ./bundler

require (
	github.com/google/logger v1.1.1
	github.com/gorilla/websocket v1.5.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/tarkov-database/website/bundler v0.0.0-00010101000000-000000000000
	github.com/zeebo/blake3 v0.2.3
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591
	golang.org/x/text v0.8.0
)

require (
	github.com/evanw/esbuild v0.15.7 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	golang.org/x/sys v0.5.0 // indirect
)
