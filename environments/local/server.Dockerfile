# build stage
FROM golang:1.19-bullseye

WORKDIR /go/src/github.com/prixfixeco/api_server

RUN	apt-get update && apt-get install -y \
	--no-install-recommends \
	entr \
	&& rm -rf /var/lib/apt/lists/*
ENV ENTR_INOTIFY_WORKAROUND=true

ENTRYPOINT echo "please wait for server to start" && find . -type f \( -iname "*.go*" ! -iname "*_test.go" \) | entr -r go run github.com/prixfixeco/api_server/cmd/server