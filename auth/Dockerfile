FROM golang:1.22

WORKDIR /go/app

COPY . .

RUN go mod download
RUN make build

EXPOSE 4080

CMD ["/go/app/bin/server"]
