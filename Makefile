-include .env

OUT := frontendserver

MAIN_PKG := github.com/tarkov-database/website
VERSION_PKG := ${MAIN_PKG}/version

BUILD_TAGS := $()

BUILD_DATE := $(shell date +%s)
BRANCH_NAME := $(shell if [ -z "${BRANCH}" ]; then git rev-parse --abbrev-ref HEAD; else echo "${BRANCH}"; fi)
COMMIT_LONG := $(shell git rev-parse HEAD)
COMMIT_SHORT := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format="%ct")

all: run

bin:
	go build -v -o ${OUT} -tags="${BUILD_TAGS}" -ldflags="-X ${VERSION_PKG}.buildDate=${BUILD_DATE} -X ${VERSION_PKG}.commitDate=${COMMIT_DATE} -X ${VERSION_PKG}.branch=${BRANCH_NAME} -X ${VERSION_PKG}.commitLong=${COMMIT_LONG} -X ${VERSION_PKG}.commitShort=${COMMIT_SHORT}"

debug: BUILD_TAGS += DEBUG
debug: bin

statics:
	go run ${MAIN_PKG}/bundler/cmd -source=static/src -out=static/public/resources/js -sourcemap
	cp node_modules/mapbox-gl/dist/mapbox-gl.css static/public/resources/css/mapbox-gl.min.css

lint:
	revive -config revive.toml -formatter stylish ./...
	test -z $(shell gofmt -l .) || (gofmt -l . && exit 1)
	npm run lint

fmt:
	go fmt ./...
	npm run fmt

test:
	npm run test

run: debug
	./${OUT}

clean:
	-@rm ${OUT}
	-@rm -rf node_modules
