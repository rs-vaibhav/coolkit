# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
# Install git for dependencies
RUN apk add --no-cache git
# Copy go.mod first for layer caching
COPY go.mod go.sum ./
RUN go mod download
# Copy source
COPY . .
# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o coolkit ./cmd/server

# Stage 2: Runtime
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
# Copy binary
COPY --from=builder /app/coolkit .
# Copy migrations
COPY --from=builder /app/migrations ./migrations
# Copy frontend
COPY --from=builder /app/frontend ./frontend
# Expose port
EXPOSE 8080
# Run
CMD ["./coolkit"]
