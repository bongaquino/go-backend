# Koneksi Backend
Koneksi Backend is a monorepo containing the microservices that power the Koneksi ecosystem. Designed for scalability and efficiency, it provides core functionalities such as authentication, resource management, analytics, and integrations. Each service is modular and optimized for high-performance distributed systems.

## 🚀 **Getting Started**

This project follows a **microservices architecture**, where each service runs in a **separate container**. The services communicate over a shared **Docker network**, allowing seamless interaction.

### **🛠 Services in the Development Setup**
The following microservices and dependencies are available in this local development environment:

| Service             | Description                              | Exposed Ports                           |
|---------------------|------------------------------------------|-----------------------------------------|
| **MongoDB**         | NoSQL Database                           | `27017`                                 |
| **Mongo Express**   | Web UI for MongoDB                       | `8082`                                  |
| **Redis**           | In-memory key-value store                | `6379`                                  |
| **Redis Commander** | Web UI for Redis                         | `8083`                                  |
| **RabbitMQ**        | Message Broker                           | `5672`, `15672`                         |
| **Elasticsearch**   | Search and analytics engine              | `9200`                                  |
| **Logstash**        | Data processing pipeline                 | `12201/udp`                             |
| **Kibana**         | Visualization for Elasticsearch          | `5601`                                  |
| **Tyk API Gateway** | API Gateway for managing requests        | `8080`                                  |
| **Orchestrator**    | Manages workflow and events              | `3000`                                  |

### **📌 Prerequisites**
Ensure you have the following installed:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### **🔹 Cloning the Repository**
```sh
git clone https://github.com/koneksi-tech/koneksi-backend
cd koneksi-backend
```

### **🛠 Running the Services**
To start all services, run:
```sh
docker compose up -d
```
This will:
- Build any missing images.
- Start all containers in detached mode (`-d`).
- Automatically create and attach the **Docker network**.

### **📌 Stopping Services**
To stop the running services:
```sh
docker compose down
```

To stop a specific service:
```sh
docker compose stop <service_name>
```
Example:
```sh
docker compose stop mongo
```

### **📌 Viewing Logs**
To see logs for all services:
```sh
docker compose logs -f
```

To view logs for a specific service:
```sh
docker compose logs -f <service_name>
```
Example:
```sh
docker compose logs -f gateway
```

## ⚙ **Configuration**

### **🔹 Environment Variables**
Each microservice has its own `.env` file. The default locations are:

- **Orchestrator**: `orchestrator/.env`
- **MongoDB**: Configured inside `docker-compose.yml`
- **Tyk API Gateway**: Configured via `gateway/tyk.conf`

Make sure you update these files with the correct values before running the services.

### **🔹 Network Configuration**
All services are attached to the `network` defined in `docker-compose.yml`, allowing them to communicate using **service names**.

Example:
- **Orchestrator can connect to MongoDB** using:  
  ```
  mongodb://root:password@mongo:27017
  ```
- **Redis can be accessed at**:  
  ```
  redis://redis:6379
  ```

## 🔧 **Adding a New Microservice**
To add a new service:
1. Create a directory under `services/` (e.g., `services/new-service`).
2. Add a `Dockerfile` and configuration files.
3. Update `docker-compose.yml` with:
   ```yaml
   new-service:
     build:
       context: services/new-service
       dockerfile: Dockerfile
     container_name: new-service
     restart: unless-stopped
     ports:
       - "808X:8080"
     env_file:
       - services/new-service/.env
     networks:
       - network
   ```
4. Start the service:
   ```sh
   docker compose up -d new-service
   ```

## 🛠 **Common Issues & Fixes**

### **❌ Port Already in Use**
**Issue:**  
```
Error: Bind for 0.0.0.0:8081 failed: port is already allocated
```
**Fix:**  
Run:
```sh
docker ps
```
Find the conflicting container and stop it:
```sh
docker stop <container_id>
```

### **❌ MongoDB Connection Issues**
**Fix:**  
Ensure MongoDB is running:
```sh
docker compose up -d mongo
```

Then check logs:
```sh
docker compose logs -f mongo
```

### **❌ Service Not Found**
**Fix:**  
Run:
```sh
docker network inspect network
```
If the service is missing, restart:
```sh
docker compose down
docker compose up -d
```
