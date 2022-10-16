# syntax=docker/dockerfile:1
FROM golang:1.19-buster

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY . .

RUN go build -o /meal_plan_task_creator github.com/prixfixeco/api_server/cmd/localdev/meal_plan_task_creator

ENTRYPOINT /meal_plan_task_creator