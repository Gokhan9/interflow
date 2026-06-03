# 🌊 Interflow

**Interflow** is a high-performance, production-ready AI Gateway designed to provide a unified interface for multiple Large Language Model (LLM) providers. Built with efficiency and scalability in mind, it simplifies the integration, management, and monitoring of AI services.

---

## 🚀 Features

- **Unified API Interface**: Access OpenAI, Gemini, Anthropic, and more through a single, consistent API.
- **Robust Authentication**: Secure your gateway with API key-based authentication.
- **Intelligent Rate Limiting**: Redis-backed rate limiting to protect your infrastructure and manage costs.
- **Detailed Analytics & Logging**: Track token usage, latency, and status codes for every request.
- **High Performance**: Developed in Go using the Gin framework for minimal overhead and maximum throughput.
- **SQLC Powered**: Type-safe database interactions with PostgreSQL.

---

## 🛠️ Tech Stack

- **Language**: [Go](https://go.dev/) (v1.25+)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **ORM/Tooling**: [sqlc](https://sqlc.dev/)
- **Caching/Rate Limiting**: [Redis](https://redis.io/)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Migrations**: Golang-migrate

---

## 📂 Project Structure

```text
.
├── cmd/
│   └── gateway/            # Application entry point
├── internal/
│   ├── analytics/          # Usage tracking and event processing
│   ├── cache/              # Redis client and caching logic
│   ├── config/             # Environment and configuration management
│   ├── database/           # SQLC generated code and DB models
│   ├── middleware/         # Auth, Rate Limit, and Logging middlewares
│   ├── provider/           # LLM provider implementations (OpenAI, etc.)
│   ├── repository/         # Data access layer
│   └── service/            # Business logic
├── migrations/             # PostgreSQL schema migrations
└── sqlc.yaml               # SQLC configuration
```

---

## ⚙️ Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install) installed.
- [PostgreSQL](https://www.postgresql.org/download/) instance.
- [Redis](https://redis.io/docs/getting-started/) instance.
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) (optional, for regenerating DB code).

### Installation

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
   ANTHROPIC_API_KEY=your_anthropic_key
   ```

3. **Run migrations**:
   (Assuming you have a migration tool or apply manually)
   ```sql
   -- Use the scripts in migrations/ folder
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

## 🛣️ Roadmap

- [ ] Implement Anthropic & Gemini providers.
- [ ] Add dynamic provider routing (failover & load balancing).
- [ ] Develop a Dashboard for usage monitoring.
- [ ] Streaming support for chat responses.
- [ ] Comprehensive Unit & Integration tests.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
