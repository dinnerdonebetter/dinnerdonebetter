# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY . .

# we need the `-tags json1` so sqlite can support JSON columns.
RUN go build -tags json1 -trimpath -o /prixfixe -v gitlab.com/prixfixe/prixfixe/cmd/server

# final stage
FROM debian:stretch

COPY --from=build-stage /prixfixe /prixfixe

ENTRYPOINT ["/prixfixe"]
