# Docker Setup Guide for Pet Service API

## 🐳 Docker Compose Configuration

This project includes complete Docker setup with:
- **PostgreSQL 15**: Database
- **MinIO**: Object storage for images
- **Pet Service API**: Golang Gin application

## 📦 Quick Start

### Production Mode

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Remove all data (including volumes)
docker-compose down -v
```

### Development Mode (with Hot Reload)

```bash
# Start development environment
docker-compose -f docker-compose.dev.yml up

# Stop development environment
docker-compose -f docker-compose.dev.yml down
```

## 🛠️ Using Makefile

```bash
# Production
make docker-up          # Start production services
make docker-down        # Stop services
make docker-logs        # View logs
make docker-restart     # Restart services
make docker-clean       # Remove containers and volumes

# Development
make docker-dev-up      # Start with hot reload
make docker-dev-down    # Stop dev environment

# Building
make docker-build       # Build Docker images
make prod-build        # Build production binary
```

## 🌐 Access Services

After running `docker-compose up -d`:

| Service | URL | Credentials |
|---------|-----|-------------|
| API | http://localhost:8001 | - |
| API Health | http://localhost:8001/health | - |
| PostgreSQL | localhost:5432 | postgres/postgres |
| MinIO API | http://localhost:9000 | minioadmin/minioadmin |
| MinIO Console | http://localhost:9001 | minioadmin/minioadmin |

## 📝 Configuration

### Environment Variables

All environment variables are defined in `docker-compose.yml`:

```yaml
environment:
  PROJECT_NAME: Pet Service API
  DEBUG: "true"
  SECRET_KEY: your-secret-key-change-in-production
  TIME_ZONE: Asia/Ho_Chi_Minh
  
  POSTGRES_HOST: postgres
  POSTGRES_PORT: 5432
  POSTGRES_USER: postgres
  POSTGRES_PASSWORD: postgres
  POSTGRES_DB: pet_service
  
  SERVER_PORT: 8001
  
  MINIO_ENDPOINT: minio:9000
  MINIO_ACCESS_KEY: minioadmin
  MINIO_SECRET_KEY: minioadmin
  MINIO_USE_SSL: "false"
  MINIO_BUCKET: pet-service
```

### Customizing Configuration

**Option 1**: Edit `docker-compose.yml` directly

**Option 2**: Create `.env` file in project root:
```env
SECRET_KEY=your-production-secret-key
POSTGRES_PASSWORD=your-secure-password
MINIO_ROOT_PASSWORD=your-minio-password
```

Then reference in docker-compose.yml:
```yaml
environment:
  SECRET_KEY: ${SECRET_KEY}
  POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
```

## 🔧 Troubleshooting

### Check Service Health

```bash
# Check all services
docker-compose ps

# Check specific service logs
docker-compose logs postgres
docker-compose logs minio
docker-compose logs pet-service

# Follow logs in real-time
docker-compose logs -f pet-service
```

### Database Connection Issues

```bash
# Restart PostgreSQL
docker-compose restart postgres

# Check PostgreSQL logs
docker-compose logs postgres

# Connect to PostgreSQL directly
docker exec -it pet-service-postgres psql -U postgres -d pet_service
```

### MinIO Issues

```bash
# Restart MinIO
docker-compose restart minio

# Check MinIO logs
docker-compose logs minio

# Access MinIO console
# Open http://localhost:9001 in browser
```

### API Not Starting

```bash
# Check API logs
docker-compose logs pet-service

# Rebuild API image
docker-compose up -d --build pet-service

# Check API health
curl http://localhost:8001/health
```

### Port Already in Use

If you see "port already allocated" error:

```bash
# Option 1: Stop conflicting services
sudo lsof -i :8001  # Find process using port
sudo kill -9 <PID>  # Kill the process

# Option 2: Change port in docker-compose.yml
ports:
  - "8002:8001"  # Use different host port
```

## 🗄️ Data Persistence

### Volumes

Docker compose creates two named volumes:

```bash
# List volumes
docker volume ls | grep pet-service

# Inspect volume
docker volume inspect pet-service-golang_postgres_data
docker volume inspect pet-service-golang_minio_data

# Backup volume
docker run --rm -v pet-service-golang_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres-backup.tar.gz /data

# Restore volume
docker run --rm -v pet-service-golang_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres-backup.tar.gz
```

### Clean Up

```bash
# Remove all containers and volumes
docker-compose down -v

# Remove unused images
docker image prune -a

# Complete cleanup
make docker-clean
```

## 🚀 Production Deployment

### 1. Update Configuration

```bash
# Edit docker-compose.yml
environment:
  DEBUG: "false"
  SECRET_KEY: "$(openssl rand -hex 32)"
  POSTGRES_PASSWORD: "strong-password-here"
```

### 2. Use Production Build

```bash
# Build optimized image
docker-compose build --no-cache

# Start services
docker-compose up -d
```

### 3. Enable HTTPS

Add nginx reverse proxy:

```yaml
# docker-compose.yml
nginx:
  image: nginx:alpine
  ports:
    - "80:80"
    - "443:443"
  volumes:
    - ./nginx.conf:/etc/nginx/nginx.conf
    - ./ssl:/etc/nginx/ssl
  depends_on:
    - pet-service
```

### 4. Monitoring

```bash
# Check resource usage
docker stats

# View all logs
docker-compose logs --tail=100

# Monitor specific service
watch docker-compose ps
```

## 🔐 Security Best Practices

1. **Change default passwords**:
   - PostgreSQL password
   - MinIO credentials
   - JWT secret key

2. **Use secrets management**:
   ```yaml
   secrets:
     db_password:
       file: ./secrets/db_password.txt
   ```

3. **Enable SSL for MinIO in production**:
   ```yaml
   MINIO_USE_SSL: "true"
   ```

4. **Limit network exposure**:
   ```yaml
   # Don't expose PostgreSQL port in production
   # ports:
   #   - "5432:5432"
   ```

5. **Use specific image versions**:
   ```yaml
   image: postgres:15.3-alpine  # Instead of 'latest'
   ```

## 📊 Health Checks

All services include health checks:

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres"]
  interval: 10s
  timeout: 5s
  retries: 5
```

Check health status:
```bash
docker-compose ps
```

## 🔄 Updates and Maintenance

### Update Application

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose up -d --build pet-service
```

### Update Dependencies

```bash
# Update Go modules
go get -u ./...
go mod tidy

# Rebuild image
docker-compose build --no-cache pet-service
```

### Database Migrations

```bash
# Run migrations (if implemented)
docker-compose exec pet-service ./main migrate-up

# Or connect to database directly
docker exec -it pet-service-postgres psql -U postgres -d pet_service
```

## 📱 Development Tips

### Hot Reload Setup

The `docker-compose.dev.yml` uses Air for hot reload:

```bash
# Start dev environment
docker-compose -f docker-compose.dev.yml up

# Code changes will automatically reload the application
```

### Debug Inside Container

```bash
# Get shell access
docker exec -it pet-service-api sh

# Check environment
docker exec pet-service-api env

# Test API from inside container
docker exec pet-service-api curl http://localhost:8001/health
```

### Local Development + Docker Services

Run only database and MinIO in Docker, run API locally:

```bash
# Start only dependencies
docker-compose up -d postgres minio

# Run API locally
go run main.go
```

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [MinIO Docker Image](https://hub.docker.com/r/minio/minio)
