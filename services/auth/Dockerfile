# Build stage
FROM golang:1.24 AS builder
ARG SERVICE_NAME
WORKDIR /bitka-exchange

COPY ../go.work .
COPY ./services/auth/go.mod ./services/auth/go.mod
COPY ./services/account/go.mod ./services/account/go.mod
COPY ./pkg/go.mod ./pkg/go.mod

RUN go mod download
COPY ./pkg ./pkg
COPY ./services/${SERVICE_NAME} ./services/${SERVICE_NAME}

RUN CGO_ENABLED=0 GOOS=linux go build \
  -o /bin/${SERVICE_NAME}-service \
  ./services/${SERVICE_NAME}/cmd/server/main.go

# Runtime stage
FROM alpine:latest
ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

COPY --from=builder /bin/${SERVICE_NAME}-service /bin/server

CMD ["/bin/server"]
