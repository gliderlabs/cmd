FROM alpine:3.5

RUN apk --no-cache add ca-certificates

COPY ./build/linux_amd64/cmd /usr/local/bin/cmd
COPY ./ui /app/ui
COPY ./com /app/com
WORKDIR /app

ENV LOCAL="false"
EXPOSE 22 80
CMD ["/usr/local/bin/cmd"]
