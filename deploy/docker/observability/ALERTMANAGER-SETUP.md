# Alertmanager Configuration Guide

## Environment Variables

Alertmanager uses environment variables for notification channel configuration. Add these to your `.env` file:

### SMTP Email Configuration

```bash
# SMTP server configuration
SMTP_HOST=smtp.gmail.com:587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-specific-password

# Email addresses
ALERT_EMAIL_FROM=alerts@pathfinder.local
ALERT_EMAIL_CRITICAL=ops-oncall@company.com
ALERT_EMAIL_HIGH=ops-team@company.com
```

#### Gmail Setup
1. Enable 2-factor authentication
2. Generate app-specific password: https://myaccount.google.com/apppasswords
3. Use app password in `SMTP_PASSWORD`

#### Other SMTP Providers
- **Office 365**: `smtp.office365.com:587`
- **SendGrid**: `smtp.sendgrid.net:587`
- **AWS SES**: `email-smtp.us-east-1.amazonaws.com:587`

### Slack Configuration

```bash
# Slack webhook URL (required for Slack notifications)
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL

# Slack channels for different severity levels
SLACK_CRITICAL_CHANNEL=#alerts-critical
SLACK_HIGH_CHANNEL=#alerts-high
SLACK_MEDIUM_CHANNEL=#alerts-medium
```

#### Slack Webhook Setup
1. Go to: https://api.slack.com/apps
2. Create New App → From scratch
3. Enable "Incoming Webhooks"
4. Click "Add New Webhook to Workspace"
5. Select channel and authorize
6. Copy webhook URL to `SLACK_WEBHOOK_URL`

### PagerDuty Configuration (Optional)

For PagerDuty integration, add to `alertmanager.yml`:

```yaml
pagerduty_configs:
  - routing_key: '${PAGERDUTY_ROUTING_KEY}'
    severity: '{{ .GroupLabels.severity }}'
    description: '{{ .Annotations.summary }}'
```

Then add to `.env`:
```bash
PAGERDUTY_ROUTING_KEY=your-integration-key
```

## Alert Routing

Alerts are routed by severity:

### Critical (P1)
- **Channels**: Slack + Email
- **Group wait**: 10 seconds
- **Repeat**: Every 30 minutes
- **Example**: Service down, data loss risk

### High (P2)
- **Channels**: Slack + Email
- **Group wait**: 30 seconds
- **Repeat**: Every 2 hours
- **Example**: Performance degradation, connection failures

### Medium (P3)
- **Channels**: Slack only
- **Group wait**: 5 minutes
- **Repeat**: Every 12 hours
- **Example**: Capacity warnings, observability issues

## Testing Notifications

### 1. Start Alertmanager
```bash
make deploy MQTT_BROKER=vernemq BROKER=rabbitmq TIMESERIES_DB=timescaledb
```

### 2. Check Alertmanager status
```bash
# Web UI
http://localhost:9093

# API status
curl http://localhost:9093/api/v2/status
```

### 3. Send test alert
```bash
curl -X POST http://localhost:9093/api/v1/alerts \
  -H 'Content-Type: application/json' \
  -d '[
    {
      "labels": {
        "alertname": "TestAlert",
        "severity": "high",
        "component": "test",
        "environment": "pathfinder-dev"
      },
      "annotations": {
        "summary": "This is a test alert",
        "description": "Testing Alertmanager notification channels",
        "impact": "No impact - this is a test",
        "action": "Verify notification received"
      }
    }
  ]'
```

### 4. Verify notifications
- Check Slack channel for message
- Check email inbox
- Check Alertmanager UI: http://localhost:9093/#/alerts

## Silence Alerts

### Via Web UI
1. Go to http://localhost:9093
2. Click alert → "Silence" button
3. Set duration and reason

### Via CLI
```bash
# Silence all alerts for maintenance window (1 hour)
amtool silence add \
  --alertmanager.url=http://localhost:9093 \
  --comment="Scheduled maintenance" \
  --duration=1h \
  alertname=~".+"

# Silence specific component
amtool silence add \
  --alertmanager.url=http://localhost:9093 \
  --comment="Database maintenance" \
  --duration=30m \
  component=timescaledb
```

## Troubleshooting

### No notifications received

1. **Check Alertmanager logs:**
```bash
docker logs alertmanager
```

2. **Verify environment variables:**
```bash
docker exec alertmanager env | grep -E "(SMTP|SLACK)"
```

3. **Test SMTP connection:**
```bash
# Install swaks for testing
docker exec alertmanager /bin/sh -c "
  echo 'Test email' | \
  mail -s 'Test Subject' \
  -S smtp=smtp.gmail.com:587 \
  -S smtp-use-starttls \
  -S smtp-auth=login \
  -S smtp-auth-user=${SMTP_USERNAME} \
  -S smtp-auth-password=${SMTP_PASSWORD} \
  test@example.com
"
```

4. **Test Slack webhook:**
```bash
curl -X POST ${SLACK_WEBHOOK_URL} \
  -H 'Content-Type: application/json' \
  -d '{"text": "Test message from Alertmanager"}'
```

### Alerts not firing

1. **Check Prometheus alerts:**
```bash
# List all alerts
curl http://localhost:9090/api/v1/alerts

# Check specific alert
curl 'http://localhost:9090/api/v1/query?query=ALERTS{alertname="MQTTBrokerDown"}'
```

2. **Verify Alertmanager connection:**
```bash
curl http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | select(.job=="alertmanager")'
```

3. **Check Prometheus config:**
```bash
docker exec prometheus cat /etc/prometheus/prometheus.yml | grep -A5 "alerting:"
```

## Configuration Files

- **Alertmanager config**: `deploy/docker/observability/alertmanager.yml`
- **Prometheus alerting**: `deploy/docker/observability/prometheus.yml`
- **Alert rules**: `deploy/docker/observability/alerting-rules.yml`
- **Environment variables**: `.env` (create from template below)

## .env Template

```bash
# SMTP Configuration
SMTP_HOST=smtp.gmail.com:587
SMTP_USERNAME=
SMTP_PASSWORD=
ALERT_EMAIL_FROM=alerts@pathfinder.local
ALERT_EMAIL_CRITICAL=ops-oncall@company.com
ALERT_EMAIL_HIGH=ops-team@company.com

# Slack Configuration
SLACK_WEBHOOK_URL=
SLACK_CRITICAL_CHANNEL=#alerts-critical
SLACK_HIGH_CHANNEL=#alerts-high
SLACK_MEDIUM_CHANNEL=#alerts-medium
```

## Next Steps

1. Configure SMTP credentials in `.env`
2. Set up Slack webhook in `.env`
3. Deploy stack: `make deploy`
4. Test notifications with test alert
5. Monitor Alertmanager UI: http://localhost:9093
6. Configure silences for maintenance windows as needed
