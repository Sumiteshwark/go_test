# FROM golang:1.23-alpine
# WORKDIR /app
# COPY . .
# RUN go build -o main main.go
# CMD ["./main"]


## FOR PRODUCTION
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main

COPY .env ./

EXPOSE 8080

CMD ["./main"]
