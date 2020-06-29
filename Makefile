-include .env

OUT := frontendserver

BUILD_DATE := $(shell date +%s)
BRANCH_NAME := $(shell if [ -z "${BRANCH}" ]; then git rev-parse --abbrev-ref HEAD; else echo "${BRANCH}"; fi)
COMMIT_LONG := $(shell git rev-parse HEAD)
COMMIT_SHORT := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format="%ct")

VERSION_PKG := github.com/tarkov-database/website/version

all: run

bin:
	go build -v -o ${OUT} -ldflags="-X ${VERSION_PKG}.buildDate=${BUILD_DATE} -X ${VERSION_PKG}.commitDate=${COMMIT_DATE} -X ${VERSION_PKG}.branch=${BRANCH_NAME} -X ${VERSION_PKG}.commitLong=${COMMIT_LONG} -X ${VERSION_PKG}.commitShort=${COMMIT_SHORT}"

lint:
	revive -config revive.toml -formatter stylish ./...
	test -z $(shell gofmt -l .) || (gofmt -l . && exit 1)
	npm run lint

fmt:
	go fmt ./...

run: bin
	./${OUT}

clean:
	-@rm ${OUT}
