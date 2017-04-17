---
date: 2017-01-31T18:00:00-06:00
title: Authoring Commands
weight: 30
---

You'll probably be using your own commands more than off-the-shelf commands.
Right now, since Docker Hub is the only source for commands, making and
publishing commands is as easy as any Docker container. After you've gone
through this once, you might be surprised at how quickly you can make and update
Cmd.io commands. It's literally build, push, use.

Up to this point, we haven't needed to use or install Docker. For the time
being, you'll need Docker to make Cmd.io commands. We highly recommend [Docker
for Mac](https://docs.docker.com/docker-for-mac/) if you're running macOS.

You'll also need a [Docker Hub](https://hub.docker.com/) account and be sure to
login with Docker (`docker login`).

### Commands based on existing utilities

The recommended way to build commands from existing open source utilities is to
install them via a package manager. To keep your experience snappy, and since
Cmd.io may enforce a size limit on images, we highly encourage you to use Alpine
Linux for all command containers.

Alpine combines the small size of Busybox  (~5MB) with a large package index
optimized for small disk footprints. You can search for packages [based on
name](http://pkgs.alpinelinux.org/packages) or [based on
contents](http://pkgs.alpinelinux.org/contents). If you can't find a package for
a utility, you can try using `ubuntu-debootstrap`, which is a minimal Ubuntu
image with `apt-get`. However it starts at ~90MB and easily bloats from there.

[Here is jq](http://pkgs.alpinelinux.org/package/v3.4/main/x86_64/jq) for Alpine
v3.4, so we can make a container for it with a simple Dockerfile that uses the
`apk` package tool:

```
FROM alpine:3.4
RUN apk add --update --no-cache jq
ENTRYPOINT ["/usr/bin/jq"]
```

The directives used in this example are nearly all that make sense to use for
Cmd.io commands, but here is the full [Dockerfile
reference](https://docs.docker.com/engine/reference/builder/).

Now we can build this with Docker, assuming we're in the directory with the
Dockerfile. Immediately after building, we can push to Docker Hub. Replace
`progrium` with your Docker ID.

```
$ docker build -t progrium/jq .
...
$ docker push progrium/jq
...
```

At this point you can now install this command on Cmd.io like before. If you
push new versions of the image to Docker Hub, Cmd.io will pull it just before
the next run.

### Commands based on scripts

Making a container for a script is not that different from making it for an
existing utility. You'll want to install the interpreter and any other utilities
the script depends on the same way as before. But you'll also be adding your
script and making it the entrypoint.

Create a file called `netpoll` and make sure it's executable with `chmod +x
netpoll`. Inside it, put this Bash script that uses netcat to poll an address
and port for roughly 10 seconds or until the port accepts a connection. If it
connects it returns. If it times out it returns non-zero. A rather handy little
script.

```
#!/bin/bash
for retry in $(seq 1 ${TIMEOUT:-10}); do
  nc -z -w 1 "$1" "$2" && break
done
```

In the same directory create a `Dockerfile` like this:

```
FROM alpine:3.4
RUN apk add --update bash netcat-openbsd
COPY ./netpoll /bin/netpoll
ENTRYPOINT ["/bin/netpoll"]
```

Build, push, and install with Cmd.io. Let's say I installed it as `netpoll`. I
can run it against `cmd.io` port `22` and it returns immediately with status 0.
Run against `cmd.io` port `23` and it blocks for at least 10 seconds before
giving up and returning status 1.

```
$ ssh alpha.cmd.io netpoll cmd.io 22; echo $?
0
$ ssh alpha.cmd.io netpoll cmd.io 23; echo $?
1
```
