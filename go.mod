module github.com/tarkov-database/website

go 1.17

replace github.com/tarkov-database/website/bundler => ./bundler

require (
	github.com/google/logger v1.1.1
	github.com/gorilla/websocket v1.5.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/tarkov-database/website/bundler v0.0.0-00010101000000-000000000000
	github.com/zeebo/blake3 v0.2.2
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f
	golang.org/x/text v0.3.7
)

require (
	github.com/evanw/esbuild v0.14.23 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	golang.org/x/sys v0.0.0-20210908233432-aa78b53d3365 // indirect
)
