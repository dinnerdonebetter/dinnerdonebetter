# syntax=docker/dockerfile:1
FROM golang:1.20-buster

WORKDIR /go/src/github.com/prixfixeco/backend
ENV SKIP_PASETO_TESTS=TRUE
COPY . .

# to debug a specific test:
# ENTRYPOINT go test -parallel 1 -v -failfast github.com/prixfixeco/backend/tests/integration -run TestIntegration/TestUsers_Reading

ENTRYPOINT go test -v github.com/prixfixeco/backend/tests/integration
