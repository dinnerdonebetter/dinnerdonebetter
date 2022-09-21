# build stage
FROM golang:1.19-buster AS build-stage

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY . .

RUN go build -trimpath -o /server github.com/prixfixeco/api_server/cmd/server

# final stage
FROM debian:stretch

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=build-stage /server /server

ENTRYPOINT ["/server"]
