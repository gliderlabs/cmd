FROM alpine
RUN apk --update add bash curl go git mercurial docker

RUN curl -Ls https://github.com/gliderlabs/sshfront/releases/download/v0.2.1/sshfront_0.2.1_Linux_x86_64.tgz \
    | tar -zxC /bin

COPY ./data /tmp/data

ENV GOPATH /go
COPY . /go/src/github.com/progrium/cmd
WORKDIR /go/src/github.com/progrium/cmd
RUN go get && CGO_ENABLED=0 go build -a -installsuffix cgo -o /bin/cmd \
  && ln -s /bin/cmd /bin/auth

EXPOSE 22
ENTRYPOINT ["/bin/bash", "-c", "/bin/sshfront -e -k /tmp/data/id_host -a /bin/auth /bin/cmd"]
