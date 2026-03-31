FROM --platform=linux/amd64 golang:1.25-alpine AS builder
WORKDIR /app
COPY proto/ /proto/
COPY services/shared/go.mod services/shared/go.sum /shared/
COPY services/base/go.mod services/base/go.sum ./
RUN go mod download
COPY services/shared/ /shared/
COPY services/base/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o base-service .

FROM --platform=linux/amd64 alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/base-service .
COPY --from=builder /app/config/ ./config/
COPY --from=builder /app/internal/resources/email/html/ ./internal/resources/email/html/
RUN chmod +x base-service
EXPOSE 8090
CMD ["/app/base-service", "serve-api"]
