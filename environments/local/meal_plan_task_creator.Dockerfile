# syntax=docker/dockerfile:1
FROM golang:1.19-buster

WORKDIR /go/src/github.com/prixfixeco/backend

COPY . .

RUN go build -o /meal_plan_task_creator github.com/prixfixeco/backend/cmd/localdev/meal_plan_task_creator

ENTRYPOINT /meal_plan_task_creator