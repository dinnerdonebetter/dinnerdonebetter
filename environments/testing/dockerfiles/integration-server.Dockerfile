# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -trimpath -o /prixfixe -v gitlab.com/prixfixe/prixfixe/cmd/server/v1

# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm audit fix && npm run build

# final stage
FROM debian:stretch

RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
USER appuser

COPY environments/testing/config_files/integration-tests-postgres.toml /etc/config.toml
COPY --from=build-stage /prixfixe /prixfixe
COPY --from=frontend-build-stage /app/dist /frontend

ENTRYPOINT ["/prixfixe"]
