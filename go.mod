module github.com/tarkov-database/website

go 1.18

replace github.com/tarkov-database/website/bundler => ./bundler

require (
	github.com/google/logger v1.1.1
	github.com/gorilla/websocket v1.5.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/tarkov-database/website/bundler v0.0.0-00010101000000-000000000000
	github.com/zeebo/blake3 v0.2.3
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5
	golang.org/x/text v0.3.7
)

require (
	github.com/evanw/esbuild v0.14.47 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
)
