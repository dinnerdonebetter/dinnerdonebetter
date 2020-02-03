FROM golang:stretch

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

ENTRYPOINT [ "go", "test", "-v", "-failfast", "-parallel=1", "gitlab.com/prixfixe/prixfixe/tests/v1/frontend" ]
