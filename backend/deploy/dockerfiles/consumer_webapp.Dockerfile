# build stage
FROM golang:1.26-trixie AS build-stage

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN ./scripts/build.sh -o /server github.com/dinnerdonebetter/backend/cmd/services/consumer

# final stage - use bookworm to match glibc requirements from golang:trixie build
FROM debian:bookworm

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /server /server

EXPOSE 80

ENTRYPOINT ["/server"]
