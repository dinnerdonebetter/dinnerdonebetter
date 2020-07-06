FROM golang:stretch

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD ../../../dockerfiles .

CMD if [ $(gofmt -l . | grep -Ev '^vendor\/' | head -c1 | wc -c) -ne 0 ]; then exit 1; fi
