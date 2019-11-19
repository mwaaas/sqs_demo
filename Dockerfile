# syntax = docker/dockerfile:experimental
FROM golang:1.12-alpine3.10 as build-env

WORKDIR /usr/src/app

RUN apk add git

COPY . .

# based on this
# https://github.com/golang/go/issues/27719#issuecomment-514747274
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build

#FROM alpine:3.10
#
#COPY --from=build-env /usr/src/app/sqs_demo /usr/local/bin/
#
#CMD ["sqs_demo"]