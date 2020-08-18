# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go test -o /integration-server -c -coverpkg \
	gitlab.com/prixfixe/prixfixe/internal/..., \
	gitlab.com/prixfixe/prixfixe/database/v1/..., \
	gitlab.com/prixfixe/prixfixe/services/v1/..., \
	gitlab.com/prixfixe/prixfixe/cmd/server/v1/ \
    gitlab.com/prixfixe/prixfixe/cmd/server/v1

# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# final stage
FROM debian:stable

COPY --from=build-stage /integration-server /integration-server
COPY --from=frontend-build-stage /app/dist /frontend

EXPOSE 80

ENTRYPOINT ["/integration-server", "-test.coverprofile=/home/integration-coverage.out"]

