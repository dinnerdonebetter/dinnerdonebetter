# syntax=docker/dockerfile:1
FROM golang:1.24-bullseye

WORKDIR /go/src/github.com/dinnerdonebetter/backend
COPY . .

# to debug a specific test:
# ENTRYPOINT go test -parallel 1 -v -failfast github.com/dinnerdonebetter/backend/tests/integration -run TestIntegration/TestValidPreparationInstruments_CompleteLifecycle

ENTRYPOINT ["go", "test", "-v", "github.com/dinnerdonebetter/backend/tests/integration"]
