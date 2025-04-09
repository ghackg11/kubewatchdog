# Kubernetes Event Monitor

A Go-based CLI tool for monitoring and querying Kubernetes events with PostgreSQL storage and real-time event tracking.

## Features

- **Real-time Event Monitoring**: Watch and store Kubernetes events in real-time
- **Event Querying**: Query events for specific resources with customizable time windows
- **Resource Discovery**: List available Kubernetes resource types (core and API groups)
- **PostgreSQL Integration**: Store and query events from a PostgreSQL database
- **Case-insensitive Search**: Flexible resource name matching

## Prerequisites

- Go 1.23 or later
- PostgreSQL database
- Kubernetes cluster access (local or in-cluster)
- kubectl configured with cluster access

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd <repository-name>
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database:
```sql
CREATE DATABASE kubewatchdog;
\c kubewatchdog
CREATE TABLE kubernetes_events (
    id TEXT,
    event_time TIMESTAMPTZ NOT NULL,
    event_type TEXT,
    reason TEXT,
    message TEXT,
    namespace TEXT,
    resource TEXT,
    resource_name TEXT,
    PRIMARY KEY (id, event_time)
);

SELECT create_hypertable('kubernetes_events', 'event_time');
```

## Configuration

Set environment variables for different environments:

```bash
# For local development
export ENV=local

# For Kubernetes cluster
export ENV=kubernetes
```

## Usage

### List Available Resource Types

```bash
# List core resource types
go run main.go list

# List all resource types including API groups
go run main.go list --all
```

### Get Events for a Resource

```bash
# Get events for a specific resource in the last 10 minutes
go run main.go get <resource-type> <resource-name>

# Examples
go run main.go get pod my-pod
go run main.go get deployment my-deployment
```

### Start Event Monitoring

```bash
# Start watching and storing Kubernetes events
go run main.go watch
```

## Database Schema

The `kubernetes_events` table stores the following information:
- `id`: Unique identifier for the event
- `event_time`: Timestamp when the event occurred
- `event_type`: Type of the event
- `reason`: Event reason (e.g., Created, Scheduled, Started)
- `message`: Detailed event message
- `namespace`: Kubernetes namespace where the event occurred
- `resource`: Type of the Kubernetes resource
- `resource_name`: Name of the Kubernetes resource

The table uses TimescaleDB's hypertable feature for efficient time-series data storage.

## Development

### Project Structure

```
.
├── src/
│   ├── k8s.go           # Kubernetes client and event handling
│   └── phi2_client.go   # LLM client for additional features
├── cli/
│   └── cmd/
│       ├── get.go       # Get events command
│       ├── list.go      # List resources command
│       └── watch.go     # Watch events command
└── main.go              # Main application entry point
```

### Adding New Commands

1. Create a new command file in `cli/cmd/`
2. Define the command using cobra
3. Add the command to the root command in `main.go`

## License

[Your License Here]

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request