# build stage
FROM golang:1.24-bullseye AS build-stage

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN go install github.com/go-delve/delve/cmd/dlv@v1.24.0
RUN go build -trimpath -o /server github.com/dinnerdonebetter/backend/cmd/services/api/http

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /server /server
# COPY --from=build-stage /go/bin/dlv /dlv

# ENTRYPOINT ["/dlv", "exec", "--allow-non-terminal-interactive=true", "/server"]
ENTRYPOINT ["/server"]
