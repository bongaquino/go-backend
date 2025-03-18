#!/bin/bash

# Check if Docker is installed
if ! command -v docker &> /dev/null
then
    echo "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker compose &> /dev/null
then
    echo "Docker Compose (v2) is not installed. Please install Docker Compose Plugin."
    exit 1
fi

# Stop all services and delete named volumes
echo "Stopping all services and deleting volumes..."
docker compose down -v

echo "Specified volumes deleted successfully."

# Start services again
echo "Starting services..."
docker compose up -d --build

echo "Services refreshed successfully."
