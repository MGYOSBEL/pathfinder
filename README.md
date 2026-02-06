# pathfinder

Data Platform project for IoT data ingestion, processing, and storage with comprehensive observability.

## Architecture Overview

Pathfinder processes MQTT messages through Benthos/Redpanda Connect processors and forwards them to timeseries databases:

```
MQTT Broker (VerneMQ/HiveMQ)
    ↓
Benthos Injector (MQTT → RabbitMQ)
    ↓
RabbitMQ Message Broker
    ↓
Benthos Writer (RabbitMQ → Database)
    ↓
TimescaleDB / InfluxDB
```

### Key Components

- **MQTT Brokers**: VerneMQ or HiveMQ for IoT device connectivity
- **Benthos Processors**: Data ingestion and transformation pipelines
- **Message Broker**: RabbitMQ for reliable message queuing
- **Timeseries Databases**: TimescaleDB (PostgreSQL) or InfluxDB for data storage
- **Object Storage**: MinIO for file storage and backups
- **Observability Stack**: Prometheus, Grafana, Alertmanager, Loki, Tempo

## Get Started

### Prerequisites

- Docker and Docker Compose
- Make (optional, for convenience commands)
- 8GB+ RAM recommended
- 20GB+ disk space

### Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd pathfinder
   ```

2. **Create environment file**
   
   Create a `.env` file in the `deploy/docker` directory with the following configuration:
   
   **Minimum Required Variables** (for basic deployment):
   ```bash
   # PostgreSQL/TimescaleDB (Required)
   POSTGRES_ADMIN_PASSWORD=changeme123
   TIMESCALE_ADMIN_USER=postgres
   TIMESCALE_ADMIN_PASSWORD=changeme123
   TIMESCALE_DBNAME=tsdb
   TIMESCALE_URL=timescaledb:5432

   # RabbitMQ (Required)
   RABBITMQ_ADMIN_USER=admin
   RABBITMQ_ADMIN_PASSWORD=changeme123
   ```

   **Optional: Alert Notifications**
   
   To enable Slack/Email alerts, add:
   ```bash
   # Email Notifications (SMTP)
   SMTP_HOST=smtp.gmail.com:587
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-app-specific-password
   ALERT_EMAIL_FROM=alerts@pathfinder.local
   ALERT_EMAIL_CRITICAL=ops-oncall@company.com
   ALERT_EMAIL_HIGH=ops-team@company.com

   # Slack Notifications
   SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
   SLACK_CRITICAL_CHANNEL=#alerts-critical
   SLACK_HIGH_CHANNEL=#alerts-high
   SLACK_MEDIUM_CHANNEL=#alerts-medium
   ```
   
   **Optional: InfluxDB** (if using InfluxDB instead of TimescaleDB):
   ```bash
   SESSION_SECRET_KEY=<Generate with: openssl rand -hex 32>
   INFLUXDB_ADMIN_TOKEN=<Generate after first deployment>
   INFLUXDB_WRITE_TOKEN=<Generate after first deployment>
   ```

   **Environment Variables Reference**:
   
   | Variable | Required | Default | Description |
   |----------|----------|---------|-------------|
   | `POSTGRES_ADMIN_PASSWORD` | Yes | - | PostgreSQL admin password |
   | `TIMESCALE_ADMIN_USER` | Yes | postgres | TimescaleDB user |
   | `TIMESCALE_ADMIN_PASSWORD` | Yes | - | TimescaleDB password |
   | `TIMESCALE_DBNAME` | Yes | tsdb | TimescaleDB database name |
   | `TIMESCALE_URL` | Yes | timescaledb:5432 | Database host:port (use service name for local) |
   | `RABBITMQ_ADMIN_USER` | Yes | - | RabbitMQ admin username |
   | `RABBITMQ_ADMIN_PASSWORD` | Yes | - | RabbitMQ admin password |
   | `SMTP_HOST` | No | - | SMTP server for email alerts (e.g., smtp.gmail.com:587) |
   | `SMTP_USERNAME` | No | - | SMTP authentication username |
   | `SMTP_PASSWORD` | No | - | SMTP authentication password (use app-specific password for Gmail) |
   | `ALERT_EMAIL_FROM` | No | alerts@pathfinder.local | Alert sender email address |
   | `ALERT_EMAIL_CRITICAL` | No | - | Email for P1 critical alerts |
   | `ALERT_EMAIL_HIGH` | No | - | Email for P2 high severity alerts |
   | `SLACK_WEBHOOK_URL` | No | - | Slack webhook URL for notifications |
   | `SLACK_CRITICAL_CHANNEL` | No | #alerts-critical | Slack channel for critical alerts |
   | `SLACK_HIGH_CHANNEL` | No | #alerts-high | Slack channel for high severity alerts |
   | `SLACK_MEDIUM_CHANNEL` | No | #alerts-medium | Slack channel for medium severity alerts |
   | `SESSION_SECRET_KEY` | InfluxDB only | - | InfluxDB session secret (generate with `openssl rand -hex 32`) |
   | `INFLUXDB_ADMIN_TOKEN` | InfluxDB only | - | InfluxDB admin token |
   | `INFLUXDB_WRITE_TOKEN` | InfluxDB only | - | InfluxDB write token |

   **SMTP Provider Examples**:
   - **Gmail**: `smtp.gmail.com:587` (requires app-specific password with 2FA enabled)
   - **Office 365**: `smtp.office365.com:587`
   - **SendGrid**: `smtp.sendgrid.net:587` (use API key as password)
   - **AWS SES**: `email-smtp.us-east-1.amazonaws.com:587`

   **Slack Webhook Setup**:
   1. Go to https://api.slack.com/apps
   2. Create New App → From scratch
   3. Enable "Incoming Webhooks"
   4. Click "Add New Webhook to Workspace"
   5. Select channel and authorize
   6. Copy webhook URL to `SLACK_WEBHOOK_URL`

   **Security Notes**:
   - Generate strong passwords in production (use `openssl rand -base64 32`)
   - Never commit `.env` file to version control
   - Use secrets manager for production (Vault, AWS Secrets Manager, etc.)
   - Rotate credentials regularly
   - For cloud databases, add `?sslmode=require` to `TIMESCALE_URL`


3. **Deploy the platform**
   ```bash
   # Using Make (recommended)
   make deploy MQTT_BROKER=vernemq BROKER=rabbitmq TIMESERIES_DB=timescaledb

   # Or using Docker Compose directly
   docker compose -f deploy/docker/docker-compose.yaml \
     -f deploy/docker/mqtt/vernemq.yaml \
     -f deploy/docker/injectors/mqtt-rabbitmq/services.yaml \
     -f deploy/docker/broker/rabbitmq.yaml \
     -f deploy/docker/writers/rabbitmq/timescaledb/services.yaml \
     -f deploy/docker/timeseriesdb/timescaledb.yaml \
     up -d
   ```

4. **Verify deployment**
   ```bash
   # Check all containers are running
   docker ps

   # Access Grafana dashboards
   open http://localhost:3000  # Default: admin / admin
   ```

### Configuration Options

#### MQTT Broker Selection
- **VerneMQ**: `MQTT_BROKER=vernemq` (default, lightweight)
- **HiveMQ**: `MQTT_BROKER=hivemq` (enterprise-grade)

#### Message Broker
- **RabbitMQ**: `BROKER=rabbitmq` (currently supported)

#### Timeseries Database
- **TimescaleDB**: `TIMESERIES_DB=timescaledb` (PostgreSQL-based, SQL queries)
- **InfluxDB**: `TIMESERIES_DB=influxdb` (purpose-built for timeseries)

### Make Commands

```bash
make deploy      # Start the full platform
make config      # Parse and validate Docker Compose configuration
make stop        # Stop containers (keep them)
make down        # Delete containers (keep data)
make destroy     # Delete containers and data
```

## Observability

Pathfinder includes production-ready observability with 16 Grafana dashboards, 26 alerting rules, and comprehensive metrics collection.

### Quick Access

| Service | URL | Credentials |
|---------|-----|-------------|
| **Grafana** | http://localhost:3000 | admin / admin |
| **Prometheus** | http://localhost:9090 | None |
| **Alertmanager** | http://localhost:9093 | None |
| **HAProxy Stats** | http://localhost:8404 | None |
| **HAProxy Metrics** | http://localhost:8405/metrics | None |

### Dashboard Organization

Dashboards are organized by architectural layer:

#### 00-OVERVIEW
- **Platform Overview**: Single-pane-of-glass view with system health, pipeline metrics, resource utilization
- **Observability Stack Health**: Monitoring infrastructure self-monitoring

#### 01-INFRASTRUCTURE
- **Infrastructure Overview**: Platform health summary, resource metrics, HAProxy stats
- **Node Dashboard**: Detailed per-node CPU, memory, disk, network metrics

#### 02-MESSAGING
- **VerneMQ/HiveMQ MQTT Broker**: MQTT broker metrics and performance
- **RabbitMQ Broker**: Message broker queues, connections, throughput
- **Message Flow Overview**: End-to-end pipeline visualization

#### 03-DATA-STORAGE
- **TimescaleDB**: Database performance, connections, query stats
- **MinIO**: S3 storage usage, operations, cluster health

#### 04-APPLICATION
- **Benthos Injector**: MQTT-to-RabbitMQ pipeline metrics
- **Benthos Writer**: RabbitMQ-to-Database pipeline metrics

### Alerting

26 production-ready alerts across 6 groups:
- **Service Availability** (6 alerts): Critical service health
- **Performance Degradation** (4 alerts): Latency and throughput SLOs
- **Reliability SLOs** (4 alerts): Data delivery and error rates
- **Resource Exhaustion** (6 alerts): Infrastructure capacity
- **Observability Stack** (4 alerts): Monitoring system health
- **SLO Error Budget** (2 alerts): Error budget tracking

**Alert Severity Levels**:
- **P1 Critical**: Immediate response (service down, data loss)
- **P2 High**: 1-hour response (performance degradation)
- **P3 Medium**: 24-hour response (capacity warnings)

### Metrics Endpoints

All metrics accessible via HAProxy on port 8405:
- `/metrics/vernemq` - VerneMQ MQTT broker
- `/metrics/hivemq` - HiveMQ MQTT broker
- `/metrics/rabbitmq` - RabbitMQ message broker
- `/metrics/postgres` - PostgreSQL/TimescaleDB
- `/metrics/minio` - MinIO S3 storage
- `/metrics/benthos-injector` - Data injector pipeline
- `/metrics/benthos-writer` - Data writer pipeline

### Documentation

Comprehensive observability documentation:
- **[Observability Runbook](deploy/docker/observability/OBSERVABILITY-RUNBOOK.md)**: Complete operational guide with troubleshooting playbooks
- **[SLI/SLO Definitions](deploy/docker/observability/SLI-SLO-DEFINITIONS.md)**: Service level indicators and objectives with PromQL queries
- **[Alertmanager Setup](deploy/docker/observability/ALERTMANAGER-SETUP.md)**: Alert notification configuration guide

### Key SLOs

- **Availability**: 99.5% uptime (21.6 min downtime/month)
- **End-to-end latency**: P99 < 5 seconds
- **Message delivery**: 99.9% success rate
- **Database writes**: P99 < 2 seconds

## Development

### Building Go Services

```bash
# Build metrics injector
go build ./cmd/metrics-injector/

# Run locally
go run ./cmd/metrics-injector/
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/mqtt/
```

### Code Style

- Follow Go conventions (see [AGENTS.md](AGENTS.md))
- Use `go fmt` and `goimports` for formatting
- Run `go vet` for static analysis

## Troubleshooting

### Common Issues

**No data in dashboards**:
```bash
# Check Prometheus is scraping
curl http://localhost:9090/targets

# Check service metrics endpoints
curl http://localhost:8405/metrics/benthos-injector
```

**Container won't start**:
```bash
# Check logs
docker logs <container-name> --tail 100

# Check resources
docker stats
```

**High queue depth**:
```bash
# Check Message Flow dashboard in Grafana
# Check writer logs
docker logs rabbitmq-timescale-writer --tail 100
```

See [Observability Runbook](deploy/docker/observability/OBSERVABILITY-RUNBOOK.md) for detailed troubleshooting procedures.

## Project Structure

```
pathfinder/
├── cmd/                          # Application entry points
├── internal/                     # Private application code
├── pkg/                          # Public libraries
├── benthos/                      # Benthos processor configs
├── deploy/docker/                # Docker Compose configurations
│   ├── observability/           # Prometheus, Grafana, Alertmanager
│   │   ├── grafana/provisioning/dashboards/  # 16 Grafana dashboards
│   │   ├── prometheus.yml       # Metrics scraping config
│   │   ├── alerting-rules.yml   # 26 alerting rules
│   │   └── *.md                 # Observability documentation
│   ├── mqtt/                    # MQTT broker configs (VerneMQ, HiveMQ)
│   ├── broker/                  # Message broker configs (RabbitMQ)
│   ├── timeseriesdb/            # Database configs (TimescaleDB, InfluxDB)
│   ├── injectors/               # Benthos injector services
│   └── writers/                 # Benthos writer services
├── Makefile                     # Build and deployment commands
└── README.md                    # This file
```

## Contributing

1. Follow code style guidelines in [AGENTS.md](AGENTS.md)
2. Write tests for new features
3. Update documentation as needed
4. Use conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`

## License

See [LICENSE](LICENSE) file for details.

## Additional Resources

- **[AGENTS.md](AGENTS.md)**: Development guidelines and code style
- **[TODO.md](TODO.md)**: Project roadmap and tasks
- **[Observability Runbook](deploy/docker/observability/OBSERVABILITY-RUNBOOK.md)**: Operations guide
