FROM --platform=linux/amd64 golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY proto/ /proto/
COPY services/shared/go.mod services/shared/go.sum /shared/
COPY services/catalog/go.mod services/catalog/go.sum ./
RUN go mod download
COPY services/shared/ /shared/
COPY services/catalog/ .
# Stub embed files that are gitignored but required by //go:embed at compile time
RUN mkdir -p cmd/data && \
    [ -f cmd/data/taxonomy_attributes.json ] || echo '{"attributes":[],"categories":[]}' > cmd/data/taxonomy_attributes.json && \
    [ -f cmd/data/standard_categories.json ] || echo '[]' > cmd/data/standard_categories.json
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM --platform=linux/amd64 alpine:3.19
RUN apk --no-cache add ca-certificates tzdata gettext
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config
COPY --from=builder /app/internal/resources ./internal/resources
COPY services/catalog/entrypoint.sh .
RUN chmod +x entrypoint.sh
EXPOSE 8090
EXPOSE 50051
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./server", "serve-api"]
