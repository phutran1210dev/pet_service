#!/bin/bash

# Pet Service Docker Stop Script

set -e

echo "🛑 Pet Service API - Docker Stop"
echo "================================"
echo ""

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed."
    exit 1
fi

# Ask user what to stop
echo "What would you like to stop?"
echo "1) Production services"
echo "2) Development services"
echo "3) Both (all services)"
echo "4) Stop and remove volumes (⚠️  deletes all data)"
echo ""
read -p "Enter choice [1-4]: " choice

case $choice in
    1)
        echo ""
        echo "Stopping production services..."
        docker-compose down
        echo "✅ Production services stopped"
        ;;
    2)
        echo ""
        echo "Stopping development services..."
        docker-compose -f docker-compose.dev.yml down
        echo "✅ Development services stopped"
        ;;
    3)
        echo ""
        echo "Stopping all services..."
        docker-compose down 2>/dev/null || true
        docker-compose -f docker-compose.dev.yml down 2>/dev/null || true
        echo "✅ All services stopped"
        ;;
    4)
        echo ""
        echo "⚠️  WARNING: This will delete all data (database, uploaded files)!"
        read -p "Are you sure? Type 'yes' to confirm: " confirm
        if [ "$confirm" == "yes" ]; then
            echo ""
            echo "Stopping and removing all services and volumes..."
            docker-compose down -v 2>/dev/null || true
            docker-compose -f docker-compose.dev.yml down -v 2>/dev/null || true
            echo "✅ All services and volumes removed"
        else
            echo "❌ Cancelled"
            exit 1
        fi
        ;;
    *)
        echo "❌ Invalid choice. Exiting."
        exit 1
        ;;
esac

echo ""
echo "🎉 Done!"
