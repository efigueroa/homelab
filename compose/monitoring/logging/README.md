# Centralized Logging Stack

Grafana Loki + Promtail + Grafana for centralized Docker container log aggregation and visualization.

## Overview

This stack provides centralized logging for all Docker containers in your homelab:

- **Loki**: Log aggregation backend (like Prometheus but for logs)
- **Promtail**: Agent that collects logs from Docker containers
- **Grafana**: Web UI for querying and visualizing logs

### Why This Stack?

- ✅ **Lightweight**: Minimal resource usage compared to ELK stack
- ✅ **Docker-native**: Automatically discovers and collects logs from all containers
- ✅ **Powerful search**: LogQL query language for filtering and searching
- ✅ **Retention**: Configurable log retention (default: 30 days)
- ✅ **Labels**: Automatic labeling by container, image, compose project
- ✅ **Integrated**: Works seamlessly with existing homelab services

## Quick Start

### 1. Configure Environment

```bash
cd ~/homelab/compose/monitoring/logging
nano .env
```

**Update:**
```env
# Change this!
GF_SECURITY_ADMIN_PASSWORD=<your-strong-password>
```

### 2. Deploy the Stack

```bash
docker compose up -d
```

### 3. Access Grafana

Go to: **https://logs.fig.systems**

**Default credentials:**
- Username: `admin`
- Password: `<your GF_SECURITY_ADMIN_PASSWORD>`

**⚠️ Change the password immediately after first login!**

### 4. View Logs

1. Click "Explore" (compass icon) in left sidebar
2. Select "Loki" datasource (should be selected by default)
3. Start querying logs!

## Usage

### Basic Log Queries

**View all logs from a container:**
```logql
{container="jellyfin"}
```

**View logs from a compose project:**
```logql
{compose_project="media"}
```

**View logs from specific service:**
```logql
{compose_service="lldap"}
```

**Filter by log level:**
```logql
{container="immich_server"} |= "error"
```

**Exclude lines:**
```logql
{container="traefik"} != "404"
```

**Multiple filters:**
```logql
{container="jellyfin"} |= "error" != "404"
```

### Advanced Queries

**Count errors per minute:**
```logql
sum(count_over_time({container="jellyfin"} |= "error" [1m])) by (container)
```

**Rate of logs:**
```logql
rate({container="traefik"}[5m])
```

**Logs from last hour:**
```logql
{container="immich_server"} | __timestamp__ >= now() - 1h
```

**Filter by multiple containers:**
```logql
{container=~"jellyfin|immich.*|sonarr"}
```

**Extract and filter JSON:**
```logql
{container="linkwarden"} | json | level="error"
```

## Configuration

### Log Retention

Default: **30 days**

To change retention period:

**Edit `.env`:**
```env
LOKI_RETENTION_PERIOD=60d  # Keep logs for 60 days
```

**Edit `loki-config.yaml`:**
```yaml
limits_config:
  retention_period: 60d  # Must match .env

table_manager:
  retention_period: 60d  # Must match above
```

**Restart:**
```bash
docker compose restart loki
```

### Adjust Resource Limits

**Edit `loki-config.yaml`:**
```yaml
limits_config:
  ingestion_rate_mb: 10          # MB/sec per stream
  ingestion_burst_size_mb: 20    # Burst size
```

### Add Custom Labels

**Edit `promtail-config.yaml`:**
```yaml
scrape_configs:
  - job_name: docker
    docker_sd_configs:
      - host: unix:///var/run/docker.sock

    relabel_configs:
      # Add custom label
      - source_labels: ['__meta_docker_container_label_environment']
        target_label: 'environment'
```

## How It Works

### Architecture

```
Docker Containers
    ↓ (logs via Docker socket)
Promtail (scrapes and ships)
    ↓ (HTTP push)
Loki (stores and indexes)
    ↓ (LogQL queries)
Grafana (visualization)
```

### Log Collection

Promtail automatically collects logs from:
1. **All Docker containers** via Docker socket
2. **System logs** from `/var/log`

Logs are labeled with:
- `container`: Container name
- `image`: Docker image
- `compose_project`: Docker Compose project name
- `compose_service`: Service name from compose.yaml
- `stream`: stdout or stderr

### Storage

Logs are stored in:
- **Location**: `./loki-data/`
- **Format**: Compressed chunks
- **Index**: BoltDB
- **Retention**: Automatic cleanup after retention period

## Integration with Services

### Option 1: Automatic (Default)

Promtail automatically discovers all containers. No changes needed!

### Option 2: Explicit Labels (Recommended)

Add labels to services for better organization:

**Edit any service's `compose.yaml`:**
```yaml
services:
  servicename:
    # ... existing config ...
    labels:
      # ... existing labels ...

      # Add logging labels
      logging: "promtail"
      log_level: "info"
      environment: "production"
```

These labels will be available in Loki for filtering.

### Option 3: Send Logs Directly to Loki

Instead of Promtail scraping, send logs directly:

**Edit service `compose.yaml`:**
```yaml
services:
  servicename:
    # ... existing config ...
    logging:
      driver: loki
      options:
        loki-url: "http://loki:3100/loki/api/v1/push"
        loki-external-labels: "container={{.Name}},compose_project={{.Config.Labels[\"com.docker.compose.project\"]}}"
```

**Note**: This requires the Loki Docker driver plugin (not recommended for simplicity).

## Grafana Dashboards

### Built-in Explore

Best way to start - use Grafana's Explore view:
1. Click "Explore" icon (compass)
2. Select "Loki" datasource
3. Use builder to create queries
4. Save interesting queries

### Pre-built Dashboards

You can import community dashboards:

1. Go to Dashboards → Import
2. Use dashboard ID: `13639` (Docker logs dashboard)
3. Select "Loki" as datasource
4. Import

### Create Custom Dashboard

1. Click "+" → "Dashboard"
2. Add panel
3. Select Loki datasource
4. Build query using LogQL
5. Save dashboard

**Example panels:**
- Error count by container
- Log volume over time
- Top 10 logging containers
- Recent errors table

## Alerting

### Create Log-Based Alerts

1. Go to Alerting → Alert rules
2. Create new alert rule
3. Query: `sum(count_over_time({container="jellyfin"} |= "error" [5m])) > 10`
4. Set thresholds and notification channels
5. Save

**Example alerts:**
- Too many errors in container
- Container restarted
- Disk space warnings
- Failed authentication attempts

## Troubleshooting

### Promtail Not Collecting Logs

**Check Promtail is running:**
```bash
docker logs promtail
```

**Verify Docker socket access:**
```bash
docker exec promtail ls -la /var/run/docker.sock
```

**Test Promtail config:**
```bash
docker exec promtail promtail -config.file=/etc/promtail/config.yaml -dry-run
```

### Loki Not Receiving Logs

**Check Loki health:**
```bash
curl http://localhost:3100/ready
```

**View Loki logs:**
```bash
docker logs loki
```

**Check Promtail is pushing:**
```bash
docker logs promtail | grep -i push
```

### Grafana Can't Connect to Loki

**Test Loki from Grafana container:**
```bash
docker exec grafana wget -O- http://loki:3100/ready
```

**Check datasource configuration:**
- Grafana → Configuration → Data sources → Loki
- URL should be: `http://loki:3100`

### No Logs Appearing

**Wait a few minutes** - logs take time to appear

**Check retention:**
```bash
# Logs older than retention period are deleted
grep retention_period loki-config.yaml
```

**Verify time range in Grafana:**
- Make sure selected time range includes recent logs
- Try "Last 5 minutes"

### High Disk Usage

**Check Loki data size:**
```bash
du -sh ./loki-data
```

**Reduce retention:**
```env
LOKI_RETENTION_PERIOD=7d  # Shorter retention
```

**Manual cleanup:**
```bash
# Stop Loki
docker compose stop loki

# Remove old data (CAREFUL!)
rm -rf ./loki-data/chunks/*

# Restart
docker compose start loki
```

## Performance Tuning

### For Low Resources (< 8GB RAM)

**Edit `loki-config.yaml`:**
```yaml
limits_config:
  retention_period: 7d              # Shorter retention
  ingestion_rate_mb: 5              # Lower rate
  ingestion_burst_size_mb: 10       # Lower burst

query_range:
  results_cache:
    cache:
      embedded_cache:
        max_size_mb: 50             # Smaller cache
```

### For High Volume

**Edit `loki-config.yaml`:**
```yaml
limits_config:
  ingestion_rate_mb: 20             # Higher rate
  ingestion_burst_size_mb: 40       # Higher burst

query_range:
  results_cache:
    cache:
      embedded_cache:
        max_size_mb: 200            # Larger cache
```

## Best Practices

### Log Levels

Configure services to log appropriately:
- **Production**: `info` or `warning`
- **Development**: `debug`
- **Troubleshooting**: `trace`

Too much logging = higher resource usage!

### Retention Strategy

- **Critical services**: 60+ days
- **Normal services**: 30 days
- **High volume services**: 7-14 days

### Query Optimization

- **Use specific labels**: `{container="name"}` not `{container=~".*"}`
- **Limit time range**: Query hours not days when possible
- **Use filters early**: `|= "error"` before parsing
- **Avoid regex when possible**: `|= "string"` faster than `|~ "reg.*ex"`

### Storage Management

Monitor disk usage:
```bash
# Check regularly
du -sh compose/monitoring/logging/loki-data

# Set up alerts when > 80% disk usage
```

## Integration with Homarr

Grafana will automatically appear in Homarr dashboard. You can also:

### Add Grafana Widget to Homarr

1. Edit Homarr dashboard
2. Add "iFrame" widget
3. URL: `https://logs.fig.systems/d/<dashboard-id>`
4. This embeds Grafana dashboards in Homarr

## Backup and Restore

### Backup

```bash
# Backup Loki data
tar czf loki-backup-$(date +%Y%m%d).tar.gz ./loki-data

# Backup Grafana dashboards and datasources
tar czf grafana-backup-$(date +%Y%m%d).tar.gz ./grafana-data ./grafana-provisioning
```

### Restore

```bash
# Restore Loki
docker compose down
tar xzf loki-backup-YYYYMMDD.tar.gz
docker compose up -d

# Restore Grafana
docker compose down
tar xzf grafana-backup-YYYYMMDD.tar.gz
docker compose up -d
```

## Updating

```bash
cd ~/homelab/compose/monitoring/logging

# Pull latest images
docker compose pull

# Restart with new images
docker compose up -d
```

## Resource Usage

**Typical usage:**
- **Loki**: 200-500MB RAM
- **Promtail**: 50-100MB RAM
- **Grafana**: 100-200MB RAM
- **Disk**: ~1-5GB per week (depends on log volume)

## Next Steps

1. ✅ Deploy the stack
2. ✅ Login to Grafana and explore logs
3. ✅ Create useful dashboards
4. ✅ Set up alerts for errors
5. ✅ Configure retention based on needs
6. ⬜ Add Prometheus for metrics (future)
7. ⬜ Add Tempo for distributed tracing (future)

## Resources

- [Loki Documentation](https://grafana.com/docs/loki/latest/)
- [LogQL Query Language](https://grafana.com/docs/loki/latest/logql/)
- [Promtail Configuration](https://grafana.com/docs/loki/latest/clients/promtail/configuration/)
- [Grafana Tutorials](https://grafana.com/tutorials/)

---

**Now you can see logs from all containers in one place!** 🎉
