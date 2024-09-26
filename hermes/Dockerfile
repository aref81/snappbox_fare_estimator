# Build stage
FROM golang:1.19-alpine AS builder
LABEL authors="mohammadhoseinaref"

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o hermes ./cmd

# Run stage
FROM alpine:latest

# Install CA certificates (required for Go app)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /app/hermes .

# Copy config.yaml if needed
COPY config.yaml .

# Command to run the binary
CMD ["./hermes"]
