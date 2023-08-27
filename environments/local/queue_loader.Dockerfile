# syntax=docker/dockerfile:1
FROM golang:1.21-bookworm

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

ENV CGO_ENABLED=0

RUN go build -o /queue_loader github.com/dinnerdonebetter/backend/cmd/localdev/queue_loader

ENTRYPOINT /queue_loader