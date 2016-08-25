
build:
	docker build -t progrium/cmd-example .

run: build
	docker run --rm progrium/cmd-example arg1 arg2

push: build
	docker push progrium/cmd-example
