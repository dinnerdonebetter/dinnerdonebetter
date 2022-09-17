# build stage
FROM golang:1.19-bullseye

WORKDIR /go/src/github.com/prixfixeco/api_server

RUN go build -trimpath -o /meal_plan_finalizer github.com/prixfixeco/api_server/cmd/localdev/meal_plan_finalizer

ENTRYPOINT /meal_plan_finalizer