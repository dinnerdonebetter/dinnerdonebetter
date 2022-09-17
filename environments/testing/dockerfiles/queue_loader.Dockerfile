# build stage
FROM golang:1.19-bullseye

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY . .

RUN go build -trimpath -o /queue_loader github.com/prixfixeco/api_server/cmd/localdev/queue_loader

ENTRYPOINT /queue_loader