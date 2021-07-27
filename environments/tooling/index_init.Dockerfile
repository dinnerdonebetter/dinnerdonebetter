# build stage
FROM golang:stretch as build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

COPY . .

RUN go build -trimpath -o /index_initializer gitlab.com/prixfixe/prixfixe/cmd/tools/index_initializer

# final stage
FROM debian:stable

COPY --from=build-stage /index_initializer /index_initializer

ENTRYPOINT ["/index_initializer"]
