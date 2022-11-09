# syntax=docker/dockerfile:1
FROM golang:1.19-buster

WORKDIR /go/src/github.com/prixfixeco/backend

COPY . .

RUN go build -o /server github.com/prixfixeco/backend/cmd/server

ENTRYPOINT /server
