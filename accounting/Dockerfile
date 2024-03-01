FROM golang:1.22-alpine3.19

RUN apk add --no-cache make

WORKDIR /go/app
COPY ./go.work.accounting ./go.work

COPY ./accounting/go.mod ./accounting/
COPY ./accounting/go.sum ./accounting/

COPY ./schema_registry/go.mod ./schema_registry/
COPY ./schema_registry/go.sum ./schema_registry/

WORKDIR /go/app/schema_registry
RUN go mod download

WORKDIR /go/app/accounting
RUN go mod download

WORKDIR /go/app
COPY ./schema_registry ./schema_registry
COPY ./accounting ./accounting

WORKDIR /go/app/accounting
RUN make build

EXPOSE 4080

CMD ["/go/app/accounting/bin/server"]
