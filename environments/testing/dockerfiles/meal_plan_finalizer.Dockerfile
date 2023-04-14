# syntax=docker/dockerfile:1
FROM golang:1.20-buster

WORKDIR /go/src/github.com/prixfixeco/backend

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY vendor vendor
COPY go.mod go.mod
COPY go.sum go.sum

RUN --mount=type=cache,target=/root/.cache/go-build go build -trimpath -o /meal_plan_tallier github.com/prixfixeco/backend/cmd/localdev/meal_plan_tallier

ENTRYPOINT /meal_plan_finalizer