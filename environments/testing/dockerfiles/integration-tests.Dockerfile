FROM golang:1.17-stretch


RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates make git gcc musl-dev

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY . .

ENV SKIP_PASETO_TESTS=TRUE

# to debug a specific test:
# ENTRYPOINT [ "go", "test", "-parallel", "1", "-v", "-failfast", "github.com/prixfixeco/api_server/tests/integration", "-run", "TestIntegration/TestMealPlanOptions_Listing" ]

ENTRYPOINT [ "go", "test", "-v", "-failfast", "github.com/prixfixeco/api_server/tests/integration" ]
