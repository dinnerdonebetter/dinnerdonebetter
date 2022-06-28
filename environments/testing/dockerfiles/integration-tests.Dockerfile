FROM golang:1.18-stretch

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates make git gcc musl-dev

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY tests tests
COPY vendor vendor
COPY go.mod go.mod
COPY go.sum go.sum

ENV SKIP_PASETO_TESTS=FALSE

# to debug a specific test:
# ENTRYPOINT [ "go", "test", "-parallel", "1", "-v", "-failfast", "github.com/prixfixeco/api_server/tests/integration", "-run", "TestIntegration/TestMealPlanOptionVotes_CompleteLifecycle" ]

ENTRYPOINT [ "go", "test", "-parallel", "1", "-v", "-failfast", "github.com/prixfixeco/api_server/tests/integration" ]
