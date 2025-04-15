# gateway/Dockerfile

# Use the official Tyk Gateway image as the base image
FROM tykio/tyk-gateway:latest

# Copy configuration files
COPY ./apps /opt/tyk-gateway/apps
COPY ./templates /opt/tyk-gateway/templates
COPY ./tyk.conf /opt/tyk-gateway/tyk.conf

# Expose the port that Tyk Gateway is using
EXPOSE 8080