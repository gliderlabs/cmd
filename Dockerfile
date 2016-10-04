FROM alpine
RUN apk --update add go curl make git docker
ENV GOPATH /usr/local
RUN curl https://glide.sh/get | sh
COPY . /usr/local/src/github.com/progrium/cmd
WORKDIR /usr/local/src/github.com/progrium/cmd
RUN glide install
RUN go install ./cmd/cmd
ENV CMD_LISTEN_ADDR=":22"
ENV CMD_CONFIG_DIR="/config"
ENV CMD_HOSTKEY_PEM="/tmp/data/id_host"
EXPOSE 22
CMD ["/usr/local/bin/cmd"]
