FROM golang:1.21.1

WORKDIR /app

COPY go.mod /app
COPY go.sum /app
COPY main.go /app

RUN go mod tidy

CMD ["go", "run", "/app/main.go"]