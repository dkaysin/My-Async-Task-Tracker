FROM golang:1.22-alpine3.19

RUN apk add --no-cache make

WORKDIR /go/app
COPY ./go.work.analytics ./go.work

COPY ./analytics/go.mod ./analytics/
COPY ./analytics/go.sum ./analytics/

COPY ./schema_registry/go.mod ./schema_registry/
COPY ./schema_registry/go.sum ./schema_registry/

WORKDIR /go/app/schema_registry
RUN go mod download

WORKDIR /go/app/analytics
RUN go mod download

WORKDIR /go/app
COPY ./schema_registry ./schema_registry
COPY ./analytics ./analytics

WORKDIR /go/app/analytics
RUN make build

EXPOSE 4080

CMD ["/go/app/analytics/bin/server"]
