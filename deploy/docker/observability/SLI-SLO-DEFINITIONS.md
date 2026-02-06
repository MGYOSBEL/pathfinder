# Service Level Indicators (SLIs) and Objectives (SLOs)

## Overview

This document defines the Service Level Indicators (SLIs) and Service Level Objectives (SLOs) for the Pathfinder IoT data platform. These metrics ensure the platform meets reliability, performance, and availability requirements for production workloads.

## SLO Measurement Period

All SLOs are measured over a **30-day rolling window** unless otherwise specified.

---

## 1. Service Availability SLOs

### 1.1 MQTT Broker Availability
- **SLI**: Percentage of time MQTT broker responds to health checks
- **SLO**: ≥ 99.5% (43.8 minutes downtime per month maximum)
- **PromQL Query**: 
  ```promql
  avg_over_time(up{instance=~"mqtt-broker:.*"}[30d]) * 100
  ```
- **Alert Threshold**: < 99.5% over 30 days
- **Criticality**: CRITICAL - Blocks all data ingestion

### 1.2 Data Injector Availability
- **SLI**: Percentage of time Benthos injector is running and healthy
- **SLO**: ≥ 99.5% 
- **PromQL Query**:
  ```promql
  avg_over_time(up{instance="data-injector:4195"}[30d]) * 100
  ```
- **Alert Threshold**: < 99.5% over 30 days
- **Criticality**: CRITICAL - Breaks MQTT to RabbitMQ pipeline

### 1.3 RabbitMQ Broker Availability
- **SLI**: Percentage of time RabbitMQ responds to health checks
- **SLO**: ≥ 99.5%
- **PromQL Query**:
  ```promql
  avg_over_time(up{instance="rabbitmq:15692"}[30d]) * 100
  ```
- **Alert Threshold**: < 99.5% over 30 days
- **Criticality**: CRITICAL - Blocks message brokering

### 1.4 Data Writer Availability
- **SLI**: Percentage of time Benthos writer is running and healthy
- **SLO**: ≥ 99.5%
- **PromQL Query**:
  ```promql
  avg_over_time(up{instance="timeseries-writer:4195"}[30d]) * 100
  ```
- **Alert Threshold**: < 99.5% over 30 days
- **Criticality**: CRITICAL - Breaks RabbitMQ to database pipeline

### 1.5 TimescaleDB Availability
- **SLI**: Percentage of time TimescaleDB responds to health checks
- **SLO**: ≥ 99.5%
- **PromQL Query**:
  ```promql
  avg_over_time(up{instance="postgres-exporter:9187"}[30d]) * 100
  ```
- **Alert Threshold**: < 99.5% over 30 days
- **Criticality**: CRITICAL - Blocks data persistence

### 1.6 MinIO Storage Availability
- **SLI**: Percentage of time MinIO S3 storage responds to health checks
- **SLO**: ≥ 99.0% (lower than critical services)
- **PromQL Query**:
  ```promql
  avg_over_time(up{instance="minio:9000"}[30d]) * 100
  ```
- **Alert Threshold**: < 99.0% over 30 days
- **Criticality**: HIGH - Used for Loki log storage

---

## 2. Performance SLOs

### 2.1 End-to-End Message Processing Latency
- **SLI**: 99th percentile latency from MQTT publish to database write
- **SLO**: p99 < 5 seconds
- **Measurement**: Combined input and output latency across injector and writer
- **PromQL Query**:
  ```promql
  histogram_quantile(0.99, 
    rate(input_latency_ns_bucket{instance="data-injector:4195"}[5m])
  ) / 1e9 +
  histogram_quantile(0.99, 
    rate(output_latency_ns_bucket{instance="timeseries-writer:4195"}[5m])
  ) / 1e9
  ```
- **Alert Threshold**: > 5 seconds for 10 consecutive minutes
- **Criticality**: HIGH - Impacts real-time data freshness

### 2.2 Database Write Latency
- **SLI**: 99th percentile write latency to TimescaleDB
- **SLO**: p99 < 2 seconds
- **PromQL Query**:
  ```promql
  histogram_quantile(0.99, 
    rate(output_latency_ns_bucket{instance="timeseries-writer:4195"}[5m])
  ) / 1e9
  ```
- **Alert Threshold**: > 2 seconds for 10 consecutive minutes
- **Criticality**: HIGH - Indicates database performance issues

### 2.3 MQTT Message Ingestion Rate
- **SLI**: Messages per second ingested by MQTT broker
- **SLO**: Support ≥ 1,000 messages/second sustained throughput
- **PromQL Query**:
  ```promql
  rate(vernemq_mqtt_publish_received[5m]) or 
  rate(com_hivemq_messages_incoming_publish_rate[5m])
  ```
- **Alert Threshold**: < 100 msg/s for 15 minutes (possible issue)
- **Criticality**: MEDIUM - Performance degradation indicator

### 2.4 RabbitMQ Queue Processing
- **SLI**: Average queue depth over time
- **SLO**: Average queue depth < 1,000 messages
- **PromQL Query**:
  ```promql
  avg_over_time(sum(rabbitmq_queue_messages{job="pathfinder-dev"})[5m])
  ```
- **Alert Threshold**: > 10,000 messages for 5 consecutive minutes
- **Criticality**: HIGH - Indicates backpressure or writer issues

---

## 3. Reliability SLOs

### 3.1 Message Delivery Success Rate
- **SLI**: Percentage of messages successfully written to database
- **SLO**: ≥ 99.9% success rate
- **Measurement**: (Messages written / Messages received) * 100
- **PromQL Query**:
  ```promql
  (
    sum(rate(output_sent{instance="timeseries-writer:4195"}[5m]))
    /
    sum(rate(input_received{instance="data-injector:4195"}[5m]))
  ) * 100
  ```
- **Alert Threshold**: < 99.0% over 15 minutes
- **Criticality**: CRITICAL - Data loss indicator

### 3.2 Database Write Success Rate
- **SLI**: Percentage of write attempts that succeed
- **SLO**: ≥ 99.5% success rate
- **Measurement**: 100% - (error rate / total attempts)
- **PromQL Query**:
  ```promql
  (1 - (
    sum(rate(output_error{instance="timeseries-writer:4195"}[5m]))
    /
    sum(rate(output_sent{instance="timeseries-writer:4195"}[5m]))
  )) * 100
  ```
- **Alert Threshold**: < 99.0% over 10 minutes
- **Criticality**: CRITICAL - Data persistence issues

### 3.3 MQTT Connection Stability
- **SLI**: Connection failure and loss rate
- **SLO**: < 0.1% connection failures
- **PromQL Query**:
  ```promql
  rate(input_connection_failed{instance="data-injector:4195"}[5m]) +
  rate(input_connection_lost{instance="data-injector:4195"}[5m])
  ```
- **Alert Threshold**: > 10 failures per minute
- **Criticality**: HIGH - Network or MQTT broker issues

### 3.4 RabbitMQ Connection Stability
- **SLI**: Broker connection failure rate
- **SLO**: < 0.1% connection failures
- **PromQL Query**:
  ```promql
  rate(output_connection_failed{instance="data-injector:4195"}[5m]) +
  rate(input_connection_failed{instance="timeseries-writer:4195"}[5m])
  ```
- **Alert Threshold**: > 5 failures per minute
- **Criticality**: HIGH - Broker connectivity issues

---

## 4. Resource Utilization SLOs

### 4.1 System CPU Utilization
- **SLI**: Average CPU usage across all cores
- **SLO**: < 85% sustained usage
- **PromQL Query**:
  ```promql
  100 - (avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)
  ```
- **Alert Threshold**: > 85% for 10 consecutive minutes
- **Criticality**: MEDIUM - Capacity planning indicator

### 4.2 System Memory Utilization
- **SLI**: Memory usage percentage
- **SLO**: < 90% sustained usage
- **PromQL Query**:
  ```promql
  (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100
  ```
- **Alert Threshold**: > 90% for 5 consecutive minutes
- **Criticality**: HIGH - Risk of OOM kills

### 4.3 Disk Space Utilization
- **SLI**: Root filesystem usage percentage
- **SLO**: < 85% usage
- **PromQL Query**:
  ```promql
  100 - ((node_filesystem_avail_bytes{mountpoint="/"} / 
         node_filesystem_size_bytes{mountpoint="/"}) * 100)
  ```
- **Alert Threshold**: > 85% usage
- **Criticality**: HIGH - Risk of write failures

### 4.4 TimescaleDB Database Size
- **SLI**: Total database size growth rate
- **SLO**: < 50 GB per month growth (adjust based on capacity)
- **PromQL Query**:
  ```promql
  sum(pg_database_size_bytes{datname="tsdb"}) / 1024 / 1024 / 1024
  ```
- **Alert Threshold**: Monitor for capacity planning
- **Criticality**: MEDIUM - Capacity planning

---

## 5. Observability Stack SLOs

### 5.1 Prometheus Scrape Success Rate
- **SLI**: Percentage of successful metric scrapes
- **SLO**: ≥ 99.0% scrape success
- **PromQL Query**:
  ```promql
  avg_over_time(
    sum(up{job="pathfinder-dev"} == 1) / 
    count(up{job="pathfinder-dev"})
  [5m]) * 100
  ```
- **Alert Threshold**: < 95% over 15 minutes
- **Criticality**: MEDIUM - Monitoring blind spots

### 5.2 Grafana Availability
- **SLI**: Grafana uptime percentage
- **SLO**: ≥ 99.0%
- **PromQL Query**:
  ```promql
  avg_over_time(up{job="grafana"}[30d]) * 100
  ```
- **Alert Threshold**: Down for 5+ minutes
- **Criticality**: MEDIUM - Dashboard access

---

## 6. SLO Error Budget

Based on 99.5% availability SLO:
- **Monthly Error Budget**: 21.6 minutes downtime
- **Weekly Error Budget**: 5.04 minutes downtime
- **Daily Error Budget**: 43.2 seconds downtime

### Error Budget Policy
- **100-75% remaining**: Normal operations, focus on features
- **75-50% remaining**: Increased monitoring, review incidents
- **50-25% remaining**: Focus on reliability, defer non-critical changes
- **< 25% remaining**: Feature freeze, all hands on stability

---

## 7. Monitoring and Alerting Strategy

### Alert Severity Levels

**CRITICAL (P1)**:
- Service completely down
- Data loss occurring
- > 50% error budget consumed in < 1 hour
- **Response Time**: Immediate (15 minutes)
- **Escalation**: After 30 minutes

**HIGH (P2)**:
- Degraded performance affecting users
- Error rates above threshold
- Resource exhaustion imminent
- **Response Time**: Within 1 hour
- **Escalation**: After 2 hours

**MEDIUM (P3)**:
- Non-critical performance degradation
- Elevated error rates within budget
- Capacity planning alerts
- **Response Time**: Within 4 hours
- **Escalation**: After 8 hours

**LOW (P4)**:
- Informational alerts
- Trend analysis
- Proactive monitoring
- **Response Time**: Next business day

### Alert Fatigue Prevention
1. Use multi-window, multi-burn-rate alerts
2. Implement alert grouping by service
3. Set appropriate alert thresholds with hysteresis
4. Regular alert tuning based on false positive rate
5. Target < 5% false positive rate per alert

---

## 8. SLO Review and Adjustment

### Review Schedule
- **Weekly**: Error budget consumption review
- **Monthly**: SLO achievement review and trend analysis
- **Quarterly**: SLO adjustment based on business needs
- **Annually**: Comprehensive SLI/SLO framework review

### Adjustment Criteria
- Consistent over-achievement (> 99.9% when SLO is 99.5%)
- Consistent under-achievement despite efforts
- Changes in business requirements
- Platform architecture changes
- User feedback and expectations

---

## 9. Dashboard Access

### Primary Dashboards
- **Platform Overview**: http://localhost:3000/d/platform-overview
- **Observability Stack Health**: http://localhost:3000/d/observability-stack-health
- **Infrastructure Overview**: http://localhost:3000/d/infrastructure-overview
- **Message Flow Overview**: http://localhost:3000/d/message-flow-overview

### SLO Tracking
Create dedicated SLO tracking dashboard using queries defined in this document.

---

## 10. References

- [Google SRE Book - Service Level Objectives](https://sre.google/sre-book/service-level-objectives/)
- [Prometheus Best Practices - Alerting](https://prometheus.io/docs/practices/alerting/)
- [Multi-window, multi-burn-rate alerts](https://sre.google/workbook/alerting-on-slos/)

---

**Last Updated**: 2026-02-06  
**Version**: 1.0  
**Owner**: Platform Engineering Team
