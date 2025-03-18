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

# Start all services
echo "Starting all services..."
docker compose up -d

# Check running containers
echo "Checking running containers..."
docker ps

echo "Services started successfully."
