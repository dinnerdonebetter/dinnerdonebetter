# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -trimpath -o /prixfixe gitlab.com/prixfixe/prixfixe/cmd/server/v1

# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# final stage
FROM debian:stable

RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
USER appuser

ENV CONFIGURATION_FILEPATH=/etc/prixfixe/config.toml

COPY config_files/development.toml /etc/prixfixe/config.toml
COPY --from=build-stage /prixfixe /prixfixe
COPY --from=frontend-build-stage /app/dist /frontend


ENTRYPOINT ["/prixfixe"]
