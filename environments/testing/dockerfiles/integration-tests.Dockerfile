# syntax=docker/dockerfile:1
FROM golang:1.19-buster

WORKDIR /go/src/github.com/prixfixeco/api_server
ENV SKIP_PASETO_TESTS=FALSE
COPY . .

# to debug a specific test:
# ENTRYPOINT go test -parallel 1 -v -failfast github.com/prixfixeco/api_server/tests/integration -run TestIntegration/TestWebhooks_Reading_Returns404ForNonexistentWebhook

ENTRYPOINT go test -v github.com/prixfixeco/api_server/tests/integration
