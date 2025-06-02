# Use the official Golang image as the base image
FROM golang:1.24.3-alpine

# Set the working directory inside the container
WORKDIR /app

# Install Air
RUN go install github.com/air-verse/air@latest

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Expose the port the application runs on
EXPOSE 3000

# Command to run the application
CMD ["air", "-c", ".air.toml"]
