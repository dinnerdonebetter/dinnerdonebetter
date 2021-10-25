# build stage
FROM golang:1.17-stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY . .

RUN go build -trimpath -o /prixfixe -v gitlab.com/prixfixe/prixfixe/cmd/server

# final stage
FROM debian:stretch

COPY --from=build-stage /prixfixe /prixfixe

ENTRYPOINT ["/prixfixe"]
