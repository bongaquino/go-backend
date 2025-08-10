<<<<<<< HEAD
# Bong Aquino Backend
Bong Aquino Backend is a monorepo containing the services that power the Bong Aquino ecosystem.
=======
# bongaquino Backend
bongaquino Backend is a monorepo containing the services that power the bongaquino ecosystem.
>>>>>>> 1348c53f1e2e1dd6f94dc4c583ed02b1e28350ee

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
| **Elasticsearch**   | Search and analytics engine              | `9200`                                  |
| **Logstash**        | Data processing pipeline                 | `12201/udp`                             |
| **Kibana**          | Visualization for Elasticsearch          | `5601`                                  |
| **Tyk API Gateway** | API Gateway for managing requests        | `8080`                                  |
<<<<<<< HEAD
| **Bong Aquino Server**  | Core backend server                      | `3000`                                  |
=======
| **bongaquino Server**  | Core backend server                      | `3000`                                  |
>>>>>>> 1348c53f1e2e1dd6f94dc4c583ed02b1e28350ee

### **📌 Prerequisites**
Ensure you have the following installed:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### **🔹 Cloning the Repository**
```sh
<<<<<<< HEAD
git clone https://github.com/bongaquino/bongaquino-backend
cd bongaquino-backend
=======
git clone https://github.com/bongaquino-tech/bongaquino-backend
cd bongaquino-backend
>>>>>>> 1348c53f1e2e1dd6f94dc4c583ed02b1e28350ee
```

### **🔹 Creating the shared Docker network**
```sh
<<<<<<< HEAD
docker network create bongaquino-network
=======
docker network create bongaquino-network
>>>>>>> 1348c53f1e2e1dd6f94dc4c583ed02b1e28350ee
```

### **🔹 Starting the Services**
To start all services, run:
```sh
./scripts/start.sh
```
This will:
- Start all containers in detached mode (`-d`).
- Automatically create and attach the **Docker network**.

### **🔹 Stopping the Services**
To stop the running services:
```sh
./scripts/stop.sh
```

To stop a specific service:
```sh
docker compose stop <service_name>
```
Example:
```sh
docker compose stop mongo
```

### **🔹 Restarting Services**
To restart all services:
```sh
./scripts/restart.sh
```

### **🔹 Rebuilding the Setup**
If you need to **rebuild everything and reset volumes**, run:
```sh
./scripts/rebuild.sh
```
This will:
- Stop all services.
- Remove named volumes.
- Restart everything with a clean state.

### **🔹 Viewing Logs**
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
docker compose logs -f server
```

## ⚙ **Configuration**

### **🔹 Environment Variables**
The main server has its own `.env` file. The location is:

<<<<<<< HEAD
- **Bong Aquino Server**: `server/.env`
=======
- **bongaquino Server**: `server/.env`
>>>>>>> 1348c53f1e2e1dd6f94dc4c583ed02b1e28350ee

Make sure you update the file with the correct values before running the services.

### **🔹 Network Configuration**
All services are attached to the `network` defined in `docker-compose.yml`, allowing them to communicate using **service names**.