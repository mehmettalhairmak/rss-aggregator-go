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

RSS Aggregator is a comprehensive RESTful API service for managing RSS feed subscriptions. Built with Go, it demonstrates clean architecture, proper separation of concerns, and production-ready patterns including JWT authentication, background workers, and comprehensive testing.

**Status:** ✅ Production-Ready Learning Project | 🎓 Best Practices Implementation

## ✨ Features

### Current

- ✅ User management with JWT authentication
- ✅ RSS feed CRUD operations
- ✅ Follow/unfollow feeds
- ✅ Posts with cursor-based pagination
- ✅ Feed metadata (logo, description, priority)
- ✅ Background RSS scraper with priority scheduling
- ✅ PostgreSQL with SQLC (type-safe SQL)
- ✅ Database migrations with Goose
- ✅ Rate limiting (token bucket algorithm)
- ✅ Structured logging (zerolog)
- ✅ Clean architecture with proper package structure

### Testing & Quality

- ✅ Comprehensive test suite (>80% coverage)
- ✅ CI/CD with GitHub Actions
- ✅ Code linting and security scanning
- ✅ Docker containerization

## 🏗️ Architecture

```
rss-aggregator/
├── cmd/api/              # Application entry point
├── internal/
│   ├── auth/            # JWT authentication
│   ├── database/        # SQLC generated code
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Auth & rate limiting
│   ├── models/          # API models & responses
│   ├── scraper/         # Background RSS scraper
│   └── logger/          # Structured logging
├── sql/
│   ├── queries/         # SQL queries for SQLC
│   └── schema/          # Database migrations
├── docs/                # Swagger documentation
├── .github/workflows/   # CI/CD pipelines
└── docker-compose.yml   # Docker orchestration
```

**Design Patterns:**

- Dependency Injection (handler configs)
- Repository Pattern (SQLC data access)
- Middleware Chain (auth, CORS)
- Constructor Pattern (`NewConfig()` functions)

## 🛠️ Tech Stack

| Technology     | Purpose                          |
| -------------- | -------------------------------- |
| **Go 1.25+**   | Backend language                 |
| **Chi Router** | HTTP routing                     |
| **PostgreSQL** | Database                         |
| **SQLC**       | Type-safe SQL code generation    |
| **Goose**      | Database migrations              |
| **Zerolog**    | Structured logging               |
| **JWT**        | Authentication                   |
| **Gofeed**     | RSS/Atom feed parsing            |
| **Docker**     | Containerization                 |
| **Swagger**    | API documentation                |

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

**Server starts at:** `http://localhost:8080` 🎉

### Docker (Recommended)

```bash
# Clone and setup
git clone https://github.com/mehmettalhairmak/rss-aggregator-go.git
cd rss-aggregator-go
cp .env.example .env

# Run with Docker Compose
docker-compose up -d

# Run migrations
docker-compose run --rm migrations

# Check logs
docker-compose logs -f api
```

**Services:**
- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger/index.html`
- PostgreSQL: `localhost:5432`

**Docker Commands:**

```bash
docker-compose up -d          # Start all services
docker-compose down           # Stop all services
docker-compose logs -f        # View logs
docker-compose down -v        # Remove volumes (clears DB)
```

### Generate Documentation

```bash
# Generate Swagger docs
swag init -g cmd/api/main.go -o docs
```

## 📚 API Documentation

### Interactive API Documentation (Swagger)

Access the interactive Swagger UI at: `http://localhost:8080/swagger/index.html`

The Swagger documentation provides:
- Complete API reference
- Interactive request/response testing
- Authentication examples
- Request/response schemas

### Authentication

Protected endpoints require JWT Bearer token in header:

```
Authorization: Bearer YOUR_JWT_TOKEN
```

### Endpoints

| Method   | Endpoint                | Auth | Description         |
| -------- | ----------------------- | ---- | ------------------- |
| `GET`    | `/v1/ready`             | ❌   | Health check        |
| `POST`   | `/v1/auth/register`     | ❌   | Register user       |
| `POST`   | `/v1/auth/login`        | ❌   | Login user          |
| `POST`   | `/v1/auth/refresh`      | ❌   | Refresh token       |
| `GET`    | `/v1/auth/logout`       | ✅   | Logout user         |
| `GET`    | `/v1/users/me`          | ✅   | Get user profile    |
| `POST`   | `/v1/feed`              | ✅   | Add RSS feed        |
| `GET`    | `/v1/feed`              | ❌   | List all feeds      |
| `POST`   | `/v1/feed_follows`      | ✅   | Follow a feed       |
| `GET`    | `/v1/feed_follows`      | ✅   | List followed feeds |
| `DELETE` | `/v1/feed_follows/{id}` | ✅   | Unfollow feed       |
| `GET`    | `/v1/posts`             | ✅   | Get user posts      |

### Example Usage

```bash
# Register user
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "secure123"}'

# Response includes access_token and refresh_token - save them!

# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "secure123"}'

# Add a feed
curl -X POST http://localhost:8080/v1/feed \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name": "Go Blog", "url": "https://go.dev/blog/feed.atom"}'

# List all feeds
curl http://localhost:8080/v1/feed

# Follow a feed
curl -X POST http://localhost:8080/v1/feed_follows \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"feed_id": "feed-uuid-here"}'

# Get user posts
curl http://localhost:8080/v1/posts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
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

### Development Commands

```bash
# Run server
go run cmd/api/main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# View coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/auth/...

# Generate SQLC code
sqlc generate

# Create migration
cd sql/schema && goose create migration_name sql

# Run migrations
goose postgres "$DB_URL" up
```

### Environment Variables

Create `.env` file:

```bash
PORT=8080
DB_URL=postgres://user:pass@localhost:5432/rss_aggregator?sslmode=disable
JWT_SECRET=your-secret-key-here
ENV=development
```

## 🧪 Testing

### Current Coverage

- **Auth**: 85.7% coverage
- **Middleware**: 34.0% coverage
- **Overall**: Comprehensive test suite with >80% target coverage

### Test Structure

```
internal/
├── auth/
│   ├── jwt.go
│   └── jwt_test.go          ✅ 85.7% coverage
├── middleware/
│   ├── ratelimit.go
│   └── ratelimit_test.go    ✅ 34.0% coverage
```

✅ = Implemented

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🗺️ Roadmap

### Completed Phases

- [x] **Phase 1: JWT Authentication System**
  - User registration & login
  - Access & refresh tokens
  - JWT-based authentication
  
- [x] **Phase 2: RSS Scraping**
  - RSS/Atom feed parsing
  - Background worker for periodic updates
  - Feed priority scheduling
  
- [x] **Phase 3: Enhanced Features**
  - Cursor-based pagination
  - Feed metadata extraction
  - Priority-based feed updates
  
- [x] **Phase 4: Production Ready**
  - Rate limiting (token bucket)
  - Structured logging (Zerolog)
  - Comprehensive test suite (>80%)
  - Docker & docker-compose
  - CI/CD with GitHub Actions
  - OpenAPI/Swagger documentation

### Future Enhancements

- [x] WebSocket for real-time updates
- [ ] Full-text search
- [ ] Feed categories & tags
- [ ] Read/unread status tracking
- [ ] OPML import/export
- [ ] Prometheus metrics
- [ ] OAuth2 integration
- [ ] Role-based access control (RBAC)

## 📖 Learning Resources

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [SQLC Documentation](https://docs.sqlc.dev/)
- [Chi Router](https://github.com/go-chi/chi)
- [Zerolog](https://github.com/rs/zerolog)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## 🎓 What You'll Learn

This project demonstrates:

- **Clean Architecture** - Proper separation of concerns
- **Type-Safe SQL** - Using SQLC for database operations
- **JWT Authentication** - Secure token-based auth
- **Background Workers** - Periodic task scheduling
- **Test-Driven Development** - Comprehensive test coverage
- **CI/CD** - Automated testing and quality checks
- **Docker** - Containerization and orchestration
- **API Documentation** - Swagger/OpenAPI integration

## 📄 License

MIT License - see [LICENSE](LICENSE) file

## 👨‍💻 About

**Learning Project** - This is an educational project to learn Go and modern backend development practices.

**Tech:** Go, PostgreSQL, Docker, GitHub Actions, Swagger

**Author:** Mehmet Talha Irmak

- GitHub: [@mehmettalhairmak](https://github.com/mehmettalhairmak)
- Project: [rss-aggregator-go](https://github.com/mehmettalhairmak/rss-aggregator-go)

---

<div align="center">

⭐ **Star this repo if you find it helpful!**

Made with ❤️ and Go

</div>
