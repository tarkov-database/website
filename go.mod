module github.com/tarkov-database/website

go 1.15

replace github.com/tarkov-database/website/bundler => ./bundler

require (
	github.com/google/logger v1.1.0
	github.com/gorilla/websocket v1.4.2
	github.com/julienschmidt/httprouter v1.3.0
	github.com/tarkov-database/website/bundler v0.0.0-00010101000000-000000000000
	github.com/zeebo/blake3 v0.1.0
	golang.org/x/net v0.0.0-20201216054612-986b41b23924
	golang.org/x/text v0.3.5
)
