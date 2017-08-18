# This Dockerfile is not for use so much as to verify a working dev environment
# from scratch. Building this with `make test-env` will show missing dependencies
# and help clarify what software is necessary for a new development setup.

FROM alpine:3.5
RUN apk --update add go curl make bash git build-base ca-certificates glide
ENV GOPATH /usr/local
COPY . /usr/local/src/github.com/gliderlabs/comlab
WORKDIR /usr/local/src/github.com/gliderlabs/comlab
RUN make install && comlab
