.PHONY: dev infra migrate proto lint test build

# ─── Infrastructure ────────────────────────────────────────
infra:
	docker compose up -d postgres redis dynamodb-local

infra-down:
	docker compose down

# ─── Proto generation ──────────────────────────────────────
proto:
	@echo "Generating Go protobuf files..."
	find proto -name "*.proto" | xargs -I{} protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		{}
	@echo "Copying protos to Node.js gateway..."
	cp proto/whatsapp/*.proto services/whatsapp-gateway/proto/

# ─── Migrations ────────────────────────────────────────────
migrate:
	@echo "Running PostgreSQL migrations..."
	docker compose exec postgres psql -U usegro -d usegro \
		-f /docker-entrypoint-initdb.d/001_init.sql

migrate-create:
	@read -p "Migration name: " name; \
	touch infrastructure/postgres/migrations/$$(date +%03d)_$$name.sql

# ─── Dev (all services) ────────────────────────────────────
dev: infra
	@echo "Starting all Go services..."
	cd services/crm && go run ./cmd &
	cd services/messaging && go run ./cmd &
	cd services/billing && go run ./cmd &
	cd services/analytics && go run ./cmd &
	cd services/gateway && go run ./cmd &
	@echo "Starting Node.js WhatsApp Gateway..."
	cd services/whatsapp-gateway && npm run dev &
	@echo "All services started. Gateway: http://localhost:8080"

# ─── Build ─────────────────────────────────────────────────
build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

# ─── Test ──────────────────────────────────────────────────
test:
	cd services/crm && go test ./...
	cd services/messaging && go test ./...
	cd services/billing && go test ./...
	cd services/analytics && go test ./...
	cd services/gateway && go test ./...
	cd services/whatsapp-gateway && npm test

# ─── Lint ──────────────────────────────────────────────────
lint:
	cd services/crm && golangci-lint run ./...
	cd services/messaging && golangci-lint run ./...
	cd services/gateway && golangci-lint run ./...
	cd services/whatsapp-gateway && npm run lint
