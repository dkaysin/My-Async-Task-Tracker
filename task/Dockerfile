FROM golang:1.22-alpine3.19

WORKDIR /go/app

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN apk add --no-cache make
RUN go mod download

COPY . .

RUN make build

EXPOSE 4080

CMD ["/go/app/bin/server"]
