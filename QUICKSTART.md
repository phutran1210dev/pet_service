# 🚀 Quick Start Guide

Get your Pet Service API up and running in 5 minutes!

## Prerequisites

- Docker & Docker Compose installed
- OR Go 1.24+ installed

## Option 1: Docker (Recommended) ⭐

### 1. Clone and Start

```bash
cd pet-service-golang

# Start all services
./start.sh

# Or manually
docker-compose up -d
```

### 2. Verify

```bash
# Check if services are running
docker-compose ps

# Test API
curl http://localhost:8001/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "Pet Service API is running"
}
```

### 3. Test the API

```bash
# Run interactive tests
./test-api.sh

# Or run all tests automatically
./test-api.sh all
```

### 4. Access Services

- **API**: http://localhost:8001
- **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin)
- **PostgreSQL**: localhost:5432

## Option 2: Local Development

### 1. Setup

```bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit .env with your settings
nano .env
```

### 2. Start Dependencies

```bash
# Start only PostgreSQL and MinIO
docker-compose up -d postgres minio
```

### 3. Run Application

```bash
# Run directly
go run main.go

# Or with hot reload
make dev
```

## 🧪 Testing the API

### Using the test script:

```bash
./test-api.sh
```

### Manual testing:

```bash
# 1. Register a user
curl -X POST http://localhost:8001/api/v1/user \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "0123456789",
    "gender": true,
    "password": "password123"
  }'

# Response includes: id, email, first_name, last_name, phone, is_admin, roles, permissions

# 2. Login
curl -X POST http://localhost:8001/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'

# Save the access_token from response

# 3. Get your profile
curl -X GET http://localhost:8001/api/v1/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 4. Create a pet
curl -X POST http://localhost:8001/api/v1/pet \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Buddy",
    "gender": true,
    "date_of_birth": "2020-01-15",
    "breed": "Golden Retriever",
    "description": "Friendly dog",
    "type": "dog"
  }'

# 5. Get all pets
curl -X GET "http://localhost:8001/api/v1/pets?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Response structure:
# {
#   "data": [...],
#   "meta": {
#     "total_items": 10,
#     "total_pages": 2,
#     "page": 1,
#     "page_size": 10
#   }
# }
```

## 📊 Monitoring

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f pet-service

# With make
make docker-logs
```

### Check Status

```bash
# Service status
docker-compose ps

# Resource usage
docker stats
```

## 🛑 Stopping

```bash
# Interactive stop
./stop.sh

# Or manually
docker-compose down

# Stop and remove all data
docker-compose down -v
```

## 🔧 Troubleshooting

### API not responding?

```bash
# Check logs
docker-compose logs pet-service

# Restart service
docker-compose restart pet-service
```

### Database connection error?

```bash
# Check PostgreSQL
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

### Port already in use?

```bash
# Find what's using the port
lsof -i :8001

# Stop the process or change port in docker-compose.yml
```

### Fresh start

```bash
# Stop everything and clean up
docker-compose down -v
docker system prune -f

# Start again
./start.sh
```

## 📚 Next Steps

1. **Read the full documentation**: [README.md](README.md)
2. **Docker details**: [DOCKER.md](DOCKER.md)
3. **API Documentation**: http://localhost:8001/docs (when DEBUG=true)
4. **Customize configuration**: Edit `.env` or `docker-compose.yml`

## 🎯 Common Tasks

### Development Mode with Hot Reload

```bash
docker-compose -f docker-compose.dev.yml up
# OR
make docker-dev-up
```

### Production Build

```bash
make prod-build
# OR
docker-compose up -d --build
```

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

## 💡 Tips

- Use `./test-api.sh` for quick API testing
- Use `make help` to see all available commands
- Check `docker-compose ps` to verify services are healthy
- Access MinIO console at http://localhost:9001 to manage files
- PostgreSQL runs on localhost:5432 for direct access

## 🆘 Getting Help

- Check logs: `docker-compose logs -f`
- Verify health: `curl http://localhost:8001/health`
- Restart services: `docker-compose restart`
- Clean restart: `docker-compose down -v && ./start.sh`

---

**Ready to go!** 🚀 Your Pet Service API is now running!
