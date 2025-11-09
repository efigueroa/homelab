# Uptime Kuma - Status & Uptime Monitoring

Beautiful uptime monitoring and alerting for all your homelab services.

## Overview

**Uptime Kuma** monitors the health and uptime of your services:

- ✅ **HTTP(s) Monitoring**: Check if web services are responding
- ✅ **TCP Port Monitoring**: Check if services are listening on ports
- ✅ **Docker Container Monitoring**: Check container status
- ✅ **Response Time**: Measure how fast services respond
- ✅ **SSL Certificate Monitoring**: Alert before certificates expire
- ✅ **Status Pages**: Public or private status pages
- ✅ **Notifications**: Email, Discord, Slack, Pushover, and 90+ more
- ✅ **Beautiful UI**: Clean, modern interface

## Quick Start

### 1. Deploy

```bash
cd ~/homelab/compose/monitoring/uptime
docker compose up -d
```

### 2. Access Web UI

Go to: **https://status.fig.systems**

### 3. Create Admin Account

On first visit, you'll be prompted to create an admin account:
- Username: `admin` (or your choice)
- Password: Strong password
- Click "Create"

### 4. Add Your First Monitor

Click **"Add New Monitor"**

**Example: Monitor Jellyfin**
- Monitor Type: `HTTP(s)`
- Friendly Name: `Jellyfin`
- URL: `https://flix.fig.systems`
- Heartbeat Interval: `60` seconds
- Retries: `3`
- Click **Save**

Uptime Kuma will now check Jellyfin every 60 seconds!

## Monitoring Your Services

### Quick Setup All Services

Here's a template for all your homelab services:

**Core Services:**
```
Name: Traefik Dashboard
Type: HTTP(s)
URL: https://traefik.fig.systems
Interval: 60s

Name: LLDAP
Type: HTTP(s)
URL: https://lldap.fig.systems
Interval: 60s

Name: Grafana Logs
Type: HTTP(s)
URL: https://logs.fig.systems
Interval: 60s
```

**Media Services:**
```
Name: Jellyfin
Type: HTTP(s)
URL: https://flix.fig.systems
Interval: 60s

Name: Immich
Type: HTTP(s)
URL: https://photos.fig.systems
Interval: 60s

Name: Jellyseerr
Type: HTTP(s)
URL: https://requests.fig.systems
Interval: 60s

Name: Sonarr
Type: HTTP(s)
URL: https://sonarr.fig.systems
Interval: 60s

Name: Radarr
Type: HTTP(s)
URL: https://radarr.fig.systems
Interval: 60s
```

**Utility Services:**
```
Name: Homarr Dashboard
Type: HTTP(s)
URL: https://home.fig.systems
Interval: 60s

Name: Backrest
Type: HTTP(s)
URL: https://backup.fig.systems
Interval: 60s

Name: Linkwarden
Type: HTTP(s)
URL: https://links.fig.systems
Interval: 60s

Name: Vikunja
Type: HTTP(s)
URL: https://tasks.fig.systems
Interval: 60s
```

### Advanced Monitoring Options

#### Monitor Docker Containers Directly

**Setup:**
1. Add New Monitor
2. Type: **Docker Container**
3. Docker Daemon: `unix:///var/run/docker.sock`
4. Container Name: `jellyfin`
5. Click Save

**Benefits:**
- Checks if container is running
- Monitors container restarts
- No network requests needed

**Note**: Requires mounting Docker socket (already configured).

#### Monitor TCP Ports

**Example: Monitor PostgreSQL**
```
Type: TCP Port
Hostname: linkwarden-postgres
Port: 5432
Interval: 60s
```

#### Check SSL Certificates

**Automatic**: When using HTTP(s) monitors, Uptime Kuma automatically:
- Checks SSL certificate validity
- Alerts when certificate expires soon (7 days default)
- Shows certificate expiry date

#### Keyword Monitoring

Check if a page contains specific text:

```
Type: HTTP(s) - Keyword
URL: https://home.fig.systems
Keyword: "Homarr"  # Check page contains "Homarr"
```

## Notifications

### Setup Alerts

1. Click **Settings** (gear icon)
2. Click **Notifications**
3. Click **Setup Notification**

### Popular Options

#### Email
```
Type: Email (SMTP)
Host: smtp.gmail.com
Port: 587
Security: TLS
Username: your-email@gmail.com
Password: your-app-password
From: alerts@yourdomain.com
To: you@email.com
```

#### Discord
```
Type: Discord
Webhook URL: https://discord.com/api/webhooks/...
(Get from Discord Server Settings → Integrations → Webhooks)
```

#### Slack
```
Type: Slack
Webhook URL: https://hooks.slack.com/services/...
(Get from Slack App → Incoming Webhooks)
```

#### Pushover (Mobile)
```
Type: Pushover
User Key: (from Pushover account)
App Token: (create app in Pushover)
Priority: Normal
```

#### Gotify (Self-hosted)
```
Type: Gotify
Server URL: https://gotify.yourdomain.com
App Token: (from Gotify)
Priority: 5
```

### Apply to Monitors

After setting up notification:
1. Edit a monitor
2. Scroll to **Notifications**
3. Select your notification method
4. Click **Save**

Or apply to all monitors:
1. Settings → Notifications
2. Click **Apply on all existing monitors**

## Status Pages

### Create Public Status Page

Perfect for showing service status to family/friends!

**Setup:**
1. Click **Status Pages**
2. Click **Add New Status Page**
3. **Slug**: `homelab` (creates /status/homelab)
4. **Title**: `Homelab Status`
5. **Description**: `Status of all homelab services`
6. Click **Next**

**Add Services:**
1. Drag monitors into "Public" or "Groups"
2. Organize by category (Core, Media, Utilities)
3. Click **Save**

**Access:**
- Private: https://status.fig.systems/status/homelab
- Or make public (no login required)

**Share with family:**
```
https://status.fig.systems/status/homelab
```

### Customize Status Page

**Options:**
- Show/hide uptime percentage
- Show/hide response time
- Custom domain
- Theme (light/dark/auto)
- Custom CSS
- Password protection

## Tags and Groups

### Organize Monitors with Tags

**Create Tags:**
1. Click **Manage Tags**
2. Add tags like:
   - `core`
   - `media`
   - `critical`
   - `production`

**Apply to Monitors:**
1. Edit monitor
2. Scroll to **Tags**
3. Select tags
4. Save

**Filter by Tag:**
- Click tag name to show only those monitors

### Create Monitor Groups

**Group by service type:**
1. Settings → Groups
2. Create groups:
   - Core Infrastructure
   - Media Services
   - Productivity
   - Monitoring

Drag monitors into groups for organization.

## Maintenance Windows

### Schedule Maintenance

Pause notifications during planned downtime:

1. Edit monitor
2. Click **Maintenance**
3. **Add Maintenance**
4. Set start/end time
5. Select monitors
6. Save

During maintenance:
- Monitor still checks but doesn't alert
- Status page shows "In Maintenance"

## Best Practices

### Monitor Configuration

**Heartbeat Interval:**
- Critical services: 30-60 seconds
- Normal services: 60-120 seconds
- Background jobs: 300-600 seconds

**Retries:**
- Set to 2-3 to avoid false positives
- Service must fail 2-3 times before alerting

**Timeout:**
- Web services: 10-30 seconds
- APIs: 5-10 seconds
- Slow services: 30-60 seconds

### What to Monitor

**Critical (Monitor these!):**
- ✅ Traefik (if this is down, everything is down)
- ✅ LLDAP (SSO depends on this)
- ✅ Core services users depend on

**Important:**
- ✅ Jellyfin, Immich (main media services)
- ✅ Sonarr, Radarr (automation)
- ✅ Backrest (backups)

**Nice to have:**
- ⬜ Utility services
- ⬜ Less critical services

**Don't over-monitor:**
- Internal components (databases, redis, etc.)
- These should be monitored via main service health

### Notification Strategy

**Alert fatigue is real!**

**Good approach:**
- Critical services → Immediate push notification
- Important services → Email
- Nice-to-have → Email digest

**Don't:**
- Alert on every blip
- Send all alerts to mobile push
- Alert on expected downtime

## Integration with Loki

Uptime Kuma and Loki complement each other:

**Uptime Kuma:**
- ✅ Is the service UP or DOWN?
- ✅ How long was it down?
- ✅ Response time trends

**Loki:**
- ✅ WHY did it go down?
- ✅ What errors happened?
- ✅ Historical log analysis

**Workflow:**
1. Uptime Kuma alerts you: "Jellyfin is down!"
2. Go to Grafana/Loki
3. Query: `{container="jellyfin"} | __timestamp__ >= now() - 15m`
4. See what went wrong

## Metrics and Graphs

### Built-in Metrics

Uptime Kuma tracks:
- **Uptime %**: 99.9%, 99.5%, etc.
- **Response Time**: Average, min, max
- **Ping**: Latency to service
- **Certificate Expiry**: Days until SSL expires

### Response Time Graph

Click any monitor to see:
- 24-hour response time graph
- Uptime/downtime periods
- Recent incidents

### Export Data

Export uptime data:
1. Settings → Backup
2. Export JSON (includes all monitors and data)
3. Store backup safely

## Troubleshooting

### Monitor Shows Down But Service Works

**Check:**
1. **SSL Certificate**: Is it valid?
2. **SSO**: Does monitor need to login first?
3. **Timeout**: Is timeout too short?
4. **Network**: Can Uptime Kuma reach the service?

**Solutions:**
- Increase timeout
- Check accepted status codes (200-299)
- Verify URL is correct
- Check Uptime Kuma logs: `docker logs uptime-kuma`

### Docker Container Monitor Not Working

**Requirements:**
- Docker socket must be mounted (✅ already configured)
- Container name must be exact

**Test:**
```bash
docker exec uptime-kuma ls /var/run/docker.sock
# Should show the socket file
```

### Notifications Not Sending

**Check:**
1. Test notification in Settings → Notifications
2. Check Uptime Kuma logs
3. Verify notification service credentials
4. Check if notification is enabled on monitor

### Can't Access Web UI

**Check:**
```bash
# Container running?
docker ps | grep uptime-kuma

# Logs
docker logs uptime-kuma

# Traefik routing
docker logs traefik | grep uptime
```

## Advanced Features

### API Access

Uptime Kuma has a WebSocket API:

**Get API Key:**
1. Settings → API Keys
2. Generate new key
3. Use with monitoring tools

### Docker Socket Monitoring

Already configured! You can monitor:
- Container status (running/stopped)
- Container restarts
- Resource usage (via Docker stats)

### Multiple Status Pages

Create different status pages:
- `/status/public` - For family/friends
- `/status/critical` - Only critical services
- `/status/media` - Media services only

### Custom CSS

Brand your status page:
1. Status Page → Edit
2. Custom CSS
3. Add styling

**Example:**
```css
body {
  background: #1a1a1a;
}
.title {
  color: #00ff00;
}
```

## Resource Usage

**Typical usage:**
- **RAM**: 50-150MB
- **CPU**: Very low (only during checks)
- **Disk**: <100MB
- **Network**: Minimal (only during checks)

**Very lightweight!**

## Backup and Restore

### Backup

**Automatic backup:**
1. Settings → Backup
2. Export

**Manual backup:**
```bash
cd ~/homelab/compose/monitoring/uptime
tar czf uptime-backup-$(date +%Y%m%d).tar.gz ./data
```

### Restore

```bash
docker compose down
tar xzf uptime-backup-YYYYMMDD.tar.gz
docker compose up -d
```

## Comparison: Uptime Kuma vs Loki

| Feature | Uptime Kuma | Loki |
|---------|-------------|------|
| **Purpose** | Uptime monitoring | Log aggregation |
| **Checks** | HTTP, TCP, Ping, Docker | Logs only |
| **Alerts** | Service down, slow | Log patterns |
| **Response Time** | ✅ Yes | ❌ No |
| **Uptime %** | ✅ Yes | ❌ No |
| **SSL Monitoring** | ✅ Yes | ❌ No |
| **Why Service Down** | ❌ No | ✅ Yes (via logs) |
| **Historical Logs** | ❌ No | ✅ Yes |
| **Status Pages** | ✅ Yes | ❌ No |

**Use both together!**
- Uptime Kuma tells you WHAT is down
- Loki tells you WHY it went down

## Next Steps

1. ✅ Deploy Uptime Kuma
2. ✅ Add monitors for all services
3. ✅ Set up notifications (Email, Discord, etc.)
4. ✅ Create status page
5. ✅ Test alerts by stopping a service
6. ⬜ Share status page with family
7. ⬜ Set up maintenance windows
8. ⬜ Review and tune check intervals

## Resources

- [Uptime Kuma GitHub](https://github.com/louislam/uptime-kuma)
- [Uptime Kuma Wiki](https://github.com/louislam/uptime-kuma/wiki)
- [Notification Services List](https://github.com/louislam/uptime-kuma/wiki/Notification-Services)

---

**Know instantly when something goes down!** 🚨
