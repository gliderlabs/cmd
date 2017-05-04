## PHONY

.PHONY: help setup clean clobber dev image docker deploy build build-all www-build www-dev test test-env test-go test-all deps-update deps-go services

## VARIABLES

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
VERSION := $(shell git rev-parse HEAD)

os = $(shell echo $(1) | cut -d"_" -f1)
arch = $(shell echo $(1) | cut -d"_" -f2)

## COMMANDS

help: ## show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: build/$(GOOS)_$(GOARCH)/cmd ## build for current system

build-all: build/linux_amd64/cmd build/darwin_amd64/cmd ## build for all systems

build/%/cmd:
	GOOS=$(call os, $*) GOARCH=$(call arch, $*) \
    go build -ldflags "-X main.Version=$(VERSION)" -o $@ ./cmd/cmd

dev: services deps-update ## run dev harness
	comlab dev

services: ## run backing services
	docker-compose -f dev/services.yml up -d

setup: deps-update ## setup development environment
	go get -u github.com/gliderlabs/comlab/...

clean: ## delete typical build artifacts
	-rm -rf build

clobber: clean ## reset dev environment
	docker-compose -f dev/services.yml down
	-rm -rf vendor
	-rm -f .git/deps-*

image: build/linux_amd64/cmd ## build docker image
	docker build -t gliderlabs/cmd .

deploy-alpha: build/linux_amd64/cmd ## deploy to alpha channel
	sigil -p -f run/channels/alpha.yaml image=$(IMAGE) | kubectl apply --namespace cmd -f -
	kubectl rollout status deployment/cmd-alpha --namespace cmd --watch

deploy-beta: build/linux_amd64/cmd ## deploy to beta channel
	sigil -p -f run/channels/beta.yaml image=$(IMAGE) | kubectl apply --namespace cmd -f -
	kubectl rollout status deployment/cmd-beta --namespace cmd --watch

## TESTS

test: test-go ## run common tests

test-all: test-go test-env ## run ALL tests

test-go: ## run golang tests
	go test -v $(shell glide nv)

test-env: services ## test dev environment
	docker build -t cmd-env -f dev/setup/Dockerfile .
	docker run --rm \
		--volume /var/run/docker.sock:/var/run/docker.sock:ro \
		--volume $(shell pwd)/.env:/usr/local/src/github.com/gliderlabs/cmd/.env:ro \
		--env DYNAMODB_ENDPOINT=http://$(shell docker inspect dev_dynamodb_1 --format '{{ .NetworkSettings.Networks.dev_default.IPAddress }}'):8000 \
		cmd-env
	-docker rmi cmd-env

## WRAPPERS

ui-build:
	$(MAKE) -C ui build

www-build:
	$(MAKE) -C www build

www-dev:
	$(MAKE) -C www dev

## DEPENDENCIES

deps-update: ## update dependencies if changed
	./dev/deps.sh

deps-go:
	glide install -v
	git log -n 1 --pretty=format:%h -- glide.yaml > .git/deps-go
