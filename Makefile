.PHONY: build build-arm build-debian run
.PHONY: test-all test-watcher test-watcher-observe test-event test-utils test-cmd lint-all vet-all

VERSION := test-build
BUILD := $$(git log -1 --pretty=%h)
BUILD_TIME := $$(date -u +"%Y%m%d.%H%M%S")

build:
	@go build \
		-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}" \
		-o ./bin/switchboard ./cmd

build-debian:
	@GOOS=linux GOARCH=amd64 go build \
		-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}" \
		-o ./bin/switchboard ./cmd

build-arm:
	@GOOS=linux GOARCH=arm GOARM=5 go build \
		-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}" \
		-o ./bin/switchboard ./cmd
run:
	@go run ./cmd/main.go

test-watcher:
	@go test -v ./watcher

test-watcher-observe:
	@go test -v ./watcher -test.run TestObserve

test-event:
	@go test -v ./event

test-utils:
	@go test -v ./utils

test-cmd:
	@go test -v ./cmd

test-all: test-all test-watcher test-watcher-observe test-event test-utils test-cmd

lint-watcher:
	@golint ./watcher

lint-event:
	@golint ./event

lint-utils:
	@golint ./utils

lint-cmd:
	@golint ./cmd

lint-all: lint-watcher lint-event lint-utils lint-cmd

vet-watcher:
	@go vet ./watcher/watcher.go
	@go vet ./watcher/watcher_test.go

vet-event:
	@go vet ./event/event.go
	@go vet ./event/event_test.go

vet-utils:
	@go vet ./utils/utils.go
	@go vet ./utils/utils_test.go

vet-cmd:
	@go vet ./cmd/main.go
	@go vet ./cmd/watch.go

vet-all: vet-watcher vet-event vet-utils vet-cmd

test-format-all:
	@gofmt -l -d .
