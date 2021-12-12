.PHONY: build run
.PHONY: test-all test-watcher test-event test-utils

VERSION := test-build
BUILD := $$(git log -1 --pretty=%h)
BUILD_TIME := $$(date -u +"%Y%m%d.%H%M%S")

build:
	@go build \
		-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}" \
		-o ./bin/switchboard ./cmd

run:
	@go run ./cmd/main.go

test-watcher:
	@gotest -v ./watcher

test-event:
	@gotest -v ./event

test-utils:
	@gotest -v ./utils

test-all: test-watcher test-event test-utils
