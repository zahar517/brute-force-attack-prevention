SERVICE := "./bin/service"
CLI := "./bin/cli"
MIGRATE := "./bin/migrate"
DOCKER_IMG="bruteforce:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(SERVICE) -ldflags "$(LDFLAGS)" ./cmd/service
	go build -v -o $(CLI) -ldflags "$(LDFLAGS)" ./cmd/cli
	go build -v -o $(MIGRATE) -ldflags "$(LDFLAGS)" ./cmd/migrate

run: build
	$(SERVICE)

migrate: build
	$(MIGRATE) up

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(SERVICE) version

generate:
	go generate ./...

test:
	go test -race -count 100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.45.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
