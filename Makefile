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
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t $(DOCKER_IMG) -f build/Dockerfile .

up: build-img
	docker-compose -f ./build/docker-compose.yml up

down:
	docker-compose -f ./build/docker-compose.yml down

restart: down up

version: build
	$(SERVICE) version

generate:
	go generate ./...

test:
	go test -race -count 100 ./internal/...

integration-test:
	set -e ;\
	docker-compose -f ./build/docker-compose.yml -f ./build/docker-compose.test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f ./build/docker-compose.yml -f ./build/docker-compose.test.yml run test go test -v -tags=integration /go/src/tests/integration || test_status_code=$$? ;\
	docker-compose -f ./build/docker-compose.yml -f ./build/docker-compose.test.yml down ;\
	exit $$test_status_code ;

test-cleanup:
	docker-compose -f ./build/docker-compose.yml -f ./build/docker-compose.test.yml down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.45.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
