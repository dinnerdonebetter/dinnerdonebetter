# syntax=docker/dockerfile:1
FROM golang:1.22-bullseye

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN go build -o /server github.com/dinnerdonebetter/backend/cmdservices/api/http

ENTRYPOINT /server
