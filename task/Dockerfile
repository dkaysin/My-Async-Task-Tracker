FROM golang:1.22-alpine3.19

RUN apk add --no-cache make

WORKDIR /go/app
COPY ./go.work.task ./go.work

COPY ./task/go.mod ./task/
COPY ./task/go.sum ./task/

COPY ./schema_registry/go.mod ./schema_registry/
COPY ./schema_registry/go.sum ./schema_registry/

WORKDIR /go/app/schema_registry
RUN go mod download

WORKDIR /go/app/task
RUN go mod download

WORKDIR /go/app
COPY ./schema_registry ./schema_registry
COPY ./task ./task

WORKDIR /go/app/task
RUN make build

EXPOSE 4080

CMD ["/go/app/task/bin/server"]
