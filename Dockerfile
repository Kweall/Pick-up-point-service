FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/cli ./cmd/point-cli/

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/cli .

CMD ["./cli"]
