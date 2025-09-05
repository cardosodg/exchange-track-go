# Stage 1: Build
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates bash

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o exchangetrack ./cmd/exchangetrack/main.go

# Stage 2: Minimal runtime
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates bash

COPY --from=builder /app/exchangetrack .

CMD ["./exchangetrack"]
