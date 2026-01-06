# Pet Service API - Golang Gin

This is a Golang Gin port of the Pet Service API, originally written in Python FastAPI.

## Design

- **Figma Design**: [Pet Care UI](https://www.figma.com/design/axle4eRnxP9TtZSS8i44mq/Pet-Care-UI?node-id=0-1&t=javBl23WbAN8qt2O-1)

## Features

- User authentication with JWT
- User management (register, login, profile, change password)
- Pet management (CRUD operations, life events, image uploads)
- Comment system for pets
- Appointment booking with email scheduling
- MinIO integration for file storage
- PostgreSQL database with GORM
- Role-based access control

## Technology Stack

- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **File Storage**: MinIO
- **Task Scheduling**: robfig/cron
- **Password Hashing**: bcrypt

## Project Structure

```bash
pet-service-golang/
├── config/          # Configuration management
├── database/        # Database connection and setup
├── dto/             # Data Transfer Objects (request/response schemas)
├── handler/         # HTTP handlers (controllers)
├── middleware/      # Middleware (auth, logging, CORS)
├── models/          # GORM models
├── repository/      # Database repository layer
├── routes/          # API route definitions
├── scheduler/       # Task scheduler
├── service/         # Business logic layer
├── storage/         # MinIO file storage
├── utils/           # Utility functions and constants
├── main.go          # Application entry point
├── go.mod           # Go module dependencies
└── .env.example     # Environment variables template
```

## Setup

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 12+
- MinIO (for file storage)

**OR** just use Docker (recommended):

- Docker
- Docker Compose

### Quick Start with Docker (Recommended)

The easiest way to get started:

```bash
# Make scripts executable (first time only)
chmod +x start.sh stop.sh

# Start the application
./start.sh

# Or manually
docker-compose up -d
```

That's it! The API will be available at <http://localhost:8001>

To stop:

```bash
./stop.sh
# Or manually
docker-compose down
```

See [DOCKER.md](DOCKER.md) for detailed Docker documentation.

### Installation

1. Clone the repository:

```bash
cd pet-service-golang
```

1. Copy `.env.example` to `.env` and configure your environment variables:

```bash
cp .env.example .env
```

1. Edit `.env` with your configuration:

```env
PROJECT_NAME=Pet Service API
DEBUG=true
SECRET_KEY=your-secret-key-here
TIME_ZONE=Asia/Ho_Chi_Minh

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=pet_service

SERVER_PORT=8001

MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET=pet-service
```

1. Install dependencies:

```bash
go mod download
```

1. Run the application:

```bash
go run main.go
```

The server will start on `http://localhost:8001`

## Docker Deployment

### Quick Start with Docker Compose

1. Start all services (PostgreSQL, MinIO, and API):

```bash
docker-compose up -d
```

1. Check logs:

```bash
docker-compose logs -f pet-service
```

1. Stop all services:

```bash
docker-compose down
```

### Development with Hot Reload

For development with hot reload using Air:

```bash
# Start development environment
docker-compose -f docker-compose.dev.yml up

# Or use Makefile
make docker-dev-up
```

### Using Makefile

The project includes a Makefile for common tasks:

```bash
# Show all available commands
make help

# Development
make dev                 # Run with hot reload locally
make docker-dev-up      # Run with Docker and hot reload

# Production
make docker-up          # Start all services
make docker-down        # Stop all services
make docker-logs        # View logs
make docker-restart     # Restart services

# Build
make build              # Build Go binary
make docker-build       # Build Docker image

# Testing & Quality
make test               # Run tests
make fmt                # Format code
make vet                # Run go vet
```

### Services

When running with docker-compose, the following services are available:

- **API**: <http://localhost:8001>
- **PostgreSQL**: localhost:5432
- **MinIO**: <http://localhost:9000> (API), <http://localhost:9001> (Console)
  - Default credentials: minioadmin / minioadmin

### Environment Variables in Docker

The docker-compose.yml includes all necessary environment variables. For production, create a `.env` file or modify the environment section in docker-compose.yml with secure values.

## API Endpoints

### Authentication

- `POST /api/v1/login` - Login
- `POST /api/v1/user` - Register new user
- `POST /api/v1/logout` - Logout (requires auth)
- `GET /api/v1/me` - Get current user info (requires auth)

### User Management

- `GET /api/v1/users` - Get all users (requires auth)
- `PATCH /api/v1/users/change-password` - Change password (requires auth)

### Pet Management

- `POST /api/v1/pet` - Create pet (requires auth)
- `GET /api/v1/pets` - Get all pets with pagination (requires auth)
  - Returns data with meta object containing: total_items, total_pages, page, page_size
- `GET /api/v1/pet/:id` - Get pet details (requires auth)
- `POST /api/v1/pet/life-event` - Create pet life event (requires auth)
- `POST /api/v1/pet/:pet_id/images` - Upload pet avatar (requires auth)
- `POST /api/v1/pet/:pet_id/gallery` - Upload pet gallery images (requires auth)

### Comments

- `POST /api/v1/post/:pet_id/comment` - Create comment (requires auth)
- `PATCH /api/v1/post/:pet_id/comment/:comment_id` - Edit comment (requires auth)
- `GET /api/v1/post/:pet_id/comments` - Get comments (requires auth)

### Appointments

- `POST /api/v1/appointment/register` - Register appointment (requires auth)

### Health Check

- `GET /health` - API health check

## Authentication

The API uses JWT Bearer token authentication. To access protected endpoints:

1. Login or register to get an access token
2. Include the token in the Authorization header:

```
Authorization: Bearer <your-token>
```

## Development

### Running with hot reload

Install air for hot reloading:

```bash
go install github.com/air-verse/air@latest
```

Create `.air.toml` configuration file, then run:

```bash
air
```

### Database Migration

The application uses GORM for database management. Models are automatically synced on startup in debug mode.

For production, consider using a proper migration tool like:

- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)

### Key Changes

1. **Type Safety**: Golang's static typing provides compile-time type checking
2. **Performance**: Generally faster execution due to compiled nature
3. **Concurrency**: Uses goroutines for concurrent operations
4. **Error Handling**: Explicit error returns instead of exceptions
5. **Middleware**: Gin middleware system vs Starlette middleware
6. **Response Structure**: Standardized pagination with meta field for better API design

### Features Parity

| Feature | Status | Description |
|---------|--------|-------------|
| User Authentication & Authorization | ✅ | JWT-based auth with role-based access control |
| Pet CRUD Operations | ✅ | Create, read, update, delete pets with owner info |
| File Upload to MinIO | ✅ | Image storage for pet avatars and galleries |
| Comment System | ✅ | User comments on pet posts |
| Appointment Booking | ✅ | Schedule appointments with email notifications |
| JWT Token Management | ✅ | Access and refresh token handling |
| Role-based Permissions | ✅ | User roles and permission management |
| Pagination Support | ✅ | Standardized pagination with meta field |

## Production Deployment

1. Set `DEBUG=false` in `.env`
2. Use a proper secret key
3. Configure PostgreSQL with proper connection pooling
4. Set up SSL/TLS for MinIO
5. Use a reverse proxy (nginx/traefik) with HTTPS
6. Consider using Docker for deployment

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please follow the Go coding standards and include tests for new features.
