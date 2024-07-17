# syntax=docker/dockerfile:1
FROM golang:1.22-bullseye

WORKDIR /go/src/github.com/dinnerdonebetter/backend/backend

COPY . .

RUN go build -o /server github.com/dinnerdonebetter/backend/backend/cmd/services/api/http

ENTRYPOINT /server
