# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

COPY . .
COPY --from=frontend-build-stage /app/public /frontend

RUN go build -trimpath -o /prixfixe gitlab.com/prixfixe/prixfixe/cmd/server/v1

# final stage
FROM debian:stretch

COPY --from=build-stage /prixfixe /prixfixe

RUN mkdir /home/appuser
RUN groupadd --gid 999 appuser && \
    useradd --system --uid 999 --gid appuser appuser
RUN chown appuser /home/appuser
WORKDIR /home/appuser
USER appuser

ENV DOCKER=true

ENTRYPOINT ["/prixfixe"]