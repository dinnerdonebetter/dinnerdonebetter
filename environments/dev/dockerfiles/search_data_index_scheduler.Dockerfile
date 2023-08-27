# build stage
FROM golang:1.21-bullseye AS build-stage

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

ENV CGO_ENABLED=0

RUN go build -trimpath -o /action github.com/dinnerdonebetter/backend/cmd/jobs/search_data_index_scheduler

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /action /action

ENTRYPOINT ["/action"]
