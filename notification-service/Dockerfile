FROM golang:1.16

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o main .

CMD ["./main"]