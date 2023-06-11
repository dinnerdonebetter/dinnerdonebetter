# build stage
FROM golang:1.20-bullseye AS build-stage

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY cmd/functions/search_data_index_scheduler .

RUN go build -trimpath -o /action

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /action /action

ENTRYPOINT ["/action"]
