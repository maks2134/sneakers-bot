FROM golang:1.24-alpine as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o snakers-bot ./cmd/main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/snakers-bot .
COPY --from=builder /app/.env .
COPY --from=builder /app/data ./data

CMD ["./snakers-bot"]