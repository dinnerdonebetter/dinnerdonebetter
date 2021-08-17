# build stage
FROM golang:buster AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY . .

RUN go build -o /loadtester gitlab.com/prixfixe/prixfixe/tests/load

# final stage
FROM debian:bullseye-slim

COPY --from=build-stage /loadtester /loadtester

ENV DOCKER=true

ENTRYPOINT ["/loadtester"]
