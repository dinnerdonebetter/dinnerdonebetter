# build stage
FROM golang:1.17-stretch AS build-stage

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY . .

RUN go build -trimpath -o /prixfixe -v github.com/prixfixeco/api_server/cmd/server

# final stage
FROM debian:stretch

COPY --from=build-stage /prixfixe /prixfixe

RUN mkdir /home/appuser
RUN groupadd --gid 999 appuser && \
    useradd --system --uid 999 --gid appuser appuser
RUN chown appuser /home/appuser
WORKDIR /home/appuser
USER appuser

COPY environments/testing/config_files/frontend-tests.toml /etc/config.toml

ENTRYPOINT ["/prixfixe"]
