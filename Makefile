APP_NAME := ssh-aliases
APP_VERSION := $(shell cat VERSION)

PACKAGES := $(shell go list ./...)
BUILD_FOLDER := target
DIST_FOLDER := dist

GIT_DESC := $(shell git describe)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

SRC := $(shell find . -type f -name '*.go')

.PHONY: all version fmt clean test build release lint lint-deps

ifneq "$(APP_VERSION)" "$(GIT_DESC)"
    APP_VERSION := $(GIT_DESC)-$(GIT_BRANCH)
endif

all: clean build

version:
	@echo $(APP_VERSION)

clean:
	@go clean -v .
	@rm -rf $(BUILD_FOLDER)
	@rm -rf $(DIST_FOLDER)

test:
	@go test -cover ./...

build:
	@go build -o $(BUILD_FOLDER)/$(APP_NAME) \
	-ldflags '-s -w -X main.Version=$(APP_VERSION)'

release: clean lint build
	@bash ./package.sh $(APP_VERSION)

fmt:
	@goimports -w $(SRC)
	@gofmt -l -s -w $(SRC)

lint:
	@golangci-lint run

