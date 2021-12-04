# build stage
FROM golang:1.17-stretch AS build-stage

WORKDIR /go/src/github.com/prixfixeco/api_server

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY . .

RUN go build -trimpath -o /workers -v github.com/prixfixeco/api_server/cmd/workers/localdev

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=build-stage /workers /workers

ENTRYPOINT ["/workers"]
