# syntax=docker/dockerfile:1
FROM golang:1.21-buster

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN go build -o /queue_loader github.com/dinnerdonebetter/backend/cmd/localdev/queue_loader

ENTRYPOINT /queue_loader