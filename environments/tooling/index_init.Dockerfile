FROM golang:stretch

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

ADD cmd/tools/index_initializer .

RUN go build -trimpath -o /index_initializer

CMD /index_initializer