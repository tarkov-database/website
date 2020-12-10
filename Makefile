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

STATIC_SRC := static/src
STATIC_DIST := static/dist
NODE_MODULES := node_modules

all: run

bin:
	go build -v -o ${OUT} -tags="${BUILD_TAGS}" -ldflags="-X ${VERSION_PKG}.buildDate=${BUILD_DATE} -X ${VERSION_PKG}.commitDate=${COMMIT_DATE} -X ${VERSION_PKG}.branch=${BRANCH_NAME} -X ${VERSION_PKG}.commitLong=${COMMIT_LONG} -X ${VERSION_PKG}.commitShort=${COMMIT_SHORT}"

debug: BUILD_TAGS += DEBUG
debug: bin

statics:
	mkdir -p ${STATIC_DIST}/resources/css ${STATIC_DIST}/resources/js ${STATIC_DIST}/resources/fonts ${STATIC_DIST}/resources/img ${STATIC_DIST}/resources/style
	cp -r ${STATIC_SRC}/fonts/* ${STATIC_DIST}/resources/fonts
	cp -r ${STATIC_SRC}/images/* ${STATIC_DIST}/resources/img
	cp -r ${STATIC_SRC}/map-styles/* ${STATIC_DIST}/resources/style
	cp ${STATIC_SRC}/manifest.json ${STATIC_DIST}/resources/
	go run ${MAIN_PKG}/bundler/cmd -source=${STATIC_SRC}/scripts -out=${STATIC_DIST}/resources/js -sourcemap
	go run ${MAIN_PKG}/bundler/cmd -source=${STATIC_SRC}/styles -out=${STATIC_DIST}/resources/css

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
