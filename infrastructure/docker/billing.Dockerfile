FROM --platform=linux/amd64 golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY proto/ /proto/
COPY services/shared/go.mod services/shared/go.sum /shared/
COPY services/billing/go.mod services/billing/go.sum ./
RUN go mod download
COPY services/shared/ /shared/
COPY services/billing/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM --platform=linux/amd64 alpine:3.19
RUN apk --no-cache add ca-certificates tzdata gettext
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config
COPY services/billing/entrypoint.sh .
RUN chmod +x entrypoint.sh
EXPOSE 8090
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./server", "serve-api"]
