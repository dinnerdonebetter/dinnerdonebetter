# syntax=docker/dockerfile:1
FROM golang:1.19-buster

WORKDIR /go/src/github.com/prixfixeco/backend

COPY . .

RUN go build -o /queue_loader github.com/prixfixeco/backend/cmd/localdev/queue_loader

ENTRYPOINT /queue_loader