# syntax=docker/dockerfile:1
FROM golang:1.20-buster AS build-stage

WORKDIR /go/src/github.com/prixfixeco/backend

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY vendor vendor
COPY go.mod go.mod
COPY go.sum go.sum

RUN --mount=type=cache,target=/root/.cache/go-build go build -trimpath -o /prixfixe github.com/prixfixeco/backend/cmd/server

# final stage
FROM debian:stretch

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=build-stage /prixfixe /prixfixe

ENTRYPOINT ["/prixfixe"]
