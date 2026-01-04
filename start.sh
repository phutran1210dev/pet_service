#!/bin/bash

# Pet Service Docker Quick Start Script

set -e

echo "🐾 Pet Service API - Docker Quick Start"
echo "========================================"
echo ""

# Check if docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    echo "   Visit: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "✅ Docker and Docker Compose are installed"
echo ""

# Ask user for mode
echo "Select mode:"
echo "1) Production (optimized build)"
echo "2) Development (with hot reload)"
echo ""
read -p "Enter choice [1-2]: " choice

case $choice in
    1)
        echo ""
        echo "🚀 Starting Production Mode..."
        echo ""
        docker-compose up -d --build
        MODE="production"
        COMPOSE_FILE="docker-compose.yml"
        ;;
    2)
        echo ""
        echo "🔧 Starting Development Mode..."
        echo ""
        docker-compose -f docker-compose.dev.yml up -d --build
        MODE="development"
        COMPOSE_FILE="docker-compose.dev.yml"
        ;;
    *)
        echo "❌ Invalid choice. Exiting."
        exit 1
        ;;
esac

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check if services are running
if [ "$MODE" == "production" ]; then
    if docker-compose ps | grep -q "pet-service-api.*Up"; then
        echo ""
        echo "✅ All services started successfully!"
    else
        echo ""
        echo "⚠️  Services may still be starting. Check status with: docker-compose ps"
    fi
else
    if docker-compose -f docker-compose.dev.yml ps | grep -q "pet-service-api-dev.*Up"; then
        echo ""
        echo "✅ All services started successfully!"
    else
        echo ""
        echo "⚠️  Services may still be starting. Check status with: docker-compose -f docker-compose.dev.yml ps"
    fi
fi

echo ""
echo "📊 Service URLs:"
echo "   🌐 API:              http://localhost:8001"
echo "   ❤️  Health Check:     http://localhost:8001/health"
echo "   🗄️  PostgreSQL:       localhost:5432"
echo "   📦 MinIO API:        http://localhost:9000"
echo "   🎛️  MinIO Console:    http://localhost:9001"
echo "      (Login: minioadmin / minioadmin)"
echo ""
echo "📝 Useful Commands:"
if [ "$MODE" == "production" ]; then
    echo "   View logs:       docker-compose logs -f"
    echo "   Stop services:   docker-compose down"
    echo "   Restart:         docker-compose restart"
else
    echo "   View logs:       docker-compose -f docker-compose.dev.yml logs -f"
    echo "   Stop services:   docker-compose -f docker-compose.dev.yml down"
    echo "   Restart:         docker-compose -f docker-compose.dev.yml restart"
fi
echo ""
echo "🎉 Setup complete! Your Pet Service API is running!"
echo ""

# Ask if user wants to see logs
read -p "Do you want to view logs now? [y/N]: " view_logs
if [[ $view_logs == "y" || $view_logs == "Y" ]]; then
    echo ""
    if [ "$MODE" == "production" ]; then
        docker-compose logs -f
    else
        docker-compose -f docker-compose.dev.yml logs -f
    fi
fi
