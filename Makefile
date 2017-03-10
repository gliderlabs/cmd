GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
VERSION := $(shell git rev-parse HEAD)

build: build/$(GOOS)_$(GOARCH)/cmd

build-all: build/linux_amd64/cmd build/darwin_amd64/cmd

os = $(shell echo $(1) | cut -d"_" -f1)
arch = $(shell echo $(1) | cut -d"_" -f2)

build/%/cmd:
	GOOS=$(call os, $*) GOARCH=$(call arch, $*) \
    go build -ldflags "-X main.Version=$(VERSION)" -o $@ ./cmd/cmd

dev:
	comlab dev

clean:
	-rm -rf build

test:
	go test -v $(shell glide nv)

image: build/linux_amd64/cmd
	docker build -t progrium/cmd .

image-dev:
	docker build -t progrium/cmd-dev -f Dockerfile.dev .

docker: image
	@docker rm -f cmd || true
	docker run -d --name cmd \
		--publish 2222:22 \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--volume $(shell pwd)/com/cmd/data/dev_host:/tmp/data/id_host \
		--volume $(shell pwd)/local:/config \
		progrium/cmd

deploy: build/linux_amd64/cmd
	convox deploy -a alpha-cmd-io --wait

deploy-beta: build/linux_amd64/cmd
	sigil -p -f run/channels/beta.yaml build=$(BUILD_NUM) | kubectl apply -f -
	kubectl rollout status deployment/beta-cmd-io --watch

dynamodb:
	docker build -t dynamodb-local ./dev/dynamodb
	docker run -p 8000:8000 dynamodb-local -inMemory -sharedDb

ui-build:
	$(MAKE) -C ui build

www-build:
	$(MAKE) -C www build

www-dev:
	$(MAKE) -C www dev

.PHONY: dev image docker deploy dynamodb build build-all www-build www-dev
