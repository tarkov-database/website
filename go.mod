module github.com/tarkov-database/website

go 1.17

replace github.com/tarkov-database/website/bundler => ./bundler

require (
	github.com/google/logger v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/julienschmidt/httprouter v1.3.0
	github.com/tarkov-database/website/bundler v0.0.0-00010101000000-000000000000
	github.com/zeebo/blake3 v0.2.1
	golang.org/x/net v0.0.0-20210903162142-ad29c8ab022f
	golang.org/x/text v0.3.7
)

require (
	github.com/evanw/esbuild v0.12.25 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)
