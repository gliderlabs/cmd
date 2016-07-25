
dev: build run

build:
	docker build -t progrium/cmd .

run:
	docker rm -f cmd || true
	docker run -d --name cmd \
		--publish 2222:22 \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		progrium/cmd

deploy: build
	docker push progrium/cmd
	ssh root@cmd.io -p 2222 docker pull progrium/cmd
	ssh root@cmd.io -p 2222 docker rm -f cmd
	ssh root@cmd.io -p 2222 docker run -d --name cmd \
		--volume /etc/ssh/ssh_host_rsa_key:/tmp/data/id_host \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--volume /var/run/cmd:/config \
		--restart always \
		--publish 22:22 \
		progrium/cmd
