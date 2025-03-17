# Use the official Golang image as the base image
FROM golang:1.23.6 AS builder

# Set the working directory inside the container
WORKDIR /app

COPY go.mod ./

RUN go mod download

# Copy the entire source code
COPY . .

RUN go build -o gocli cli/main.go

CMD ["./gocli", "watch"]