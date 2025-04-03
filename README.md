# 📎 shortlink-api-gateway

A side project for experimenting with distributed architecture, observability (OpenTelemetry), and service communication (gRPC).

This is the **API Gateway** service for the distributed URL shortener system.  
It handles client requests, forwards them to internal services (via gRPC), and returns the results.

---

## 📌 Features

- Accepts **HTTP REST API** from clients
- Forwards requests to `shortlink-url-service` via **gRPC**
- Provides endpoints to:
  - Shorten a URL
  - Expand a short link
- Integrates with **OpenTelemetry** and **Grafana Tempo** for distributed tracing

---

## 🧱 Project Structure

```
shortlink-api-gateway/
├── go.mod
├── cmd/
│   └── gateway/
│       └── main.go
├── internal/
│   ├── server/              # Server + router
│   ├── handler/             # HTTP handlers
│   ├── service/             # gRPC client to url-service
│   ├── config/              # Configuration loader
│   └── logger/              # zap logger integration
│   └── otel/                # OpenTelemetry setup
├── proto/
│   └── public/              # Proto files
├── Dockerfile
├── docker-compose.yml       # Compose file for local dev
├── go.mod / go.sum  
├── Makefile
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- Docker + Docker Compose

### Run locally with Docker Compose

```bash
docker-compose up --build
```

This will launch:
- `gateway` on `http://localhost:8080`
- `Grafana` on `http://localhost:3000` (user: `admin`, pass: `admin`)
- `Tempo` OTLP receiver on port `4318`

### Run locally (without Docker)

```bash
go run ./cmd/gateway
```

---

## 🧪 API Endpoints

| Method | Path         | Description           |
|--------|--------------|-----------------------|
| POST   | `/shorten`   | Shortens a long URL   |
| GET    | `/:shortID`  | Redirects to original |

---

## 🧬 gRPC Public API

Defined in `proto/public/public.proto`.

```proto
service UrlPublicAPI {
  rpc ShortenUrl(ShortenRequest) returns (ShortenResponse);
  rpc ExpandUrl(ExpandRequest) returns (ExpandResponse);
}
```

---

## 📦 TODO (Next Steps)

- [x] Add OpenTelemetry tracing via stdout
- [x] Replace stdout exporter with OTLP exporter
- [x] Dockerize Gateway + Tempo + Grafana stack
- [ ] Implement gRPC client to URL service
- [ ] Unit testing and integration tests
- [ ] Inject handler
- [ ] Add RateLimiter

---

### View Traces in Grafana

1. Visit [http://localhost:3000](http://localhost:3000)
2. Go to **Explore > Tempo**
3. Use `service.name = "api-gateway"` to filter traces