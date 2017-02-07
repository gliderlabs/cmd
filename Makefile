
dev:
	comlab dev

build:
	go build -a -o ./build/cmd ./cmd/cmd

test:
	go test -v $(shell glide nv)

image:
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

deploy:
	convox deploy -a alpha-cmd-io --wait

dynamodb:
	docker build -t dynamodb-local ./dev/dynamodb
	docker run -p 8000:8000 dynamodb-local -inMemory -sharedDb

www-build:
	make -C www build

www-dev:
	make -C www dev

.PHONY: dev image docker deploy dynamodb build www-build www-dev
