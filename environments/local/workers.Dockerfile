# build stage
FROM golang:1.17-stretch

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN	apt-get update && apt-get install -y \
	--no-install-recommends \
	entr \
	&& rm -rf /var/lib/apt/lists/*
ENV ENTR_INOTIFY_WORKAROUND=true

ENTRYPOINT echo "please wait for workers to start" && sleep 15 && find . -type f \( -iname "*.go*" ! -iname "*_test.go" \) | entr -r go run gitlab.com/prixfixe/prixfixe/cmd/workers
