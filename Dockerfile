# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/migrate ./cmd/migrate

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H -u 10001 appuser

WORKDIR /app

COPY --from=builder /app/bin/api /app/bin/api
COPY --from=builder /app/bin/migrate /app/bin/migrate
COPY migrations /app/migrations
COPY docker-entrypoint.sh /app/docker-entrypoint.sh

RUN sed -i 's/\r$//' /app/docker-entrypoint.sh \
    && chmod +x /app/docker-entrypoint.sh \
    && mkdir -p /app/storage \
    && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
