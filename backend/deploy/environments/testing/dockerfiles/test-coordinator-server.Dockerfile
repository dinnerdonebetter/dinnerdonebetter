# syntax=docker/dockerfile:1
FROM golang:1.23-bullseye AS build-stage

WORKDIR /go/src/github.com/dinnerdonebetter/backend

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY vendor vendor
COPY go.mod go.mod
COPY go.sum go.sum

RUN --mount=type=cache,target=/root/.cache/go-build go build -trimpath -o /coordination-server github.com/dinnerdonebetter/backend/cmd/tools/test_coordination_server

# final stage
FROM debian:bullseye

COPY --from=build-stage /coordination-server /coordination-server

ENTRYPOINT ["/coordination-server"]
