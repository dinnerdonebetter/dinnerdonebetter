# build stage
FROM golang:1.19-bullseye AS build-stage

WORKDIR /go/src/github.com/prixfixeco/backend

COPY . .

RUN go build -trimpath -o /server github.com/prixfixeco/backend/cmd/server

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Create the user
RUN groupadd --gid 1000 server-runner && useradd --uid 1000 --gid 1000 -m server-runner

# ********************************************************
# * Anything else you want to do like clean up goes here *
# ********************************************************

# [Optional] Set the default user. Omit if you want to keep the default as root.
USER $server-runner

COPY --from=build-stage /server /server

ENTRYPOINT ["/server"]
