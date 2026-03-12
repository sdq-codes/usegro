FROM golang:1.25-alpine

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Install Air for hot-reloading
RUN go install github.com/air-verse/air@latest

# Copy proto first so the replace directive (../../proto) resolves
COPY proto/ /proto/

# Copy go.mod and go.sum to cache dependencies
COPY services/base/go.mod services/base/go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY services/base/ .

EXPOSE 8090

CMD ["air", "serve-api"]
