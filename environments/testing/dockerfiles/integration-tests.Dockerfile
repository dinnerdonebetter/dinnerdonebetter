FROM golang:1.17-stretch

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

COPY . .

# to debug a specific test:
ENTRYPOINT [ "go", "test", "-parallel", "1", "-v", "-failfast", "gitlab.com/prixfixe/prixfixe/tests/integration", "-run", "TestIntegration/TestRecipeStepIngredients" ]

# ENTRYPOINT [ "go", "test", "-v", "-failfast", "gitlab.com/prixfixe/prixfixe/tests/integration" ]
