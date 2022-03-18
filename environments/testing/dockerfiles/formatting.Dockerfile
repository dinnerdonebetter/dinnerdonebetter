FROM golang:1.18-stretch

WORKDIR /go/src/github.com/prixfixeco/api_server

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY . .

CMD if [ $(gofmt -l . | grep -Ev '^vendor\/' | head -c1 | wc -c) -ne 0 ]; then exit 1; fi
