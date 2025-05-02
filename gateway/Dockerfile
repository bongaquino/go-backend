# gateway/Dockerfile

# Use the official Tyk Gateway image as the base image
FROM tykio/tyk-gateway:latest

# Copy configuration files
# For deplyoment to production environment
COPY ./apps/koneksi-server.json /opt/tyk-gateway/apps/koneksi-server.json
COPY ./apps/koneksi-ipfs-gateway.json /opt/tyk-gateway/apps/koneksi-ipfs-gateway.json
COPY ./apps/koneksi-ipfs.json /opt/tyk-gateway/apps/koneksi-ipfs.json
COPY ./apps/koneksi-ipfs-upload.json /opt/tyk-gateway/apps/koneksi-ipfs-upload.json
COPY ./templates /opt/tyk-gateway/templates
COPY ./tyk.conf /opt/tyk-gateway/tyk.conf
COPY ./policy/policies.json /opt/tyk-gateway/policies/policies.json

# Expose the port that Tyk Gateway is using
EXPOSE 8080