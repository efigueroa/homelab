# Centralized Logging with Loki

Guide for setting up and using the centralized logging stack (Loki + Promtail + Grafana).

## Overview

The logging stack provides centralized log aggregation and visualization for all Docker containers:

- **Loki**: Log aggregation backend (stores and indexes logs)
- **Promtail**: Agent that collects logs from Docker containers
- **Grafana**: Web UI for querying and visualizing logs

### Why Centralized Logging?

**Problems without it:**
- Logs scattered across many containers
- Hard to correlate events across services
- Logs lost when containers restart
- No easy way to search historical logs

**Benefits:**
- ✅ Single place to view all logs
- ✅ Powerful search and filtering (LogQL)
- ✅ Persist logs even after container restarts
- ✅ Correlate events across services
- ✅ Create dashboards and alerts
- ✅ Configurable retention (30 days default)

## Quick Setup

### 1. Configure Grafana Password

```bash
cd ~/homelab/compose/monitoring/logging
nano .env
```

**Update:**
```env
GF_SECURITY_ADMIN_PASSWORD=<your-strong-password>
```

**Generate password:**
```bash
openssl rand -base64 20
```

### 2. Deploy

```bash
cd ~/homelab/compose/monitoring/logging
docker compose up -d
```

### 3. Access Grafana

Go to: **https://logs.fig.systems**

**Login:**
- Username: `admin`
- Password: `<your GF_SECURITY_ADMIN_PASSWORD>`

### 4. Start Exploring Logs

1. Click **Explore** (compass icon) in left sidebar
2. Loki datasource should be selected
3. Start querying!

## Basic Usage

### View Logs from a Container

```logql
{container="jellyfin"}
```

### View Last Hour's Logs

```logql
{container="immich_server"} | __timestamp__ >= now() - 1h
```

### Filter for Errors

```logql
{container="traefik"} |= "error"
```

### Exclude Lines

```logql
{container="traefik"} != "404"
```

### Multiple Containers

```logql
{container=~"jellyfin|immich.*"}
```

### By Compose Project

```logql
{compose_project="media"}
```

## Advanced Queries

### Count Errors

```logql
sum(count_over_time({container="jellyfin"} |= "error" [5m]))
```

### Error Rate

```logql
rate({container="traefik"} |= "error" [5m])
```

### Parse JSON Logs

```logql
{container="linkwarden"} | json | level="error"
```

### Top 10 Error Messages

```logql
topk(10,
  sum by (container) (
    count_over_time({job="docker"} |= "error" [24h])
  )
)
```

## Creating Dashboards

### Import Pre-built Dashboard

1. Go to **Dashboards** → **Import**
2. Dashboard ID: **13639** (Docker logs)
3. Select **Loki** as datasource
4. Click **Import**

### Create Custom Dashboard

1. Click **+** → **Dashboard**
2. **Add panel**
3. Select **Loki** datasource
4. Build query
5. Choose visualization (logs, graph, table, etc.)
6. **Save**

**Example panels:**
- Error count by container
- Log volume over time
- Recent errors (table)
- Top logging containers

## Setting Up Alerts

### Create Alert Rule

1. **Alerting** → **Alert rules** → **New alert rule**
2. **Query:**
   ```logql
   sum(count_over_time({container="jellyfin"} |= "error" [5m])) > 10
   ```
3. **Condition**: Alert when > 10 errors in 5 minutes
4. **Configure** notification channel (email, webhook, etc.)
5. **Save**

**Example alerts:**
- Too many errors in service
- Service stopped logging (might have crashed)
- Authentication failures
- Disk space warnings

## Configuration

### Change Log Retention

**Default: 30 days**

Edit `.env`:
```env
LOKI_RETENTION_PERIOD=60d  # 60 days
```

Edit `loki-config.yaml`:
```yaml
limits_config:
  retention_period: 60d

table_manager:
  retention_period: 60d
```

Restart:
```bash
docker compose restart loki
```

### Adjust Resource Limits

For low-resource systems, edit `loki-config.yaml`:

```yaml
limits_config:
  retention_period: 7d              # Shorter retention
  ingestion_rate_mb: 5              # Lower rate

query_range:
  results_cache:
    cache:
      embedded_cache:
        max_size_mb: 50             # Smaller cache
```

### Add Labels to Services

Make services easier to find by adding labels:

**Edit service `compose.yaml`:**
```yaml
services:
  myservice:
    labels:
      logging: "promtail"
      environment: "production"
      tier: "frontend"
```

Query with these labels:
```logql
{environment="production", tier="frontend"}
```

## Troubleshooting

### No Logs Appearing

**Wait a few minutes** - initial log collection takes time

**Check Promtail:**
```bash
docker logs promtail
```

**Check Loki:**
```bash
docker logs loki
```

**Verify Promtail can reach Loki:**
```bash
docker exec promtail wget -O- http://loki:3100/ready
```

### Grafana Can't Connect to Loki

**Test from Grafana:**
```bash
docker exec grafana wget -O- http://loki:3100/ready
```

**Check datasource:** Grafana → Configuration → Data sources → Loki
- URL should be: `http://loki:3100`

### High Disk Usage

**Check size:**
```bash
du -sh compose/monitoring/logging/loki-data
```

**Reduce retention:**
```env
LOKI_RETENTION_PERIOD=7d
```

**Manual cleanup (CAREFUL):**
```bash
docker compose stop loki
rm -rf loki-data/chunks/*
docker compose start loki
```

### Slow Queries

**Optimize queries:**
- Use specific labels: `{container="name"}` not `{container=~".*"}`
- Limit time range: Hours not days
- Filter early: `|= "error"` before parsing
- Avoid complex regex

## Best Practices

### Log Verbosity

Configure appropriate log levels per environment:
- **Production**: `info` or `warning`
- **Debugging**: `debug` or `trace`

Too verbose = wasted resources!

### Retention Strategy

Match retention to importance:
- **Critical services**: 60-90 days
- **Normal services**: 30 days
- **High-volume services**: 7-14 days

### Useful Queries to Save

Create saved queries for common tasks:

**Recent errors:**
```logql
{job="docker"} |= "error" | __timestamp__ >= now() - 15m
```

**Service health check:**
```logql
{container="traefik"} |= "request"
```

**Failed logins:**
```logql
{container="lldap"} |= "failed" |= "login"
```

## Integration Tips

### Embed in Homarr

Add Grafana dashboards to Homarr:

1. Edit Homarr dashboard
2. Add **iFrame widget**
3. URL: `https://logs.fig.systems/d/<dashboard-id>`

### Use with Backups

Include logging data in backups:

```bash
cd ~/homelab/compose/monitoring/logging
tar czf logging-backup-$(date +%Y%m%d).tar.gz loki-data/ grafana-data/
```

### Combine with Metrics

Later you can add Prometheus for metrics:
- Loki for logs
- Prometheus for metrics (CPU, RAM, disk)
- Both in Grafana dashboards

## Common LogQL Patterns

### Filter by Time

```logql
# Last 5 minutes
{container="name"} | __timestamp__ >= now() - 5m

# Specific time range (in Grafana UI time picker)
# Or use: __timestamp__ >= "2024-01-01T00:00:00Z"
```

### Pattern Matching

```logql
# Contains
{container="name"} |= "error"

# Does not contain
{container="name"} != "404"

# Regex match
{container="name"} |~ "error|fail|critical"

# Regex does not match
{container="name"} !~ "debug|trace"
```

### Aggregations

```logql
# Count
count_over_time({container="name"}[5m])

# Rate
rate({container="name"}[5m])

# Sum
sum(count_over_time({job="docker"}[1h])) by (container)

# Average
avg_over_time({container="name"} | unwrap bytes [5m])
```

### JSON Parsing

```logql
# Parse JSON and filter
{container="name"} | json | level="error"

# Extract field
{container="name"} | json | line_format "{{.message}}"

# Filter on JSON field
{container="name"} | json status_code="500"
```

## Resource Usage

**Typical usage:**
- **Loki**: 200-500MB RAM, 1-5GB disk/week
- **Promtail**: 50-100MB RAM
- **Grafana**: 100-200MB RAM, ~100MB disk
- **Total**: ~400-700MB RAM

**For 20 containers with moderate logging**

## Next Steps

1. ✅ Explore your logs in Grafana
2. ✅ Create useful dashboards
3. ✅ Set up alerts for critical errors
4. ⬜ Add Prometheus for metrics (future)
5. ⬜ Add Tempo for distributed tracing (future)
6. ⬜ Create log-based SLA tracking

## Resources

- [Loki Documentation](https://grafana.com/docs/loki/latest/)
- [LogQL Reference](https://grafana.com/docs/loki/latest/logql/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
- [Community Dashboards](https://grafana.com/grafana/dashboards/?search=loki)

---

**Now debug issues 10x faster with centralized logs!** 🔍
