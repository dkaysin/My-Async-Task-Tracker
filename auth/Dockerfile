FROM golang:1.22-alpine3.19

RUN apk add --no-cache make

WORKDIR /go/app
COPY ./go.work.auth ./go.work

COPY ./auth/go.mod ./auth/
COPY ./auth/go.sum ./auth/

COPY ./schema_registry/go.mod ./schema_registry/
COPY ./schema_registry/go.sum ./schema_registry/

WORKDIR /go/app/schema_registry
RUN go mod download

WORKDIR /go/app/auth
RUN go mod download

WORKDIR /go/app
COPY ./schema_registry ./schema_registry
COPY ./auth ./auth

WORKDIR /go/app/auth
RUN make build

EXPOSE 4080

CMD ["/go/app/auth/bin/server"]
