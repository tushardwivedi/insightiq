# backend/Dockerfile.debug - More secure Alpine-based option
FROM golang:alpine AS builder

# Create non-root user for building
RUN adduser -D -g '' appuser

# Install only essential packages and update
RUN apk update && apk upgrade && \
    apk add --no-cache \
        git \
        ca-certificates \
        tzdata \
        curl && \
    rm -rf /var/cache/apk/* /tmp/*

WORKDIR /app

# Ensure the app directory is owned by appuser
RUN chown -R appuser:appuser /app

# Copy and download dependencies as non-root
USER appuser
COPY --chown=appuser:appuser go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and build
COPY --chown=appuser:appuser . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -trimpath \
    -a -installsuffix cgo \
    -o main cmd/agent/main.go

# Final stage - minimal Alpine with security updates
FROM alpine:3.20

# Install security updates and minimal runtime deps
RUN apk update && apk upgrade && \
    apk add --no-cache \
        ca-certificates \
        curl \
        tzdata && \
    rm -rf /var/cache/apk/* /tmp/* /var/tmp/*

# Create non-root user
RUN adduser -D -g '' -s /sbin/nologin appuser

# Copy binary with proper ownership
COPY --from=builder --chown=appuser:appuser /app/main /app/main

# Switch to non-root user
USER appuser
WORKDIR /app

EXPOSE 8080

# Health check with curl
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/health || exit 1

CMD ["./main"]