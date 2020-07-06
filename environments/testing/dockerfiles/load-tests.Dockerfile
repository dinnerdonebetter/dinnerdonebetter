# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -o /loadtester gitlab.com/prixfixe/prixfixe/tests/v1/load

# final stage
FROM debian:stable

COPY --from=build-stage /loadtester /loadtester

ENV DOCKER=true

ENTRYPOINT ["/loadtester"]
