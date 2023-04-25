# build stage
FROM golang:1.20-bullseye AS build-stage

WORKDIR /go/src/github.com/prixfixeco/backend

COPY . .

RUN go build -trimpath -o /server github.com/prixfixeco/backend/cmd/server

# final stage
FROM debian:bullseye

# RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /server /server

ENTRYPOINT ["/server"]
