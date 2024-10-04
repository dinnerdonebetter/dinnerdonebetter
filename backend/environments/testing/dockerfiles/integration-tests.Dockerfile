# syntax=docker/dockerfile:1
FROM golang:1.22-bullseye

WORKDIR /go/src/github.com/dinnerdonebetter/backend
COPY . .

# TestIntegration/TestHouseholds_ChangingMemberships
# TestIntegration/TestHouseholds_OwnershipTransfer

# to debug a specific test:
ENTRYPOINT go test -parallel 1 -v -failfast github.com/dinnerdonebetter/backend/tests/integration -run TestIntegration/TestHouseholds_ChangingMemberships

# ENTRYPOINT go test -v github.com/dinnerdonebetter/backend/tests/integration
