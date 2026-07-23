# Build stage
FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /taskflow-api \
    ./cmd/api

# Runtime Stage
FROM alpine:3.23

RUN addgroup -S taskflow && \
    adduser -S taskflow -G taskflow

WORKDIR /app

COPY --from=builder /taskflow-api /app/taskflow-api

USER taskflow

EXPOSE 8080

CMD ["/app/taskflow-api"]