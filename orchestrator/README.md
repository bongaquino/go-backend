# Argo

Powerful, opinionated API server template for Go.

### Key Features

- **Modular structure**: Clean separation of concerns with dedicated folders for controllers, models, repositories, and middleware.
- **Built-in Swagger documentation**: Easily view and test your API endpoints.
- **Dockerized services**: Ready-to-use Docker Compose setup for API, MongoDB, and Redis.
- **Custom error handling**: Centralized exception management for better debugging and error responses.
- **Live reload**: Live reloading during development with `.air.toml` configuration.

## Pre-requisites

- **Go**: Version 1.22 or higher.
- **Docker**: Version 27.4 or higher.

## Local Environment Setup

To set up the local environment using Docker Compose, follow these steps:

1. Ensure you have Docker and Docker Compose installed on your machine.

2. Create a `.env` file in the root directory of your project and add the following environment variables:

   - **PORT**: The port number the server will run on.
   - **MODE**: The mode the server will run in (debug or release).
   - **MONGO\_HOST**: The hostname for the MongoDB database.
   - **MONGO\_PORT**: The port number for the MongoDB database.
   - **MONGO\_USER**: The username for the MongoDB database.
   - **MONGO\_PASSWORD**: The password for the MongoDB database.
   - **MONGO\_DATABASE**: The name of the MongoDB database.
   - **REDIS\_HOST**: The hostname for the Redis database.
   - **REDIS\_PORT**: The port number for the Redis database.
   - **REDIS\_PASSWORD**: The password for the Redis database.

3. Run the following command to start the services:

   ```sh
   docker compose up -d
   ```

   This command will start the following services:

   - **API Server**
   - **MongoDB**
   - **Mongo Express**
   - **Redis**
   - **Redis Commander**

4. Access the web app interface in your browser:

   - **API Documentation**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
   - **Mongo Express**: [http://localhost:8081](http://localhost:8081)
   - **Redis Commander**: [http://localhost:8082](http://localhost:8082)

5. To stop the services, run:

   ```sh
   docker compose down
   ```

## Documentation

The API documentation is generated using Swagger. Run the following command in the root directory of your project to regenerate the Swagger documentation:

```sh
swag init -g main.go
```
