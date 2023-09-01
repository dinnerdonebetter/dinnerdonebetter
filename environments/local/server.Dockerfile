# syntax=docker/dockerfile:1
FROM golang:1.21-bullseye

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN go build -o /server github.com/dinnerdonebetter/backend/cmd/services/api/http

ENTRYPOINT /server
