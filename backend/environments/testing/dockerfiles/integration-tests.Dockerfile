# syntax=docker/dockerfile:1
FROM golang:1.23-bullseye

WORKDIR /go/src/github.com/dinnerdonebetter/backend
COPY . .

# to debug a specific test:
ENTRYPOINT go test -parallel 1 -v -failfast github.com/dinnerdonebetter/backend/tests/integration -run TestIntegration/TestRecipes_Searching

# ENTRYPOINT go test github.com/dinnerdonebetter/backend/tests/integration
