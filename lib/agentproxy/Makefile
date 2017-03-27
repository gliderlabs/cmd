.PHONY: push image

push: image
	docker push gliderlabs/ssh-agent-proxy

image: build/agentproxy
	docker build -t gliderlabs/ssh-agent-proxy .

build/agentproxy:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o build/agentproxy ./cmd

clean:
	rm -rf build
