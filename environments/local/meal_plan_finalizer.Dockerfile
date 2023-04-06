# syntax=docker/dockerfile:1
FROM golang:1.20-buster

WORKDIR /go/src/github.com/prixfixeco/backend

COPY . .

RUN go build -o /meal_plan_finalizer github.com/prixfixeco/backend/cmd/localdev/meal_plan_finalizer

ENTRYPOINT /meal_plan_finalizer