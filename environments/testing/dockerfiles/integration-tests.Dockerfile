FROM golang:stretch

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

ADD . .

ENTRYPOINT [ "go", "test", "-v", "-failfast", "gitlab.com/prixfixe/prixfixe/tests/v1/integration" ]

# for a more specific test:
# ENTRYPOINT [ "go", "test", "-v", "gitlab.com/prixfixe/prixfixe/tests/v1/integration", "-run", "InsertTestNameHere" ]
