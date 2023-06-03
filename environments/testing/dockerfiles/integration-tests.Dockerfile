# syntax=docker/dockerfile:1
FROM golang:1.20-buster

WORKDIR /go/src/github.com/dinnerdonebetter/backend
ENV SKIP_PASETO_TESTS=TRUE
COPY . .

# to debug a specific test:
# ENTRYPOINT go test -parallel 1 -v -failfast github.com/dinnerdonebetter/backend/tests/integration -run TestIntegration/TestRecipes_Realistic

ENTRYPOINT go test -v github.com/dinnerdonebetter/backend/tests/integration
