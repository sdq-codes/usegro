FROM --platform=linux/amd64 alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY dist/crm-service .
COPY services/crm/config/ ./config/
RUN chmod +x crm-service
EXPOSE 8090
EXPOSE 50051
CMD ["/app/crm-service", "serve-api"]
