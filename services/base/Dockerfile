# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy proto first so the replace directive (../../proto) resolves
COPY proto/ /proto/

# Copy go.mod and go.sum to cache dependencies
COPY services/base/go.mod services/base/go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY services/base/ .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates gettext

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

COPY services/base/entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8090

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./server", "serve-api"]
