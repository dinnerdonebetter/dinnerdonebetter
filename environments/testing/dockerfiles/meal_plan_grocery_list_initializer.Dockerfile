# syntax=docker/dockerfile:1
FROM golang:1.21-bullseye

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY vendor vendor
COPY go.mod go.mod
COPY go.sum go.sum

RUN --mount=type=cache,target=/root/.cache/go-build go build -trimpath -o /meal_plan_grocery_list_initializer github.com/dinnerdonebetter/backend/cmd/localdev/meal_plan_grocery_list_initializer

ENTRYPOINT /meal_plan_grocery_list_initializer
