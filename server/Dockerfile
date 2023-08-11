FROM golang:1.20

ENV GO111MODULE=on

WORKDIR ./app

COPY . .

RUN go mod download

RUN go build -o main ./main.go

CMD ["./main"]
