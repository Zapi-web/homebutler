VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)

.PHONY: build clean test build-web build-all

build:
	go build -ldflags "$(LDFLAGS)" -o homebutler .

build-web:
	cd web && npm install && npm run build
	rm -rf internal/server/web_dist/*
	cp -r web/dist/* internal/server/web_dist/

build-all: build-web build

clean:
	rm -f homebutler
	rm -rf web/dist web/node_modules

test:
	go test ./...
