# 🚀 Bong Aquino Backend Services

> **High-performance microservices architecture built with Go, featuring API Gateway, comprehensive monitoring, and enterprise-grade security.**

[![Go](https://img.shields.io/badge/Go-1.22.3-00ADD8?logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?logo=docker)](https://docker.com)
[![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?logo=go)](https://gin-gonic.com)
[![MongoDB](https://img.shields.io/badge/MongoDB-Database-47A248?logo=mongodb)](https://mongodb.com)

## ✨ Features

- **🏗️ Microservices Architecture** - Scalable, maintainable service design
- **🛡️ API Gateway Integration** - Centralized request management with Tyk
- **🔒 Enterprise Security** - JWT authentication, MFA, role-based access
- **📊 Comprehensive Monitoring** - ELK stack integration and health checks
- **🚀 High Performance** - Optimized Go services with Redis caching
- **🔄 Auto-deployment** - CI/CD pipeline with Docker containerization

## 🚀 Quick Start

### Prerequisites
- **Go** 1.22.3+
- **Docker** & Docker Compose
- **Make** (optional, for shortcuts)

### One-Command Setup
```bash
# Clone and start everything
git clone https://github.com/bongaquino/go-backend.git
cd go-backend
./scripts/start.sh
```

### Manual Setup
```bash
# 1. Create Docker network
docker network create bongaquino-network

# 2. Start all services
docker-compose up -d

# 3. Verify services
curl http://localhost:3000/health
```

## 🏗️ Architecture Overview

### Service Ecosystem
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   🌐 Client     │────│  🛡️ API Gateway │────│  🚀 Go Server   │
│   (Frontend)    │    │   (Tyk:8080)    │    │   (Port:3000)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                    ┌─────────────────┐    ┌─────────────────┐
                    │  📊 Monitoring  │    │  💾 Data Layer  │
                    │  ELK Stack      │    │  MongoDB+Redis  │
                    └─────────────────┘    └─────────────────┘
```

### Core Services
| Service | Port | Purpose | Health Check |
|---------|------|---------|--------------|
| **🚀 Go Server** | 3000 | Main API backend | `/health` |
| **🛡️ Tyk Gateway** | 8080 | API management | `/hello` |
| **💾 MongoDB** | 27017 | Primary database | Auto-monitored |
| **⚡ Redis** | 6379 | Caching & sessions | Auto-monitored |
| **🔍 Elasticsearch** | 9200 | Search & logging | `/_cluster/health` |
| **📊 Kibana** | 5601 | Log visualization | `/app/home` |

## 📊 Monitoring & Management

### Service Dashboards
| Interface | URL | Purpose |
|-----------|-----|---------|
| **🖥️ Mongo Express** | [localhost:8082](http://localhost:8082) | Database management |
| **⚡ Redis Commander** | [localhost:8083](http://localhost:8083) | Cache monitoring |
| **📊 Kibana** | [localhost:5601](http://localhost:5601) | Log analytics |
| **🛡️ Tyk Dashboard** | [localhost:8080](http://localhost:8080) | API management |

### Health Monitoring
```bash
# Check all services
curl http://localhost:3000/health

# Individual service logs
docker-compose logs -f server
docker-compose logs -f mongo
docker-compose logs -f redis
```

## 🔧 Development Environment

### Local Development
```bash
# Hot reload development
cd server
go run main.go

# Or with Air (recommended)
air

# Run with environment
cp .env.example .env
# Edit .env with your settings
go run main.go
```

### Environment Configuration
```env
# Database
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=bongaquino_db

# Redis Cache
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Authentication
JWT_SECRET=your-super-secret-key
JWT_EXPIRY=24h

# Email Service
POSTMARK_API_KEY=your-postmark-key
FROM_EMAIL=noreply@example.com

# IPFS Integration
IPFS_HOST=localhost
IPFS_PORT=5001
```

## 🛡️ Security Features

### Authentication & Authorization
- **🔐 JWT Tokens** - Secure stateless authentication
- **🔑 Multi-Factor Authentication** - TOTP-based 2FA
- **👤 Role-Based Access Control** - Granular permissions
- **🔒 Password Security** - Bcrypt hashing
- **⏰ Session Management** - Configurable expiry

### API Security
```go
// Protected route example
router.Use(middleware.AuthRequired())
router.Use(middleware.RoleRequired("admin"))
router.GET("/admin/users", adminController.ListUsers)
```

### Security Headers
- **🛡️ CORS** - Configurable cross-origin requests
- **🔒 Rate Limiting** - Request throttling
- **📋 Input Validation** - Request sanitization
- **🔐 Encryption** - Data encryption at rest

## 📁 Project Structure

```
server/
├── 🚀 main.go                    # Application entry point
├── 📁 app/                       # Core application logic
│   ├── 🎮 controller/            # HTTP request handlers
│   │   ├── auth/                 # Authentication endpoints
│   │   ├── users/                # User management
│   │   ├── files/                # File operations
│   │   └── admin/                # Admin functions
│   ├── 🧩 service/               # Business logic layer
│   ├── 💾 repository/            # Data access layer
│   ├── 🔌 provider/              # External integrations
│   └── 🛡️ middleware/           # Request middleware
├── ⚙️ config/                    # Configuration management
├── 🗄️ database/                  # DB migrations & seeders
└── 🚀 start/                     # Application bootstrap
```

## 🔌 API Endpoints

### Authentication
```http
POST   /auth/register          # User registration
POST   /auth/login             # User login
POST   /auth/refresh           # Token refresh
POST   /auth/logout            # User logout
POST   /auth/forgot-password   # Password reset
```

### User Management
```http
GET    /users/profile          # Get user profile
PUT    /users/profile          # Update profile
POST   /users/change-password  # Change password
GET    /users/settings         # User preferences
```

### File Operations
```http
POST   /files/upload           # Upload file
GET    /files/:id              # Get file details
PUT    /files/:id              # Update file
DELETE /files/:id              # Delete file
POST   /files/:id/share        # Share file
```

### Admin Functions
```http
GET    /admin/users            # List all users
POST   /admin/users            # Create user
PUT    /admin/users/:id        # Update user
DELETE /admin/users/:id        # Delete user
```

## 🚀 Deployment

### Docker Production
```bash
# Build production image
docker build -t bongaquino/backend:latest .

# Run production stack
docker-compose -f docker-compose.prod.yml up -d

# Scale services
docker-compose up -d --scale server=3
```

### CI/CD Pipeline
```yaml
# GitHub Actions workflow
name: Deploy Backend
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build & Deploy
        run: |
          docker build -t backend .
          docker push registry/backend:latest
```

## 📊 Performance & Monitoring

### Metrics
- **⚡ Response Time**: < 100ms average
- **🔄 Throughput**: 1000+ requests/second
- **💾 Memory Usage**: < 512MB per instance
- **🚀 CPU Usage**: < 50% under normal load

### Logging
```go
// Structured logging
logger.Info("User authenticated", 
    zap.String("user_id", userID),
    zap.String("ip", clientIP),
    zap.Duration("duration", elapsed),
)
```

## 🧪 Testing

### Running Tests
```bash
# Unit tests
go test ./...

# Integration tests
go test ./app/... -tags=integration

# Load testing
hey -n 1000 -c 10 http://localhost:3000/health
```

### Test Coverage
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🛠️ Management Scripts

### Available Scripts
```bash
# Start all services
./scripts/start.sh

# Stop all services
./scripts/stop.sh

# Restart services
./scripts/restart.sh

# Complete rebuild
./scripts/rebuild.sh
```

### Database Management
```bash
# Run migrations
go run database/migration.go

# Seed test data
go run database/seeder.go

# Backup database
mongodump --db bongaquino_db
```

## 🤝 Contributing

### Development Workflow
1. **🍴 Fork** the repository
2. **🌿 Create** feature branch: `git checkout -b feature/new-endpoint`
3. **💻 Develop** with tests: `go test ./...`
4. **📝 Document** API changes
5. **🚀 Submit** pull request

### Code Standards
- **📝 Go fmt** - Standard formatting
- **📋 Linting** - golangci-lint
- **🧪 Testing** - 80%+ coverage
- **📖 Documentation** - Godoc comments

## 📞 Support

- **📖 API Documentation**: [Swagger UI](http://localhost:3000/swagger)
- **🐛 Bug Reports**: [GitHub Issues](https://github.com/bongaquino/go-backend/issues)
- **💬 Support**: [Discord Community](https://discord.gg/bongaquino)
- **📧 Contact**: admin@example.com

---

<div align="center">

**Built with ❤️ by [Bong Aquino](https://github.com/bongaquino)**

*Enterprise Go Backend | Secure • Scalable • Performant*

</div>