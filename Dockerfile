# Use the official Golang image as the base image
FROM golang:1.23.0 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire source code
COPY . .

CMD ["go", "run", "."]