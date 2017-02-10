FROM alpine:3.5

RUN apk --no-cache add ca-certificates

COPY ./build/linux_amd64/cmd /usr/local/bin/cmd
COPY ./ui /app/ui
COPY ./com /app/com
WORKDIR /app

ENV LOCAL="false"
ENV CMD_LISTEN_ADDR=":22"
ENV CMD_CONFIG_DIR="/config"
ENV CMD_HOSTKEY_PEM="/tmp/data/id_host"
ENV WEB_LISTEN_ADDR=":80"
EXPOSE 22 80
CMD ["/usr/local/bin/cmd"]
