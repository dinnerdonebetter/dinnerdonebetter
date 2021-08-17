# build stage
FROM golang:buster AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

COPY . .

RUN go build -trimpath -o /prixfixe -v gitlab.com/prixfixe/prixfixe/cmd/server

# frontend-build-stage
FROM node:lts AS frontend-build-stage

WORKDIR /app

COPY frontend .

RUN npm install && npm audit fix && npm run build

# final stage
FROM debian:bullseye-slim

COPY --from=build-stage /prixfixe /prixfixe

RUN mkdir /home/appuser
RUN groupadd --gid 999 appuser && \
    useradd --system --uid 999 --gid appuser appuser
RUN chown appuser /home/appuser
WORKDIR /home/appuser
USER appuser

COPY environments/testing/config_files/frontend-tests.toml /etc/config.toml
COPY --from=frontend-build-stage /app/dist /frontend

ENTRYPOINT ["/prixfixe"]
