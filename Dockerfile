FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o cli ./cmd/cli/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cli .

# ENTRYPOINT ["tail", "-f", "/dev/null"]

CMD ["./cli"]

