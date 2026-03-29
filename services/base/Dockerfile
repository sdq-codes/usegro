FROM golang:1.25-alpine

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Install Air for hot-reloading
RUN go install github.com/air-verse/air@latest

# Copy proto first so the replace directive (../../proto) resolves
COPY proto/ /proto/

# Copy shared go.mod for dependency resolution, then service go.mod
COPY services/shared/go.mod services/shared/go.sum /shared/
COPY services/base/go.mod services/base/go.sum ./
RUN go mod download

# Copy shared and service source
COPY services/shared/ /shared/
COPY services/base/ .

EXPOSE 8090

CMD ["air", "serve-api"]
