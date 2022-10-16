# syntax=docker/dockerfile:1
FROM golang:1.19-buster

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY vendor vendor
COPY go.mod go.mod
COPY go.sum go.sum

RUN --mount=type=cache,target=/root/.cache/go-build go build -trimpath -o /queue_loader github.com/prixfixeco/api_server/cmd/localdev/queue_loader

ENTRYPOINT /queue_loader