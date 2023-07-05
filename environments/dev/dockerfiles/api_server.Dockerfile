# build stage
FROM golang:1.20-bullseye AS build-stage

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN GOOS=js GOARCH=wasm go build -o internal/services/wasm/assets/helpers.wasm github.com/dinnerdonebetter/backend/cmd/wasm
RUN go build -trimpath -o /server github.com/dinnerdonebetter/backend/cmd/services/api/http

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /server /server

ENTRYPOINT ["/server"]
