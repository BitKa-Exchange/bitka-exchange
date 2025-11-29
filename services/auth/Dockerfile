# Stage 1: Builder
FROM golang:1.25-alpine AS builder

# Set the working directory to the root of the workspace inside the container
WORKDIR /workspace

# Install build dependencies (if needed, mostly for CGO, but we will disable CGO)
# RUN apk add --no-cache git

# 1. Copy the workspace definition file
COPY go.work .

# 2. Copy the go.mod and go.sum files for ALL modules
#    We do this before copying source code to cache dependencies.
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/auth/go.mod services/auth/go.sum ./services/auth/

# 3. Download dependencies
#    go.work makes Go aware of the module relationship
RUN go mod download

# 4. Copy the actual source code
COPY pkg/ ./pkg/
COPY services/auth/ ./services/auth/

# 5. Build the binary
#    We disable CGO for a static binary (works better in Alpine/Scratch)
WORKDIR /workspace/services/auth
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-service cmd/main.go

# Stage 2: Runner
FROM alpine:latest

# Install CA certificates (Required for JWKS/HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /bin/auth-service .

# Copy the .env file if you want to bundle it (Optional, usually better to inject via docker-compose)
# COPY services/auth/.env .

# Expose the application port
EXPOSE 3000

# Run the binary
CMD ["./auth-service"]