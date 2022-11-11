# syntax=docker/dockerfile:1
FROM golang:1.19-buster

WORKDIR /go/src/github.com/prixfixeco/backend
ENV SKIP_PASETO_TESTS=FALSE
COPY . .

# to debug a specific test:
# ENTRYPOINT go test -parallel 1 -v -failfast github.com/prixfixeco/backend/tests/integration -run TestIntegration/TestMealPlans_CompleteLifecycleForSomeVotesReceived

ENTRYPOINT go test -v github.com/prixfixeco/backend/tests/integration
