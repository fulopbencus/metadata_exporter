# Step 1: Build the Go binary
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o metadata_exporter .

# Step 2: Create a lightweight container to run the binary
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/metadata_exporter .

# Ensure the binary has execution permissions
RUN chmod +x metadata_exporter

# Expose the port for Prometheus to scrape
EXPOSE 9091

# Run the binary
ENTRYPOINT ["./metadata_exporter"]
