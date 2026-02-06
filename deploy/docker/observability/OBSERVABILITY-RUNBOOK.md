# Pathfinder Observability Runbook

## Table of Contents
1. [Overview](#overview)
2. [Quick Access](#quick-access)
3. [Metrics Endpoints](#metrics-endpoints)
4. [Dashboard Guide](#dashboard-guide)
5. [Alert Response Procedures](#alert-response-procedures)
6. [Troubleshooting Playbooks](#troubleshooting-playbooks)
7. [Common Operations](#common-operations)

---

## Overview

Pathfinder's observability stack provides complete visibility into the IoT data pipeline from MQTT ingestion through timeseries storage. The stack includes:

- **Prometheus**: Metrics collection and alerting engine
- **Grafana**: Visualization and dashboards (16 dashboards across 5 layers)
- **Alertmanager**: Alert routing and notification management
- **Loki**: Log aggregation (optional)
- **Tempo**: Distributed tracing (optional)

**Architecture**: MQTT → Benthos Injector → RabbitMQ → Benthos Writer → TimescaleDB/InfluxDB

---

## Quick Access

| Service | URL | Credentials |
|---------|-----|-------------|
| Grafana | http://localhost:3000 | admin / admin (change on first login) |
| Prometheus | http://localhost:9090 | None (public access) |
| Alertmanager | http://localhost:9093 | None (public access) |
| HAProxy Stats | http://localhost:8404 | None (public access) |
| HAProxy Metrics | http://localhost:8405/metrics | Centralized metrics gateway |

### Dashboard Quick Links
- **Platform Overview**: http://localhost:3000/d/platform-overview
- **Infrastructure Overview**: http://localhost:3000/d/infrastructure-overview
- **Message Flow**: http://localhost:3000/d/message-flow-overview
- **Observability Stack**: http://localhost:3000/d/observability-stack-health

---

## Metrics Endpoints

All metrics are accessible internally via service ports and externally via HAProxy on port 8405.

### Infrastructure Metrics

| Service | Internal Endpoint | HAProxy Route | Metrics Format |
|---------|------------------|---------------|----------------|
| Node Exporter | node-exporter:9100/metrics | N/A | Prometheus |
| cAdvisor | cadvisor:8080/metrics | N/A | Prometheus |
| HAProxy | haproxy:8404/metrics | /metrics | Prometheus |

**Key Metrics:**
```promql
# CPU Usage
100 - (avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

# Memory Usage
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100

# Disk Usage
100 - ((node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100)
```

### MQTT Broker Metrics

| Broker | Internal Endpoint | HAProxy Route | Metrics Format |
|--------|------------------|---------------|----------------|
| VerneMQ | mqtt-broker:8888/metrics | /metrics/vernemq | Prometheus |
| HiveMQ | mqtt-broker:9399/metrics | /metrics/hivemq | Prometheus |

**Key Metrics:**
```promql
# VerneMQ
vernemq_mqtt_publish_received # Messages received
vernemq_mqtt_publish_sent      # Messages sent
vernemq_queue_message_in       # Queue in-flight

# HiveMQ
com_hivemq_messages_incoming_publish_count_total  # Messages received
com_hivemq_messages_outgoing_publish_count_total  # Messages sent
com_hivemq_clients_connected                      # Connected clients
```

### Message Broker Metrics

| Service | Internal Endpoint | HAProxy Route | Metrics Format |
|---------|------------------|---------------|----------------|
| RabbitMQ | rabbitmq:15692/metrics | /metrics/rabbitmq | Prometheus |

**Key Metrics:**
```promql
rabbitmq_queue_messages              # Total messages in queues
rabbitmq_queue_messages_ready        # Messages ready for delivery
rabbitmq_connections                 # Active connections
rabbitmq_consumers                   # Active consumers
rabbitmq_queue_messages_published_total  # Published messages rate
```

### Application Metrics (Benthos)

| Service | Internal Endpoint | HAProxy Route | Metrics Format |
|---------|------------------|---------------|----------------|
| Data Injector | data-injector:4195/benthos/metrics | /metrics/benthos-injector | Prometheus |
| Timeseries Writer | timeseries-writer:4195/benthos/metrics | /metrics/benthos-writer | Prometheus |

**Key Metrics:**
```promql
# Input/Output
input_received{instance="data-injector:4195"}       # MQTT messages received
output_sent{instance="data-injector:4195"}          # RabbitMQ messages sent
input_received{instance="timeseries-writer:4195"}   # RabbitMQ messages received
output_sent{instance="timeseries-writer:4195"}      # Database writes sent

# Latency
histogram_quantile(0.99, rate(input_latency_ns_bucket[5m])) / 1e9  # P99 input latency
histogram_quantile(0.99, rate(output_latency_ns_bucket[5m])) / 1e9  # P99 output latency

# Errors
input_error                  # Input errors
output_error                 # Output errors
input_connection_failed      # Connection failures
```

### Database Metrics

| Service | Internal Endpoint | HAProxy Route | Metrics Format |
|---------|------------------|---------------|----------------|
| PostgreSQL Exporter | postgres-exporter:9187/metrics | /metrics/postgres | Prometheus |
| MinIO | minio:9000/minio/v2/metrics/cluster | /metrics/minio | Prometheus |

**Key Metrics:**
```promql
# TimescaleDB
pg_database_size_bytes{datname="tsdb"}     # Database size
pg_stat_database_tup_inserted              # Rows inserted
pg_stat_database_xact_commit               # Transactions committed
pg_stat_database_blks_hit / pg_stat_database_blks_read  # Cache hit ratio

# MinIO
minio_s3_requests_total                    # S3 API requests
minio_bucket_usage_total_bytes             # Storage usage
minio_cluster_nodes_online                 # Online nodes
```

### Observability Stack Metrics

| Service | Internal Endpoint | HAProxy Route | Metrics Format |
|---------|------------------|---------------|----------------|
| Prometheus | prometheus:9090/metrics | N/A | Prometheus |
| Grafana | grafana:3000/metrics | N/A | Prometheus |
| Loki | loki:3100/metrics | N/A | Prometheus |
| Tempo | tempo:3200/metrics | N/A | Prometheus |

---

## Dashboard Guide

### 00-OVERVIEW Layer

#### Platform Overview
**Purpose**: Single-pane-of-glass view of entire platform health  
**Use Case**: Daily operations, incident detection, executive reporting  
**Sections**:
- System Health Status: 8 service health indicators (MQTT, Injector, RabbitMQ, Writer, DB, MinIO, Prometheus, Grafana)
- Data Pipeline Metrics: Message flow, queue depth, write throughput
- Resource Utilization: CPU, memory, disk usage with gauges
- Error Rates: Application and infrastructure errors

**When to Use**:
- ✅ Morning health check
- ✅ Incident triage (what's broken?)
- ✅ Capacity planning overview
- ✅ Executive status reporting

#### Observability Stack Health
**Purpose**: Monitor the monitoring infrastructure itself  
**Use Case**: Ensure Prometheus/Grafana/Loki are healthy  
**Sections**:
- Prometheus Health: Targets, TSDB size, ingestion rate
- Service Status: Grafana, Loki, Tempo availability
- Storage & Resources: Data retention, disk usage

**When to Use**:
- ✅ Troubleshooting missing metrics
- ✅ Prometheus storage capacity planning
- ✅ Verifying observability stack after deployment

### 01-INFRASTRUCTURE Layer

#### Infrastructure Overview
**Purpose**: Detailed infrastructure monitoring  
**Sections**:
- Platform Health Summary: All services up/down status
- Platform Overall Metrics: CPU, memory, disk across all nodes
- HAProxy Statistics: Frontend/backend performance

**When to Use**:
- ✅ Investigating resource exhaustion alerts
- ✅ Analyzing HAProxy load balancing
- ✅ Infrastructure capacity planning

#### Node Dashboard
**Purpose**: Deep dive into individual node metrics  
**Variables**: $environment, $node  
**Sections**:
- CPU, memory, network, disk I/O per node

**When to Use**:
- ✅ Investigating high CPU/memory alerts
- ✅ Analyzing disk I/O bottlenecks
- ✅ Network troubleshooting

### 02-MESSAGING Layer

#### VerneMQ / HiveMQ Dashboards
**Purpose**: MQTT broker monitoring  
**Variables**: $environment  
**Sections**:
- Connections, messages in/out, subscriptions
- Queue depth, session metrics
- Network bandwidth

**When to Use**:
- ✅ MQTT broker performance issues
- ✅ Connection drop troubleshooting
- ✅ Message throughput analysis

#### RabbitMQ Dashboard
**Purpose**: Message broker monitoring  
**Variables**: $environment  
**Sections**:
- Overview: Connections, channels, consumers
- Queues: Depth, publish/deliver rates
- Node metrics: Memory, file descriptors

**When to Use**:
- ✅ Queue depth alerts investigation
- ✅ Consumer lag analysis
- ✅ RabbitMQ capacity planning

#### Message Flow Overview
**Purpose**: End-to-end pipeline visualization  
**Variables**: $environment  
**Sections**:
- MQTT → Injector → RabbitMQ → Writer → Database flow
- Message rates at each stage
- Queue depths and backpressure detection

**When to Use**:
- ✅ Pipeline bottleneck identification
- ✅ Data loss investigation
- ✅ Performance optimization

### 03-DATA-STORAGE Layer

#### TimescaleDB Dashboard
**Purpose**: Database performance monitoring  
**Variables**: $environment  
**Sections**:
- Database size, connections, TPS
- Cache hit ratio, query performance
- Table sizes, indexes

**When to Use**:
- ✅ Database performance degradation
- ✅ Connection pool exhaustion
- ✅ Storage capacity planning

#### MinIO Dashboard
**Purpose**: S3 object storage monitoring  
**Variables**: $environment  
**Sections**:
- Storage usage, bucket metrics
- S3 operations (GET, PUT, DELETE)
- Network I/O, cluster health

**When to Use**:
- ✅ S3 performance issues
- ✅ Storage capacity planning
- ✅ Bucket usage analysis

### 04-APPLICATION Layer

#### Benthos Injector Dashboard
**Purpose**: MQTT-to-RabbitMQ pipeline monitoring  
**Variables**: $environment, $injector (auto-discovered)  
**Sections**:
- Input: MQTT messages received, connections
- Output: RabbitMQ messages sent
- Latency: P50, P95, P99 processing time
- Errors: Connection failures, processing errors

**When to Use**:
- ✅ Data ingestion issues
- ✅ MQTT connection problems
- ✅ Pipeline latency analysis

#### Benthos Writer Dashboard
**Purpose**: RabbitMQ-to-Database pipeline monitoring  
**Variables**: $environment, $writer (auto-discovered with regex)  
**Sections**:
- Input: RabbitMQ messages consumed
- Output: Database writes
- Latency: Write performance
- Errors: Database connection/write failures

**When to Use**:
- ✅ Database write issues
- ✅ Consumer lag investigation
- ✅ Write throughput optimization

---

## Alert Response Procedures

### Alert Severity Levels

| Severity | Response Time | Escalation | Examples |
|----------|--------------|------------|----------|
| **P1 Critical** | Immediate | After 5 min | Service down, data loss |
| **P2 High** | Within 1 hour | After 1 hour | Performance degradation |
| **P3 Medium** | Within 24 hours | After 2 days | Capacity warnings |
| **P4 Low** | Best effort | After 1 week | Observability issues |

### P1 Critical Alerts

#### MQTTBrokerDown
**Impact**: All data ingestion is blocked  
**Response**:
1. Check broker container status: `docker ps -a | grep mqtt-broker`
2. Check logs: `docker logs mqtt-broker --tail 100`
3. Restart if needed: `docker restart mqtt-broker`
4. Verify connectivity: `curl http://localhost:8888/metrics` (VerneMQ) or `curl http://localhost:9399/metrics` (HiveMQ)
5. Check Grafana MQTT dashboard for recovery

**Root Causes**:
- Container crash (check logs for OOM, panic)
- Network issues (firewall, DNS)
- Resource exhaustion (CPU, memory)
- Configuration error after deployment

#### DataInjectorDown / DataWriterDown
**Impact**: Pipeline broken, data not flowing  
**Response**:
1. Check container: `docker ps -a | grep -E "injector|writer"`
2. Check logs: `docker logs <container> --tail 100`
3. Verify dependencies are up (MQTT broker, RabbitMQ, TimescaleDB)
4. Check Benthos metrics: `curl http://localhost:8405/metrics/benthos-injector`
5. Restart if needed: `docker restart <container>`

**Root Causes**:
- Dependency unavailable (broker, queue, database)
- Configuration error in benthos YAML
- Network connectivity issues
- Resource exhaustion

#### RabbitMQDown
**Impact**: Message brokering unavailable, pipeline broken  
**Response**:
1. Check container: `docker ps -a | grep rabbitmq`
2. Check logs: `docker logs rabbitmq --tail 100`
3. Check disk space: `df -h` (RabbitMQ requires disk space)
4. Restart if needed: `docker restart rabbitmq`
5. Verify management UI: http://localhost:15672

**Root Causes**:
- Disk space exhaustion (RabbitMQ blocks on low disk)
- Memory exhaustion (check memory alarms)
- Configuration error
- Corrupted mnesia database

#### TimescaleDBDown
**Impact**: Data persistence unavailable  
**Response**:
1. Check container: `docker ps -a | grep timescaledb`
2. Check logs: `docker logs timescaledb --tail 100`
3. Check disk space: `df -h`
4. Verify postgres-exporter: `docker ps | grep postgres-exporter`
5. Test connection: `docker exec timescaledb psql -U ${TIMESCALE_ADMIN_USER} -d tsdb -c "SELECT version();"`

**Root Causes**:
- Container crash (check logs for panic, OOM)
- Disk full (cannot write WAL)
- Connection limit reached
- Corrupted data files

#### MultipleServicesDown
**Impact**: Platform-wide outage  
**Response**:
1. Check infrastructure: `docker ps` (how many containers down?)
2. Check system resources: `top`, `df -h`, `free -h`
3. Check Docker daemon: `systemctl status docker`
4. Review recent changes (deployments, config updates)
5. Consider full stack restart: `make stop && make deploy`

**Root Causes**:
- System resource exhaustion (CPU, memory, disk)
- Docker daemon issues
- Network problems
- Recent deployment issue

### P2 High Alerts

#### HighEndToEndLatency
**Impact**: Real-time data freshness degraded  
**Response**:
1. Check Message Flow dashboard for bottleneck location
2. Check RabbitMQ queue depth: Is queue growing?
3. Check database write performance: Slow queries?
4. Check resource usage: CPU, memory, disk I/O
5. Consider scaling if resources exhausted

**Investigation Queries**:
```promql
# Identify slowest component
histogram_quantile(0.99, rate(input_latency_ns_bucket{instance="data-injector:4195"}[5m])) / 1e9
histogram_quantile(0.99, rate(output_latency_ns_bucket{instance="timeseries-writer:4195"}[5m])) / 1e9
```

#### HighRabbitMQQueueDepth
**Impact**: Backpressure, writer falling behind  
**Response**:
1. Check RabbitMQ dashboard: Which queues are deep?
2. Check writer throughput: `rate(output_sent{instance="timeseries-writer:4195"}[5m])`
3. Check database performance: Is database slow?
4. Check for writer errors: `rate(output_error{instance="timeseries-writer:4195"}[5m])`
5. Consider scaling writers or optimizing database writes

#### LowMessageDeliveryRate
**Impact**: Potential data loss  
**Response**:
1. Compare ingestion vs. write rates in Platform Overview
2. Check for errors in injector and writer
3. Check RabbitMQ: Messages being lost or expired?
4. Review logs for error patterns
5. Verify end-to-end connectivity

**Investigation**:
```bash
# Check error rates
curl -s "http://localhost:9090/api/v1/query?query=rate(input_error{instance='data-injector:4195'}[5m])" | jq .
curl -s "http://localhost:9090/api/v1/query?query=rate(output_error{instance='timeseries-writer:4195'}[5m])" | jq .
```

### P3 Medium Alerts

#### HighCPUUsage / HighMemoryUsage
**Impact**: System performance degradation likely  
**Response**:
1. Identify top consumers: `docker stats`
2. Check Node dashboard for historical trends
3. Check if this is expected load increase
4. Review capacity planning
5. Consider scaling horizontally or vertically

#### PrometheusLowScrapeSuccessRate
**Impact**: Monitoring blind spots  
**Response**:
1. Check Prometheus targets: http://localhost:9090/targets
2. Identify which targets are failing
3. Check failed target logs
4. Verify network connectivity
5. Fix target configuration if needed

---

## Troubleshooting Playbooks

### No Data in Dashboards

**Symptoms**: Grafana dashboards show "No Data"

**Diagnosis**:
1. Check Prometheus: http://localhost:9090
2. Query for data: `up{job="pathfinder-dev"}`
3. Check Prometheus targets: http://localhost:9090/targets
4. Verify service is running: `docker ps | grep <service>`

**Solutions**:
- Service down → Restart container
- Prometheus not scraping → Check prometheus.yml configuration
- Wrong query → Check dashboard query matches metric names
- Time range issue → Adjust Grafana time range

### High Message Queue Depth

**Symptoms**: RabbitMQ queues growing continuously

**Diagnosis**:
1. Check Message Flow dashboard
2. Compare ingestion rate vs. write rate
3. Check writer container: `docker logs timeseries-writer --tail 100`
4. Check database performance

**Solutions**:
- Writer errors → Fix writer configuration or database connection
- Database slow → Optimize queries, check indexes, increase resources
- Writer under-provisioned → Scale writers
- Burst traffic → Normal, will drain over time

### Pipeline Data Loss

**Symptoms**: Messages ingested but not written to database

**Diagnosis**:
1. Check Platform Overview: Compare MQTT received vs. DB written
2. Check RabbitMQ: Messages in queue but not being consumed?
3. Check writer errors: `curl http://localhost:8405/metrics/benthos-writer | grep error`
4. Check logs for both injector and writer

**Solutions**:
- Writer not consuming → Check RabbitMQ connection, restart writer
- Database connection failed → Check database health, credentials
- Messages expired → Increase RabbitMQ TTL or reduce latency
- Benthos configuration error → Review benthos YAML files

### High Latency

**Symptoms**: P99 latency exceeds SLO (> 5s)

**Diagnosis**:
1. Check Message Flow dashboard: Where is the bottleneck?
2. Check database write latency
3. Check network latency
4. Check resource utilization

**Solutions**:
- Database bottleneck → Optimize queries, add indexes, increase connection pool
- Network issues → Check network I/O, latency between services
- Resource exhaustion → Scale infrastructure
- Large payload → Optimize message size

### Container Won't Start

**Symptoms**: Container in restart loop or exited state

**Diagnosis**:
1. Check status: `docker ps -a | grep <container>`
2. Check logs: `docker logs <container> --tail 100`
3. Check resources: `docker stats` (is system out of resources?)
4. Check configuration: Review environment variables, mounted files

**Solutions**:
- Configuration error → Fix config file, check environment variables
- Dependency not ready → Ensure dependencies started first
- Port conflict → Check port mappings, ensure no conflicts
- Resource limit → Increase Docker resources or system capacity

### Metrics Not Updating

**Symptoms**: Metrics in Grafana are stale

**Diagnosis**:
1. Check Prometheus: Is it scraping? http://localhost:9090/targets
2. Check target service: Is it exposing metrics?
3. Check Grafana datasource: Settings → Data Sources → Prometheus
4. Check dashboard time range

**Solutions**:
- Target down → Restart target service
- Prometheus not scraping → Check prometheus.yml, restart Prometheus
- Grafana can't reach Prometheus → Check network, datasource URL
- Dashboard time range → Adjust to "Last 5 minutes" or enable auto-refresh

---

## Common Operations

### Deploy Platform

```bash
# Full deployment
make deploy MQTT_BROKER=vernemq BROKER=rabbitmq TIMESERIES_DB=timescaledb

# Check status
docker ps
```

### Restart Observability Stack

```bash
# Restart only observability services
docker restart prometheus grafana alertmanager

# Verify
curl http://localhost:9090/-/healthy
curl http://localhost:3000/api/health
curl http://localhost:9093/api/v2/status | jq .
```

### Query Prometheus Manually

```bash
# Check if service is up
curl -s "http://localhost:9090/api/v1/query?query=up{instance='mqtt-broker:8888'}" | jq .

# Get current queue depth
curl -s "http://localhost:9090/api/v1/query?query=sum(rabbitmq_queue_messages)" | jq .

# Get P99 latency
curl -s "http://localhost:9090/api/v1/query?query=histogram_quantile(0.99,rate(input_latency_ns_bucket[5m]))/1e9" | jq .
```

### View Logs

```bash
# View logs for any service
docker logs <service-name> --tail 100 --follow

# Examples
docker logs mqtt-rabbit-injector --tail 100
docker logs rabbitmq-timescale-writer --tail 100
docker logs timescaledb --tail 100
```

### Test Alert

```bash
# Send test alert to Alertmanager
curl -X POST http://localhost:9093/api/v2/alerts \
  -H 'Content-Type: application/json' \
  -d '[{
    "labels": {"alertname": "Test", "severity": "high"},
    "annotations": {"summary": "Test alert"}
  }]'

# Verify in Alertmanager UI
open http://localhost:9093
```

### Silence Alert

```bash
# Via Alertmanager UI
# 1. Go to http://localhost:9093
# 2. Click alert → "Silence" button
# 3. Set duration and reason

# Via API (if amtool installed)
amtool silence add \
  --alertmanager.url=http://localhost:9093 \
  --comment="Maintenance window" \
  --duration=1h \
  alertname="MQTTBrokerDown"
```

### Check HAProxy Metrics

```bash
# All services metrics via HAProxy
curl http://localhost:8405/metrics/vernemq
curl http://localhost:8405/metrics/hivemq
curl http://localhost:8405/metrics/rabbitmq
curl http://localhost:8405/metrics/postgres
curl http://localhost:8405/metrics/minio
curl http://localhost:8405/metrics/benthos-injector
curl http://localhost:8405/metrics/benthos-writer

# HAProxy own stats
curl http://localhost:8405/metrics
```

### Backup Prometheus Data

```bash
# Snapshot Prometheus data
curl -X POST http://localhost:9090/api/v1/admin/tsdb/snapshot

# Copy snapshot (if using persistent volume)
docker exec prometheus ls -la /prometheus/snapshots/
```

### Update Dashboard

```bash
# Edit dashboard JSON file
vim deploy/docker/observability/grafana/provisioning/dashboards/00-OVERVIEW/platform-overview.json

# Restart Grafana to reload
docker restart grafana

# Verify in Grafana UI (may need to refresh browser)
```

### Scale Writer

```bash
# If using Docker Compose scale feature (not currently configured)
docker compose -f <compose-files> up -d --scale timeseries-writer=3

# Verify
docker ps | grep writer
```

### Check Database Size

```bash
# TimescaleDB
docker exec timescaledb psql -U ${TIMESCALE_ADMIN_USER} -d tsdb -c "\l+"

# Query via Prometheus
curl -s "http://localhost:9090/api/v1/query?query=pg_database_size_bytes{datname='tsdb'}" | jq .
```

---

## Maintenance Windows

### Planned Maintenance Procedure

1. **Announce maintenance**: Update status page, notify team
2. **Silence alerts**: Use Alertmanager to silence all alerts for duration
   ```bash
   # Silence all alerts for 1 hour
   curl -X POST http://localhost:9093/api/v2/silences \
     -H 'Content-Type: application/json' \
     -d '{
       "matchers": [{"name": "alertname", "value": ".*", "isRegex": true}],
       "startsAt": "'$(date -u +%Y-%m-%dT%H:%M:%S.000Z)'",
       "endsAt": "'$(date -u -d '+1 hour' +%Y-%m-%dT%H:%M:%S.000Z)'",
       "comment": "Planned maintenance",
       "createdBy": "ops-team"
     }'
   ```
3. **Perform maintenance**: Update services, restart as needed
4. **Verify health**: Check Platform Overview dashboard
5. **Remove silence**: Alerts will automatically resume after silence expires
6. **Announce completion**: Update status page

---

## SLI/SLO Reference

See `SLI-SLO-DEFINITIONS.md` for complete definitions.

**Key SLOs**:
- **Availability**: 99.5% (21.6 min downtime/month)
- **End-to-end latency**: P99 < 5 seconds
- **Message delivery**: 99.9% success rate
- **Database writes**: P99 < 2 seconds

**Error Budget**:
- Monthly budget: 0.5% (216 minutes)
- When < 25% remaining: Feature freeze, focus on reliability
- Burn rate alert: Triggers if will exhaust budget in < 2 days

---

## Additional Resources

- **SLI/SLO Definitions**: `deploy/docker/observability/SLI-SLO-DEFINITIONS.md`
- **Alertmanager Setup**: `deploy/docker/observability/ALERTMANAGER-SETUP.md`
- **Prometheus Configuration**: `deploy/docker/observability/prometheus.yml`
- **Alert Rules**: `deploy/docker/observability/alerting-rules.yml`
- **Project Documentation**: `README.md`, `AGENTS.md`

---

**Last Updated**: 2026-02-06  
**Version**: 1.0  
**Maintainer**: Platform Operations Team
