# Build stage
FROM golang:1.23-alpine AS builder

# Install necessary dependencies
RUN apk add --no-cache git make build-base

# Define the working directory
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/main .
COPY --from=builder /app/config/config-prod.yml ./config/config.yml

# Expose the port
EXPOSE 8080

# Configure gin to run in production mode
ENV GIN_MODE=release

# Run the binary
CMD ["/app/main"]
