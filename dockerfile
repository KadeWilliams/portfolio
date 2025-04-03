# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app /app/app
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static

expose 8080

CMD ["/app/app"]
