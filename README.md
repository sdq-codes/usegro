# useGro

**The Operating System for SMEs**

useGro is a global all-in-one business platform for small and medium enterprises. It replaces the 6–10 disconnected tools an SME uses daily — CRM, storefront, inventory, invoicing, payments, WhatsApp, social media, and analytics — with a single connected workspace.

---

## What useGro Does

| Module | What it covers |
|---|---|
| **CRM & Leads** | Contacts, pipeline, lifecycle stages, activity feed, segmentation |
| **Commerce** | Storefront (My Store), products, inventory, POS, invoicing, payments |
| **Messaging & WhatsApp BSP** | Unified inbox, WhatsApp Business API, campaigns, chatbot, automation |
| **Social Channels** | TikTok Shop, Instagram Shopping, Facebook Ads, Email, SMS |
| **Analytics** | Sales, campaigns, store, CSAT, revenue attribution |

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend services | Go + gRPC |
| WhatsApp / Meta interface | Node.js microservice |
| Frontend | Vue 3 + TypeScript |
| Primary database | PostgreSQL |
| Message & event store | DynamoDB |
| Event bus + cache | Redis (pub/sub) |
| Inter-service comms | gRPC |
| File storage | S3 / Cloudflare R2 |
| Payments | Stripe + Paystack + Flutterwave |
| Containers | Docker + Docker Compose → Kubernetes |

---

## Services

| Service | Lang | Port | Responsibility |
|---|---|---|---|
| `gateway` | Go | 8080 | REST API Gateway — single entry point for frontend, JWT auth, WebSocket |
| `crm` | Go | 50051 | Contacts, leads, pipeline, lifecycle tracking, activity feed |
| `commerce` | Go | 50052 | Storefront, products, inventory, POS, invoicing, payments |
| `messaging` | Go | 50053 | Unified inbox, WhatsApp, campaigns, templates, broadcast, chatbot |
| `channels` | Go | 50054 | TikTok, Instagram, Facebook Ads, Email (SendGrid), SMS (Twilio) |
| `analytics` | Go | 50055 | Reports, CSAT, revenue attribution, dashboards |
| `billing` | Go | 50056 | useGro subscription plans, usage metering, Stripe |
| `whatsapp-gateway` | Node.js | 3000 | Meta Cloud API — webhook ingestion, message send, media, WA Flows |

---

## Architecture

```
Vue 3 Frontend
      │
      │  REST / WebSocket
      ▼
┌─────────────────────┐
│   Go API Gateway    │  :8080
└─────────────────────┘
      │
      │  gRPC
      ▼
┌──────────┬───────────┬──────────┬───────────┬──────────┬──────────┐
│  crm     │ commerce  │messaging │ channels  │analytics │ billing  │
│ :50051   │ :50052    │ :50053   │ :50054    │ :50055   │ :50056   │
└──────────┴───────────┴──────────┴───────────┴──────────┴──────────┘
      │
      │  Redis pub/sub
      ▼
┌─────────────────────┐
│ whatsapp-gateway    │  Node.js  :3000
└─────────────────────┘
      │
      │  HTTPS
      ▼
   Meta Cloud API
```

**Redis channel conventions:**
```
usegro:inbound:message:{wabaId}   — Node.js → Go  (inbound WhatsApp message)
usegro:inbound:status:{wabaId}    — Node.js → Go  (delivery status update)
usegro:outbound:send:{wabaId}     — Go → Node.js  (trigger message send)
```

**Data storage split:**
- **PostgreSQL** — contacts, leads, products, inventory, invoices, payments, orders, campaigns, billing
- **DynamoDB** — conversations, messages, webhook events, activity feed, analytics events (high-volume, time-series)
- **Redis** — pub/sub event bus, session cache, agent presence, rate limiting

---

## Project Structure

```
usegro/
├── proto/
│   ├── whatsapp/          # WA message types, inbound/outbound, status updates
│   ├── crm/               # Tenant, contact, lead, auth RPC definitions
│   ├── messaging/         # Conversation, message, template, campaign RPCs
│   ├── commerce/          # Product, order, invoice, inventory RPCs
│   ├── channels/          # Social channel integration RPCs
│   ├── analytics/         # Report and metrics RPCs
│   └── billing/           # Subscription and usage RPCs
│
├── services/
│   ├── gateway/           # Go — REST API gateway + WebSocket hub
│   ├── crm/               # Go — CRM service
│   │   ├── cmd/           # main.go, server.go
│   │   └── internal/      # contacts, leads, tenants, auth, lifecycle
│   ├── commerce/          # Go — Commerce service
│   │   ├── cmd/
│   │   └── internal/      # products, inventory, storefront, pos, invoicing, payments
│   ├── messaging/         # Go — Messaging service
│   │   ├── cmd/           # main.go, subscriber.go, handler.go
│   │   └── internal/      # inbox, campaigns, templates, broadcast, chatbot
│   ├── channels/          # Go — Social channels service
│   ├── analytics/         # Go — Analytics service
│   ├── billing/           # Go — Billing service
│   └── whatsapp-gateway/  # Node.js — Meta Cloud API interface
│       └── src/
│           ├── webhook/   # Inbound handler + Meta signature verification
│           ├── sender/    # Meta Cloud API message sender
│           ├── redis/     # Publisher + subscriber
│           └── grpc/      # gRPC server
│
├── frontend/              # Vue 3 + TypeScript
│   └── src/
│       ├── views/         # Pages: inbox, contacts, pipeline, store, analytics
│       ├── components/    # Shared UI components
│       ├── stores/        # Pinia stores
│       └── api/           # API client (typed, auto-generated from proto)
│
├── infrastructure/
│   ├── docker/            # Dockerfile per service
│   └── postgres/
│       └── migrations/    # Numbered SQL migration files
│
├── docker-compose.yml     # Full local dev stack
├── Makefile               # make dev / build / migrate / proto / test
└── .env.example           # All environment variables documented
```

---

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker + Docker Compose
- `protoc` + Go and Node.js proto plugins (for `make proto`)

### 1. Clone and configure

```bash
git clone https://github.com/your-org/usegro.git
cd usegro
cp .env.example .env
# Fill in your Meta, Stripe, AWS, and Redis credentials
```

### 2. Start infrastructure

```bash
make infra
# Starts: PostgreSQL, Redis, DynamoDB local
```

### 3. Run migrations

```bash
make migrate
```

### 4. Generate protobuf files

```bash
make proto
```

### 5. Start all services (dev)

```bash
make dev
# Starts all Go services + Node.js WhatsApp Gateway
# API Gateway available at http://localhost:8080
```

### 6. Start frontend

```bash
cd frontend
npm install
npm run dev
# Frontend at http://localhost:5173
```

### Docker (full stack)

```bash
docker compose up -d
```

---

## Environment Variables

See `.env.example` for the full list. Key groups:

```bash
# Meta / WhatsApp BSP
WHATSAPP_ACCESS_TOKEN=
WHATSAPP_WEBHOOK_VERIFY_TOKEN=
WHATSAPP_APP_SECRET=

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_DB=usegro

# DynamoDB
DYNAMODB_REGION=eu-west-1
DYNAMODB_ENDPOINT=http://localhost:8000   # local only

# Redis
REDIS_HOST=localhost

# Auth
JWT_SECRET=

# Payments
STRIPE_SECRET_KEY=
PAYSTACK_SECRET_KEY=
FLUTTERWAVE_SECRET_KEY=

# Email / SMS
SENDGRID_API_KEY=
TWILIO_ACCOUNT_SID=
```

---

## Feature Priority Matrix

> **P0** = launch blocker · **P1** = competitive must-have · **P2** = growth feature

### CRM & Leads

| Feature | Description | Priority | Phase |
|---|---|:---:|:---:|
| Multi-tenant auth | Workspace per SME, user roles, JWT | **P0** | 1 |
| Contact management | Full profiles, custom fields, tags, activity feed | **P0** | 1 |
| Lead pipeline | Kanban + list view, customisable stages, deal value | **P0** | 1 |
| Lifecycle stage tracking | Lead → MQL → SQL → Customer → Churned | **P0** | 1 |
| Contact segmentation | Filter by tag, stage, spend, channel source | P1 | 1 |
| Opt-in / opt-out compliance | Consent tracking per contact | P1 | 2 |
| Lead scoring | Automated scoring based on activity and engagement | P2 | 3 |

### Commerce

| Feature | Description | Priority | Phase |
|---|---|:---:|:---:|
| Product catalogue | Products, services, variants, pricing | **P0** | 1 |
| Inventory tracking | Stock levels, low-stock alerts, adjustment logs | **P0** | 1 |
| My Store (storefront) | Public product page with unique shareable URL | **P0** | 1 |
| Invoicing (core) | Create, send, track — draft / sent / paid / overdue | **P0** | 1 |
| Stripe payments | Card payments, payment links, refunds | **P0** | 1 |
| Invoice PDF + client portal | PDF generation, self-service customer payment portal | P1 | 1 |
| Point of Sale (POS) | In-person sales, cart, receipts, inventory sync | P1 | 2 |
| Paystack / Flutterwave | African payment processor integration | P1 | 2 |
| Multi-location inventory | Stock tracking across multiple locations | P2 | 2 |
| Recurring invoices | Subscription billing, auto-send on schedule | P2 | 2 |
| Custom store domain | Map your own domain to My Store | P2 | 3 |

### Messaging & WhatsApp BSP

| Feature | Description | Priority | Phase |
|---|---|:---:|:---:|
| WhatsApp BSP (core) | Cloud API, webhook, send/receive, multi-tenant WABAs | **P0** | 2 |
| Unified inbox | Multi-agent, assignment, labels, real-time WebSocket | **P0** | 2 |
| Message templates | Create, submit to Meta, track approval status | **P0** | 2 |
| Rich media messages | Images, video, docs, audio, location | P1 | 2 |
| Campaign management | Broadcast, segmentation, scheduling, analytics | P1 | 2 |
| Email campaigns | Templates, send, delivery + open rate tracking | P1 | 2 |
| SMS campaigns | Twilio / Africa's Talking integration | P1 | 2 |
| CSAT surveys | Post-conversation satisfaction scoring | P1 | 2 |
| Rule-based chatbot | Keyword triggers, decision trees, quick replies | P1 | 2 |
| WhatsApp Flows | Native structured UX forms inside WhatsApp | P2 | 3 |
| Automation flow builder | Drag-and-drop triggers, conditions, actions | P2 | 3 |
| AI chatbot | LLM-powered, intent detection, human handoff | P2 | 3 |
| WhatsApp Calls | Voice calling via WhatsApp | P2 | 3 |
| In-chat payments | WhatsApp Pay + payment links in conversation | P2 | 3 |

### Social Channels

| Feature | Description | Priority | Phase |
|---|---|:---:|:---:|
| Instagram DMs + Shopping | DMs in unified inbox, product catalogue sync | P1 | 2 |
| Facebook Messenger + Ads | Messenger in inbox, click-to-WhatsApp ads | P1 | 2 |
| TikTok Shop + Ads | Product sync, order management, ad campaigns | P2 | 3 |

### Analytics

| Feature | Description | Priority | Phase |
|---|---|:---:|:---:|
| Analytics dashboard | Sales, CRM, campaign, and store metrics | P1 | 2 |
| Revenue attribution | Which campaign/channel drove each sale | P2 | 3 |
| Custom report builder | User-defined reports and scheduled exports | P2 | 3 |
| Mobile app | iOS + Android for agents and business owners | P2 | 3 |

---

## Build Phases

### Phase 1 — Revenue MVP (Weeks 1–10)
Auth & multi-tenancy · Contacts & lead pipeline · Product catalogue · Inventory · Invoicing · Stripe payments · My Store (storefront) · Go API Gateway · Vue frontend shell

### Phase 2 — Full Messaging & Commerce (Weeks 11–22)
WhatsApp BSP (full) · Unified inbox · Campaign management · POS · Email & SMS · Instagram + Facebook · CSAT · Analytics dashboard · Paystack/Flutterwave

### Phase 3 — Market Leader (Weeks 23–36)
TikTok Shop + Ads · AI chatbot · Drag-and-drop automation builder · WhatsApp Flows · In-chat payments · Revenue attribution · Mobile app (iOS + Android)

---

## Key Design Decisions

- **Node.js is isolated to Meta API only** — all business logic lives in Go. The whatsapp-gateway is a thin adapter, not a business service.
- **Redis pub/sub as the event bus** — decouples the hot inbound path. Node.js publishes, Go subscribes. No tight gRPC coupling on webhook delivery.
- **DynamoDB for messages** — conversation history is high-volume, append-only, and time-series. No complex joins required.
- **PostgreSQL for everything relational** — contacts, invoices, pipeline, billing. Row-level security enforces tenant isolation.
- **Monorepo** — shared proto definitions, consistent tooling, easier cross-service refactoring.
- **Multi-tenant by workspace** — every table has `tenant_id`. One useGro account = one SME workspace.

---

## Meta BSP Status

useGro is applying for Meta BSP (Business Solution Provider) status to operate the WhatsApp Business API directly. Until approved, the whatsapp-gateway runs against a test WABA. Production client WABAs require BSP approval (~2–4 weeks from Meta).

---

## Contributing

1. Branch from `main` — `feat/your-feature` or `fix/your-fix`
2. Run `make lint` and `make test` before pushing
3. Proto changes require running `make proto` and committing generated files
4. Database changes require a new numbered migration file in `infrastructure/postgres/migrations/`

---

*useGro — grow without the chaos.*