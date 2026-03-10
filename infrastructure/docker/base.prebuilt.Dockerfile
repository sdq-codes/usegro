FROM --platform=linux/amd64 alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY dist/base-service .
COPY services/base/config/ ./config/
RUN chmod +x base-service
EXPOSE 8090
CMD ["/app/base-service", "serve-api"]
