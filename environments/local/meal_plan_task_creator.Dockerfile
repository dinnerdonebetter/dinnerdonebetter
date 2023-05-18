# syntax=docker/dockerfile:1
FROM golang:1.20-buster

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN go build -o /meal_plan_task_creator github.com/dinnerdonebetter/backend/cmd/localdev/meal_plan_task_creator

ENTRYPOINT /meal_plan_task_creator