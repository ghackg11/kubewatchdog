# KubeWatch

[![Go Report Card](https://goreportcard.com/badge/github.com/yourorg/kubewatch)](https://goreportcard.com/report/github.com/yourorg/kubewatch)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A lightweight command-line tool for real-time monitoring of Kubernetes events and database changes.

## Overview

KubeWatch provides comprehensive visibility into your Kubernetes clusters and associated databases. By tracking both Kubernetes resource events and database modifications in real-time, it helps developers and DevOps teams maintain system stability and respond quickly to potential issues.

### Key Features

- **Real-time Kubernetes event monitoring** for pods, deployments, services, and more
- **Database change tracking** with TimescaleDB integration
- **Simple CLI interface** built with Cobra for intuitive command execution
- **Low resource footprint** for efficient operation in production environments
- **Docker support** for easy deployment and integration with CI/CD pipelines

## Installation

### Using Go

```bash
go install github.com/yourorg/kubewatch@latest
```

### Using Docker

```bash
docker pull yourorg/kubewatch:latest
docker run -it --rm yourorg/kubewatch --help
```

## Usage

### Basic Commands

```bash
# Start monitoring Kubernetes events
kubewatch k8s --namespace=default

# Monitor database changes
kubewatch db --connection="postgres://user:password@localhost:5432/dbname"

# Monitor both Kubernetes and database events
kubewatch watch --namespace=default --connection="postgres://user:password@localhost:5432/dbname"

# Get help
kubewatch --help
```

### Configuration

KubeWatch can be configured using either command-line flags or a configuration file:

```bash
# Using a config file
kubewatch --config=/path/to/config.yaml

# Sample config.yaml
kubernetes:
  namespace: default
  resources:
    - pods
    - deployments
    - services
database:
  connection: "postgres://user:password@localhost:5432/dbname"
  tables:
    - users
    - products
```

## Technical Details

KubeWatch is built with Go and follows a modular architecture:

- **CLI Package**: Handles command-line interactions using the Cobra package
- **Core Package**: Contains the main logic for Kubernetes event tracking and database monitoring
- **Kubernetes Integration**: Uses client-go to establish event watchers on Kubernetes resources
- **Database Monitoring**: Integrates with TimescaleDB via GORM for efficient time-series data tracking

## Development

### Prerequisites

- Go 1.18+
- Docker (for containerized development)
- Access to a Kubernetes cluster (or minikube for local development)
- TimescaleDB instance

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourorg/kubewatch.git
cd kubewatch

# Build the binary
go build -o kubewatch main.go

# Run tests
go test ./...
```

### Docker Build

```bash
docker build -t yourorg/kubewatch:latest .
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache 2.0 License - see the LICENSE file for details.