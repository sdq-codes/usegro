FROM --platform=linux/amd64 golang:1.25-alpine AS builder
WORKDIR /app
COPY proto/ /proto/
COPY services/shared/go.mod services/shared/go.sum /shared/
COPY services/billing/go.mod services/billing/go.sum ./
RUN go mod download
COPY services/shared/ /shared/
COPY services/billing/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o billing-service ./cmd

FROM --platform=linux/amd64 alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/billing-service .
COPY --from=builder /app/config/ ./config/
RUN chmod +x billing-service
EXPOSE 8090
CMD ["/app/billing-service", "serve-api"]
