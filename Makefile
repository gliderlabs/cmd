
dev:
	comlab dev

build:
	go build -a -o ./build/cmd ./cmd/cmd

image:
	docker build -t progrium/cmd .

docker: image
	@docker rm -f cmd || true
	docker run -d --name cmd \
		--publish 2222:22 \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--volume $(shell pwd)/com/cmd/data/dev_host:/tmp/data/id_host \
		--volume $(shell pwd)/local:/config \
		progrium/cmd

deploy: image
	docker save progrium/cmd | ssh root@cmd.io -p 2222 docker load
	@ssh root@cmd.io -p 2222 docker rm -f cmd || true
	ssh root@cmd.io -p 2222 docker run -d --name cmd \
	  --env WEB_LISTEN_ADDR=:8080 \
		--volume /etc/ssh/ssh_host_rsa_key:/tmp/data/id_host \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--volume /var/run/cmd:/config \
		--restart always \
		--publish 22:22 \
		--publish 80:8080 \
		progrium/cmd

dynamodb:
	docker build -t dynamodb-local ./dev/dynamodb
	docker run -p 8000:8000 dynamodb-local -inMemory -sharedDb

 .PHONY: dev image docker deploy dynamodb build
