FROM golang:1.22.5

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o main .

EXPOSE 8081

CMD ["./main"]