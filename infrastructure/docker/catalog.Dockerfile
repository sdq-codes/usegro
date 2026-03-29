FROM --platform=linux/amd64 golang:1.25-alpine AS builder
WORKDIR /app
COPY proto/ /proto/
COPY services/shared/go.mod services/shared/go.sum /shared/
COPY services/catalog/go.mod services/catalog/go.sum ./
RUN go mod download
COPY services/shared/ /shared/
COPY services/catalog/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o catalog-service ./cmd

FROM --platform=linux/amd64 alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/catalog-service .
COPY --from=builder /app/config/ ./config/
RUN chmod +x catalog-service
EXPOSE 8090
EXPOSE 50051
CMD ["/app/catalog-service", "serve-api"]
