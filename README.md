# 🌊 Interflow

**Interflow** is a high-performance, production-ready AI Gateway designed to provide a unified interface for multiple Large Language Model (LLM) providers. Built with Go, it simplifies the integration, management, and monitoring of AI services with built-in authentication, rate limiting, and asynchronous analytics.

---

## 🚀 Features

- **Unified API Interface**: Access multiple LLM providers through a single, consistent `/v1/chat` endpoint.
- **Multi-Provider Support**: Supports **OpenAI** and **Gemini** (Google AI) out of the box.
- **Robust Authentication**: Secure gateway access using `X-API-Key` header-based authentication.
- **Intelligent Rate Limiting**: Redis-backed rate limiting (default: 10 requests/minute per user) to protect your infrastructure.
- **Asynchronous Analytics**: High-performance usage tracking using a dedicated worker pool to minimize request latency.
- **SQLC Powered**: Type-safe database interactions with PostgreSQL for logging and management.
- **Graceful Shutdown**: Ensures all pending analytics data is flushed to the database and all connections are closed safely before the service stops.

---

## 🛠️ Tech Stack

- **Language**: [Go](https://go.dev/) (v1.25.4)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: [PostgreSQL](https://www.postgresql.org/) (via [pgx](https://github.com/jackc/pgx))
- **ORM/Tooling**: [sqlc](https://sqlc.dev/)
- **Caching/Rate Limiting**: [Redis](https://redis.io/)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Migrations**: Golang-migrate

---

## 📂 Project Structure

```text
.
├── cmd/
│   └── gateway/            # Application entry point (main.go)
├── internal/
│   ├── analytics/          # Usage event definitions
│   ├── cache/              # Redis client and rate limiting logic
│   ├── config/             # Environment and configuration management
│   ├── database/           # SQLC generated code and DB models
│   ├── handler/            # HTTP handlers (Chat, Health)
│   ├── middleware/         # Auth, Rate Limit, and Logging middlewares
│   ├── provider/           # LLM provider implementations (OpenAI, Gemini)
│   ├── repository/         # Database connection and initialization
│   └── service/            # Business logic (Analytics worker pool)
├── migrations/             # PostgreSQL schema migrations
└── sqlc.yaml               # SQLC configuration
```

---

## ⚙️ Getting Started

### Prerequisites

- **Go**: v1.25.4 or higher
- **PostgreSQL**: A running instance for persistent storage
- **Redis**: A running instance for rate limiting
- **sqlc**: (Optional) For regenerating database code

### Installation & Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/interflow.git
   cd interflow
   ```

2. **Setup environment variables**:
   Create a `.env` file in the root directory:
   ```env
   PORT=8080
   DATABASE_URL=postgres://user:password@localhost:5432/interflow?sslmode=disable
   REDIS_URL=localhost:6379
   OPENAI_API_KEY=your_openai_key
   GEMINI_API_KEY=your_gemini_key
   ```

3. **Database Setup**:
   Apply the migrations located in `migrations/` to your PostgreSQL database.
   ```bash
   # Example using golang-migrate
   migrate -path migrations/ -database "$DATABASE_URL" up
   ```

4. **Install dependencies**:
   ```bash
   go mod tidy
   ```

5. **Run the application**:
   ```bash
   go run cmd/gateway/main.go
   ```

---

## 📡 API Usage

### Authentication
All requests require an `X-API-Key` header.
```text
X-API-Key: your_generated_api_key
```

### Provider Selection
You can specify the LLM provider using the `X-Provider` header. If omitted, it defaults to `openai`.
- `openai`
- `gemini`

### Chat Completion
**Endpoint**: `POST /v1/chat`

**Request Body**:
```json
{
  "model": "gpt-4o",
  "messages": [
    {
      "role": "user",
      "content": "Hello, how are you?"
    }
  ],
  "temperature": 0.7,
  "max_tokens": 500
}
```

**Example Requests**:

#### OpenAI Request
```bash
curl -X POST http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test-api-key" \
  -H "X-Provider: openai" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Say hello!"}]
  }'
```

#### Gemini Request
```bash
curl -X POST http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test-api-key" \
  -H "X-Provider: gemini" \
  -d '{
    "model": "gemini-1.5-flash",
    "messages": [{"role": "user", "content": "Say hello!"}]
  }'
```

---

## 🛣️ Roadmap

- [x] **Multi-Provider Support**: Initial implementation for OpenAI and Gemini.
- [ ] **Dynamic Routing**: Automatic failover and load balancing between providers.
- [ ] **Streaming**: Support for Server-Sent Events (SSE) for real-time responses.
- [ ] **Additional Providers**: Add Anthropic (Claude) and local models (Ollama).
- [ ] **Admin Dashboard**: Web interface for managing API keys and viewing analytics.
- [ ] **Custom Policies**: Per-key rate limits and cost quotas.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
