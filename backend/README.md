# Learn Platform Backend

Go backend API for Learn Platform - a full-stack learning platform.

## 🏗️ Architecture

```
Clean Architecture Pattern:
├── Handlers      (HTTP Layer)
├── Services      (Business Logic)
├── Repositories  (Data Access)
└── Models        (Data Structures)
```

## 📦 Dependencies

### Core
- **Echo**: Fast HTTP web framework
- **GORM**: Object-relational mapping for databases
- **PostgreSQL Driver**: Database driver for PostgreSQL

### Authentication
- **JWT (golang-jwt)**: JSON Web Tokens for authentication
- **bcrypt**: Password hashing

### Security
- **golang/crypto**: Cryptographic functions

## 🚀 Quick Start

### Prerequisites
```bash
# Go 1.22 or higher
go version

# PostgreSQL 16+
psql --version
```

### Installation

```bash
# Clone and navigate
cd apps/backend   # dari root repo; atau cd backend dari folder apps/

# Create environment file
cp .env.example .env

# Edit .env with your settings
# Important: Change JWT_SECRET and DATABASE_URL

# Install dependencies
go mod download

# Run
go run main.go

# Or with hot reload (requires air)
# go install github.com/cosmtrek/air@latest
air
```

## 📁 Project Structure

```
backend/
├── main.go                 # Entry point
├── go.mod                  # Dependencies
├── Makefile                # Build commands
├── Dockerfile              # Container image
├── .env.example            # Environment template
│
├── internal/               # Private application code
│   ├── config/
│   │   └── config.go       # Configuration loading
│   │
│   ├── model/              # Data models
│   │   ├── user.go         # User model
│   │   └── auth.go         # Auth models
│   │
│   ├── repository/         # Data access layer
│   │   ├── user_repository.go
│   │   └── auth_repository.go
│   │
│   ├── service/            # Business logic
│   │   ├── user_service.go
│   │   └── auth_service.go
│   │
│   └── handler/            # HTTP handlers
│       ├── auth_handler.go
│       └── user_handler.go
│
└── pkg/                    # Public utilities
    ├── database/           # Database setup
    │   └── database.go
    └── logger/             # Logging
        └── logger.go
```

## 🔑 API Endpoints

### Health Check
```
GET /health
```

### Authentication
```
POST   /api/v1/auth/register     # User registration
POST   /api/v1/auth/login        # User login
POST   /api/v1/auth/refresh      # Refresh access token
```

### Users (Protected)
```
GET    /api/v1/users             # List all users
GET    /api/v1/users/:id         # Get specific user
GET    /api/v1/users/me          # Get current user
PUT    /api/v1/users/:id         # Update user
DELETE /api/v1/users/:id         # Delete user (admin only)
```

## 🔒 Security Features

### Authentication
- JWT with 15-minute expiration
- Refresh tokens with 7-day expiration
- Password hashing with bcrypt
- Secure token validation

### Authorization
- Role-based access control (RBAC)
- User can only update their own profile
- Admin-only delete operations

### Middleware
- CORS protection
- Rate limiting (100 requests/minute)
- Security headers
- Request ID tracking
- Gzip compression
- Request logging
- Panic recovery

### Input Validation
- Email format validation
- Password minimum length (8 chars)
- Required field checks
- SQL injection prevention (parameterized queries)

### Error Handling
- Generic error messages (no information leakage)
- Proper HTTP status codes
- Context-aware error logging

## 🗄️ Database

### Connection
- PostgreSQL with GORM
- Connection pooling enabled
- Auto-migration on startup

### Schema
```sql
users:
  id: UUID (Primary Key)
  email: VARCHAR (Unique)
  password: VARCHAR (Hashed)
  name: VARCHAR
  role: VARCHAR (default: 'user')
  is_active: BOOLEAN (default: true)
  created_at: TIMESTAMP
  updated_at: TIMESTAMP
  deleted_at: TIMESTAMP (Soft delete)
```

## 🛠️ Development Commands

```bash
# Build
make build

# Run
make run

# Hot reload development
make dev

# Testing
make test
make test-cover

# Code quality
make lint
make fmt

# Dependencies
make deps

# Docker
make docker-build
make docker-up
make docker-down
make docker-logs
```

## 📊 Key Implementation Details

### JWT Tokens
- **Access Token**: 15 minutes validity
- **Refresh Token**: 7 days validity
- **Algorithm**: HS256 (HMAC-SHA256)
- **Claims**: User ID, Email, Role, Expiration

### Password Security
- **Hashing**: bcrypt with cost factor 10
- **Verification**: Compare hash with input
- **Never** logged or returned in responses

### User Roles
- **user**: Regular user (can view own profile)
- **admin**: Administrator (can manage all users)

### Error Responses

```json
{
  "message": "Generic error message"
}
```

HTTP Status Codes:
- 200: OK
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## 🧪 Testing

```bash
# Run all tests with verbose output
make test

# With race condition detection
go test -race ./...

# With coverage report
make test-cover

# View coverage in browser
open coverage.html
```

### Test Structure
```
./internal/
├── handler/handler_test.go
├── service/service_test.go
└── repository/repository_test.go
```

## 🔄 Request Flow

```
HTTP Request
    ↓
Middleware (CORS, logging, rate limit, etc.)
    ↓
Handler (validate input, call service)
    ↓
Service (business logic)
    ↓
Repository (database access)
    ↓
Database
    ↓
Response (JSON)
```

## 🚀 Deployment

### Docker

Build image:
```bash
make docker-build
```

Run container:
```bash
docker run -p 8080:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/db \
  -e JWT_SECRET=your-secret \
  learn-platform-backend:latest
```

### Environment Variables for Production

```env
ENVIRONMENT=production
PORT=8080
DATABASE_URL=postgres://prod_user:prod_pass@prod_host:5432/learn_platform_prod
JWT_SECRET=your-very-secure-random-string
CORS_ORIGIN=https://yourdomain.com
```

## 📈 Performance Optimization

- [x] Connection pooling
- [x] Database indexes on email (unique)
- [x] Prepared statements (via GORM)
- [x] Gzip compression
- [x] Request caching headers (future)
- [ ] Redis caching (ready for integration)
- [ ] Batch operations
- [ ] Query optimization

## 🐛 Common Issues

### Database Connection Error
```
Error: dial tcp localhost:5432: connection refused
```
Solution: Ensure PostgreSQL is running and DATABASE_URL is correct

### JWT Secret Error
```
Error: invalid signing method
```
Solution: Ensure JWT_SECRET is set and valid

### Port Already in Use
```bash
lsof -i :8080
kill -9 <PID>
```

## 📚 Additional Resources

- [Echo Documentation](https://echo.labstack.com/)
- [GORM Documentation](https://gorm.io/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc7519)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 🔗 Related Documentation

- [Frontend README](../frontend/README.md)
- [Main Project README](../README.md)

---

**Technology Stack**: Go 1.22 | Echo | GORM | PostgreSQL | JWT | Docker

**Security First**: All endpoints protected, inputs validated, errors generic, passwords hashed, tokens refreshable.
