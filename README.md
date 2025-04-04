# 📎 shortlink-api-gateway

A side project for experimenting with distributed architecture, observability (OpenTelemetry), and service communication (gRPC).

This is the **API Gateway** service for the distributed URL shortener system.  
It handles client requests, forwards them to internal services (via gRPC), and returns the results.

## 🎬 Demo

![Shortlink Demo](demo/shortlink-demo.gif)

---

## 📌 Features

- Accepts **HTTP REST API** from clients
- Forwards requests to `shortlink-url-service` via **gRPC**
- Provides endpoints to:
  - Shorten a URL
  - Expand a short link
- Integrates with **OpenTelemetry** and **Grafana Tempo** for distributed tracing

---

## 🔍 Observability Stack

This project implements a modern observability stack using OpenTelemetry, Tempo and Grafana:

### Components

- **OpenTelemetry**: Framework for collecting traces, metrics and logs
  - Automatically creates spans for HTTP requests via `otelgin` middleware
  - Captures gRPC calls with `otelgrpc` stats handler
  - Provides trace context propagation across service boundaries

- **Tempo**: High-scale, cost-effective distributed tracing backend
  - Stores all trace data without sampling
  - Efficiently queries traces by TraceID
  - Retains traces for 7 days by default

- **Grafana**: Visualization platform for all observability data
  - Provides trace exploration UI
  - Shows service graphs and span details
  - Offers unified view of the system behavior

### How It Works

1. When a request arrives at the gateway:
   - HTTP middleware creates a root span with unique TraceID
   - Context with TraceID propagates through the application
   - Child spans are created for downstream operations

2. When the gateway calls the URL service via gRPC:
   - The parent context is passed to maintain trace continuity
   - gRPC handler automatically creates child spans
   - Trace shows complete request flow across services

3. All spans are exported to Tempo via OTLP protocol:
   - HTTP spans show route, method, status code
   - gRPC spans show service name, method, status
   - Custom attributes can be added to provide more context

### Using the Observability Tools

1. **View Request Traces**:
   - Go to Grafana: `http://localhost:3000`
   - Navigate to Explore > Tempo
   - Search by service name: `service.name = "api-gateway"`
   - Click any trace to see the detailed timeline

2. **Analyze Performance**:
   - Examine span durations to identify bottlenecks
   - Look for anomalies in request patterns
   - Filter by HTTP status codes to find errors

3. **Debug Issues**:
   - Trace through the entire request lifecycle
   - See exactly where failures occur
   - Compare successful vs failed requests

---

## 🧱 Project Structure

```
shortlink-gateway/
├── cmd/
│   └── gateway/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration loader
│   ├── engine/                  # Gin engine setup
│   ├── handler/                 # HTTP handlers
│   │   ├── expand.go            # URL expansion handler
│   │   └── shorten.go           # URL shortening handler
│   ├── logger/                  # Zap logger integration
│   ├── middleware/              # HTTP middleware
│   ├── otel/                    # OpenTelemetry setup
│   ├── server/                  # Server and router
│   └── service/                 # Service layer implementation
│       ├── url_service.go       # URLService interface and Mock implementation
│       └── url_grpc_client.go   # gRPC client implementation
├── proto/                       # Protocol Buffers definitions
│   ├── shortlink.proto          # Service and message definitions
│   ├── shortlink.pb.go          # Generated proto code
│   └── shortlink_grpc.pb.go     # Generated gRPC code
├── tempo/                       # Tempo distributed tracing config
├── grafana/                     # Grafana visualization config
├── .env                         # Environment variables
├── .env_example                 # Environment variables example
├── config.yaml                  # Application configuration
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Docker build file
├── go.mod / go.sum              # Go module dependencies
├── Makefile                     # Project build and management
└── README.md                    # Project documentation
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

Defined in `proto/shortlink.proto`.

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
- [x] Implement gRPC client to URL service
- [ ] Unit testing and integration tests
- [x] Inject handler
- [ ] Add RateLimiter
- [ ] Error Handle improvement