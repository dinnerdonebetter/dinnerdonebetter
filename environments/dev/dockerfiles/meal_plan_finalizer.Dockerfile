# build stage
FROM golang:1.17-stretch AS build-stage

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY . .

RUN go build -trimpath -o /finalizer -v github.com/prixfixeco/api_server/cmd/meal_plan_finalizer

# final stage
FROM debian:stretch

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=build-stage /finalizer /finalizer

ENTRYPOINT ["/finalizer"]
