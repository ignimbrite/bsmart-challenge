FROM golang:1.22-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./cmd/api

FROM alpine:3.19
WORKDIR /app
RUN adduser -D -u 10001 appuser
COPY --from=builder /app/server /app/server
COPY .env.example /app/.env.example
ENV APP_ENV=production \
    HTTP_PORT=8080

USER appuser
EXPOSE 8080

CMD ["/app/server"]
