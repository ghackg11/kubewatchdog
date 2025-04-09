# Use the official Golang image as the base image
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and download dependencies first (cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go binary
RUN GOOS=linux GOARCH=amd64 go build -o /app/kubewatch cli/main.go

# Use a minimal image for the final stage
FROM debian:latest

WORKDIR /root/

RUN apt update

RUN apt -y install curl

# Copy only the compiled binary
COPY --from=builder /app/kubewatch /usr/local/bin/kubewatch

# Ensure the binary is executable
RUN chmod +x /usr/local/bin/kubewatch

# Set the entrypoint and command properly
ENTRYPOINT ["/usr/local/bin/kubewatch"]
CMD ["watch"]
