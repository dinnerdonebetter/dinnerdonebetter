# syntax=docker/dockerfile:1
FROM golang:1.20-buster

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN go build -o /server github.com/dinnerdonebetter/backend/cmd/server/http

ENTRYPOINT /server
