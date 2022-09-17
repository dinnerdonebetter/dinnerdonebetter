# build stage
FROM golang:1.19-bullseye

WORKDIR /go/src/github.com/prixfixeco/api_server

COPY . .

RUN go build -trimpath -o /advanced_prep_step_creator github.com/prixfixeco/api_server/cmd/localdev/advanced_prep_step_creator

ENTRYPOINT /advanced_prep_step_creator
