# 📎 shortlink-api-gateway

This is the **API Gateway** service for the distributed URL shortener system.  
It handles client requests, forwards them to internal services (via gRPC), and returns the results.

---

## 📌 Features

- Accepts **HTTP REST API** from clients
- Forwards requests to `shortlink-url-service` via **gRPC`
- Provides endpoints to:
  - Shorten a URL
  - Expand a short link

---

## 🧱 Project Structure

```
shortlink-api-gateway/
├── go.mod
├── cmd/
│   └── gateway/
│       └── main.go
├── internal/
│   ├── server/          # regist server and router
│   ├── handler/         # HTTP/gRPC handlers
│   ├── service/         # gRPC client to url-service
│   ├── logger/          # logger
│   ├── otel/            # opentelemetry
│   └── config/          # Configuration loading
├── proto/
│   └── public/          # Proto file for client-facing API
├── api/                 # HTTP router setup (e.g. Gin/Echo)
├── Dockerfile
├── Makefile
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.20+
- protoc & `protoc-gen-go`, `protoc-gen-go-grpc`
- Docker (optional)
- `shortlink-url-service` running locally or remotely

### Install dependencies

```bash
go mod tidy
```

### Run locally

```bash
go run ./cmd/main.go
```

By default, the server runs at `http://localhost:8080`.

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

## 🐳 Docker

```bash
docker build -t shortlink-api-gateway .
docker run -p 8080:8080 shortlink-api-gateway
```

---

## 📦 TODO (Next Steps)

- [ ] Add gRPC client to connect with `url-service`
- [ ] Add unit tests
- [x] Add tracing and logging
- [ ] Implement retry logic

---

## 🧠 License

MIT License © YourName
