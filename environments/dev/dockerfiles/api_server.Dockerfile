# build stage
FROM golang:1.20-bullseye AS build-stage

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build go build -trimpath -o /server github.com/dinnerdonebetter/backend/cmd/server/http

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /server /server

ENTRYPOINT ["/server"]
