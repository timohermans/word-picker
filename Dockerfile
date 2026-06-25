# Build stage
FROM golang:1.26 AS builder

ARG APP_NAME="app"

WORKDIR /app

# Download dependencies first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ${APP_NAME} ./cmd/${APP_NAME}

# Runtime stage
FROM alpine:3.24.1

ARG APP_NAME="app"

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/${APP_NAME} .

EXPOSE 8080

CMD ["./${APP_NAME}"]