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

# Function to start a service and wait until it's healthy
start_service() {
    local service_name=$1
    echo "Starting $service_name..."
    docker compose up -d --build $service_name
    echo "Waiting for $service_name to be healthy..."
    while ! docker inspect --format='{{.State.Health.Status}}' $(docker compose ps -q $service_name) 2>/dev/null | grep -q "healthy"; do
        echo "Waiting for $service_name..."
        sleep 3
    done
    echo "$service_name is now running."
}

# Start services one by one
start_service mongo
start_service redis
start_service rabbitmq
start_service elasticsearch
start_service logstash
start_service kibana
start_service api-gateway
start_service account-service
start_service backup-service
start_service dashboard-service

# Check running containers
echo "Checking running containers..."
docker ps

echo "All services started successfully."
