# build stage
FROM golang:buster AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY . .

RUN go build -trimpath -o /prixfixe -v gitlab.com/prixfixe/prixfixe/cmd/server

# final stage
FROM debian:bullseye-slim

COPY --from=build-stage /prixfixe /prixfixe

ENTRYPOINT ["/prixfixe"]
