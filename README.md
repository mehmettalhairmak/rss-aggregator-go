<div align="center">

# 📰 RSS Aggregator

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

**A production-ready RSS feed aggregator built with Go, featuring clean architecture, type-safe SQL, and scalable design patterns.**

[Features](#-features) • [Quick Start](#-quick-start) • [API Docs](#-api-documentation) • [Roadmap](#-roadmap)

</div>

---

## 🎯 Overview

RSS Aggregator is a RESTful API service for managing RSS feed subscriptions. Built with Go, it demonstrates clean architecture, proper separation of concerns, and production-ready patterns.

**Status:** ✅ Core API Complete | 🚧 Feed Scraping In Development

## ✨ Features

### Current

- ✅ User management with API key authentication
- ✅ RSS feed CRUD operations
- ✅ Follow/unfollow feeds
- ✅ PostgreSQL with SQLC (type-safe SQL)
- ✅ Database migrations with Goose
- ✅ Clean architecture with proper package structure

### Coming Soon

- 🚧 Automatic RSS feed fetching
- 🚧 Posts storage and API
- 🚧 Background worker for periodic updates
- 📋 See full [Roadmap](#-roadmap)

## 🏗️ Architecture

```
rss-aggregator/
├── cmd/api/              # Application entry point
├── internal/
│   ├── auth/            # Authentication helpers
│   ├── database/        # SQLC generated code
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Auth middleware
│   └── models/          # API models & responses
└── sql/
    ├── queries/         # SQL queries for SQLC
    └── schema/          # Database migrations
```

**Design Patterns:**

- Dependency Injection (handler configs)
- Repository Pattern (SQLC data access)
- Middleware Chain (auth, CORS)
- Constructor Pattern (`NewConfig()` functions)

## 🛠️ Tech Stack

| Technology     | Purpose          |
| -------------- | ---------------- |
| **Go 1.25+**   | Backend language |
| **Chi Router** | HTTP routing     |
| **PostgreSQL** | Database         |
| **SQLC**       | Type-safe SQL    |
| **Goose**      | Migrations       |

## 🚀 Quick Start

### Prerequisites

```bash
go version    # 1.25+
psql --version   # PostgreSQL 12+
```

### Installation

```bash
# Clone repository
git clone https://github.com/mehmettalhairmak/rss-aggregator-go.git
cd rss-aggregator-go

# Install dependencies
go mod download

# Create database
createdb rss_aggregator

# Setup environment
cp .env.example .env
# Edit .env with your database URL

# Run migrations
cd sql/schema
goose postgres "$DB_URL" up
cd ../..

# Run server
go run cmd/api/main.go
```

Server starts at `http://localhost:8080` 🎉

## 📚 API Documentation

### Authentication

Protected endpoints require API Key in header:

```
Authorization: ApiKey YOUR_API_KEY
```

### Endpoints

| Method   | Endpoint                | Auth | Description         |
| -------- | ----------------------- | ---- | ------------------- |
| `GET`    | `/v1/ready`             | ❌   | Health check        |
| `POST`   | `/v1/users`             | ❌   | Create user         |
| `GET`    | `/v1/users`             | ✅   | Get user profile    |
| `POST`   | `/v1/feed`              | ✅   | Add RSS feed        |
| `GET`    | `/v1/feed`              | ❌   | List all feeds      |
| `POST`   | `/v1/feed_follows`      | ✅   | Follow a feed       |
| `GET`    | `/v1/feed_follows`      | ✅   | List followed feeds |
| `DELETE` | `/v1/feed_follows/{id}` | ✅   | Unfollow feed       |

### Example Usage

```bash
# Create user
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe"}'

# Response includes api_key - save it!

# Add a feed
curl -X POST http://localhost:8080/v1/feed \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -d '{"name": "Go Blog", "url": "https://go.dev/blog/feed.atom"}'

# List all feeds
curl http://localhost:8080/v1/feed

# Follow a feed
curl -X POST http://localhost:8080/v1/feed_follows \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -d '{"feed_id": "feed-uuid-here"}'
```

### Response Format

**Success:**

```json
{
  "id": "uuid",
  "created_at": "2025-10-05T10:00:00Z",
  "updated_at": "2025-10-05T10:00:00Z",
  ...
}
```

**Error:**

```json
{
  "error": "Descriptive error message"
}
```

## 🧪 Development

### Commands

```bash
# Run in development
go run cmd/api/main.go

# Build binary
go build -o bin/api ./cmd/api

# Run tests
go test ./...

# Generate SQLC code (after editing SQL)
sqlc generate

# Create new migration
cd sql/schema && goose create migration_name sql

# Run migrations
goose postgres "$DB_URL" up
```

### Environment Variables

Create `.env` file:

```bash
PORT=8080
DB_URL=postgres://user:pass@localhost:5432/rss_aggregator?sslmode=disable
```

## 🗺️ Roadmap

### Phase 1: JWT Authentication System

- [x] **User Authentication**
  - Login/logout endpoints
  - User registration with email & password
  - Password hashing with bcrypt
  - Email validation
- [x] **JWT Implementation**
  - Access & refresh token generation
  - Token expiration & rotation
  - Refresh token storage in database
  - Secure JWT secret management

### Phase 2: RSS Scraping

- [x] Posts table schema & migration
- [x] RSS XML parser (RSS 2.0 & Atom)
- [x] Background worker for periodic fetching
- [x] Posts API endpoints with pagination

### Phase 3: Enhanced Features

- [x] Feed URL validation before creation
- [x] Cursor-based pagination
- [x] Feed metadata (logo, description)
- [x] Better error handling & logging
- [x] Feed update priority system

### Phase 4: Production Ready

- [ ] Rate limiting (token bucket)
- [ ] Structured logging (zap/zerolog)
- [ ] Comprehensive test suite (>80% coverage)
- [ ] Docker & docker-compose
- [ ] CI/CD with GitHub Actions
- [ ] OpenAPI/Swagger documentation

### Phase 5: Advanced Features

- [ ] WebSocket for real-time updates
- [ ] Full-text search (PostgreSQL or ElasticSearch)
- [ ] Feed categories & tags
- [ ] Read/unread status tracking
- [ ] OPML import/export
- [ ] Prometheus metrics
- [ ] **Advanced Authentication**
  - OAuth2 integration (Google, GitHub)
  - Role-based access control (RBAC)
  - API key scopes & permissions
  - Two-factor authentication (2FA)
  - Session management

### Phase 6: Scaling

- [ ] Redis caching layer
- [ ] Message queue (RabbitMQ/NATS)
- [ ] Database read replicas
- [ ] Horizontal scaling support
- [ ] Performance benchmarks

## 📖 Learning Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [SQLC Documentation](https://docs.sqlc.dev/)
- [Goose Migrations](https://github.com/pressly/goose)

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

**Guidelines:**

- Follow Go conventions
- Add tests for new features
- Update documentation
- Run `go fmt` and `go vet`

## 📄 License

MIT License - see [LICENSE](LICENSE) file

## 👨‍💻 Author

**Mehmet Talha Irmak**

- GitHub: [@mehmettalhairmak](https://github.com/mehmettalhairmak)
- Project: [rss-aggregator-go](https://github.com/mehmettalhairmak/rss-aggregator-go)

---

<div align="center">

⭐ **Star this repo if you find it helpful!**

Made with ❤️ and Go

</div>
