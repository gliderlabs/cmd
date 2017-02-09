FROM alpine:3.5

ENV GOPATH /go
ENV GOBIN /usr/local/bin

COPY . /go/src/github.com/progrium/cmd
WORKDIR /go/src/github.com/progrium/cmd

RUN apk --no-cache add go git glide build-base ca-certificates \
  && git config --global http.https://gopkg.in.followRedirects true \
  && glide install --strip-vendor \
  && go install ./cmd/cmd \
  && glide cc && rm -r ./vendor \
  && apk --no-cache del go git glide build-base


ENV LOCAL="false"
ENV CMD_LISTEN_ADDR=":22"
ENV CMD_CONFIG_DIR="/config"
ENV CMD_HOSTKEY_PEM="/tmp/data/id_host"
ENV WEB_LISTEN_ADDR=":80"
EXPOSE 22 80
CMD ["/usr/local/bin/cmd"]
