FROM golang:1.14

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ./bin/main .

COPY config.json /app/bin

WORKDIR /app/bin

CMD ["/app/bin/main"]


