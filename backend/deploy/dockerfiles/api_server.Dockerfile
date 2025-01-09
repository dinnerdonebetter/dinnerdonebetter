# build stage
FROM golang:1.23-bullseye AS build-stage

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN go build -trimpath -o /server github.com/dinnerdonebetter/backend/cmd/services/api/http

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /server /server

ENTRYPOINT ["/server"]



echo -e "GET /_meta_/ready HTTP/1.1\r\nHost: dinner-done-better-api-svc\r\nConnection: close\r\n\r\n" | nc dinner-done-better-api-svc 8000
