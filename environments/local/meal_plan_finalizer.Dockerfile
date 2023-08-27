# syntax=docker/dockerfile:1
FROM golang:1.21-buster

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

ENV CGO_ENABLED=0

RUN go build -o /meal_plan_finalizer github.com/dinnerdonebetter/backend/cmd/localdev/meal_plan_finalizer

ENTRYPOINT /meal_plan_finalizer