# build stage
FROM golang:1.19-bullseye

WORKDIR /go/src/github.com/prixfixeco/api_server

ENTRYPOINT go run github.com/prixfixeco/api_server/cmd/localdev/advanced_prep_step_creator
