FROM golang:1.25-alpine AS builder
WORKDIR /workspace
COPY go.work .
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/auth/go.mod services/auth/go.sum ./services/auth/
COPY services/account/go.mod services/account/go.sum ./services/account/
RUN go mod download
COPY pkg/ ./pkg/
COPY services/account/ ./services/account/
WORKDIR /workspace/services/account
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/account-service cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /bin/account-service .
EXPOSE 3001
CMD ["./account-service"]