# build stage
FROM golang:1.19-buster

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY . .

RUN go build -o /server github.com/prixfixeco/api_server/cmd/server

ENTRYPOINT /server
