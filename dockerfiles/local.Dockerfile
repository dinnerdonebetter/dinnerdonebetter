## NOTE: This dockerfile doesn't contain a build phase for the frontend.
##       The reason for this is so you run the frontend build phase in a different process

# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

COPY . .

RUN go build -trimpath -o /prixfixe gitlab.com/prixfixe/prixfixe/cmd/server/v1

# final stage
FROM debian:stretch

COPY --from=build-stage /prixfixe /prixfixe
COPY config_files config_files

RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
USER appuser

ENV DOCKER=true

ENTRYPOINT ["/prixfixe"]
