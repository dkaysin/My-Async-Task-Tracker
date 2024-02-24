FROM --platform=linux/amd64 golang:1.22-alpine3.19

WORKDIR /go/app

COPY . .

RUN apk add --no-cache make
RUN go mod download
RUN make build

EXPOSE 4080

CMD ["/go/app/bin/server"]
