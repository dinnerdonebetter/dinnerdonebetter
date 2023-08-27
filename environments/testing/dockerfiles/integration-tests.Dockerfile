# syntax=docker/dockerfile:1
FROM golang:1.21-bookworm

WORKDIR /go/src/github.com/dinnerdonebetter/backend
ENV SKIP_PASETO_TESTS=TRUE
COPY . .

# to debug a specific test:
ENTRYPOINT go test -parallel 1 -v -failfast github.com/dinnerdonebetter/backend/tests/integration -run TestIntegration/TestRecipes_Cloning

# ENTRYPOINT go test -v github.com/dinnerdonebetter/backend/tests/integration
