APP_NAME := ssh-aliases
PACKAGES := $(shell go list ./... | grep -v /vendor/)
BUILD_FOLDER := target
DIST_FOLDER := dist

GIT_REV=$(shell git rev-parse --verify --short HEAD)

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all fmt clean test build release lint lint-deps

all: clean build

clean:
	@go clean -v .
	@rm -rf $(BUILD_FOLDER)
	@rm -rf $(DIST_FOLDER)

test:
	@go test -cover ./...

build: test
	@go build -o $(BUILD_FOLDER)/$(APP_NAME) -ldflags '-s -w -X main.Version=$(GIT_REV) -extldflags "-static"'

release: clean lint build
	@bash ./package.sh

fmt: lint-deps
	@goimports -w $(SRC)
	@gofmt -l -s -w $(SRC)

lint: lint-deps
	@gometalinter.v1 --config=gometalinter.json ./...

lint-deps:
	@which gometalinter.v1 > /dev/null || (go get gopkg.in/alecthomas/gometalinter.v1 && gometalinter.v1 --install)
