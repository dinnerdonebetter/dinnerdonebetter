FROM golang:buster

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

COPY . .

ENTRYPOINT [ "go", "test", "-parallel", "1", "-v", "-failfast", "gitlab.com/prixfixe/prixfixe/tests/integration" ]

# to debug a specific test:
# ENTRYPOINT [ "go", "test", "-parallel", "1", "-v", "-failfast", "gitlab.com/prixfixe/prixfixe/tests/integration", "-run", "InsertTestNameHere" ]
