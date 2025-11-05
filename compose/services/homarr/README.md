# Homarr Dashboard

Modern, customizable dashboard with automatic Docker service discovery.

## Features

- 🎨 **Modern UI** - Beautiful, responsive design
- 🔍 **Auto-Discovery** - Automatically finds Docker services
- 📊 **Widgets** - System stats, weather, calendar, RSS, etc.
- 🏷️ **Labels** - Organize services by category
- 🔗 **Integration** - Connects to *arr apps, Jellyfin, etc.
- 🎯 **Customizable** - Drag-and-drop layout
- 🌙 **Dark Mode** - Built-in dark theme
- 📱 **Mobile Friendly** - Works on all devices

## Access

- **URL:** https://home.fig.systems or https://home.edfig.dev
- **Port:** 7575 (if accessing directly)

## First-Time Setup

### 1. Deploy Homarr

```bash
cd compose/services/homarr
docker compose up -d
```

### 2. Access Dashboard

Open https://home.fig.systems in your browser.

### 3. Auto-Discovery

Homarr will automatically detect services with these labels:

```yaml
labels:
  homarr.name: "Service Name"
  homarr.group: "Category"
  homarr.icon: "/icons/service.png"
  homarr.href: "https://service.fig.systems"
```

## Adding Services to Dashboard

### Automatic (Recommended)

Add labels to your service's `compose.yaml`:

```yaml
labels:
  # Traefik labels...
  traefik.enable: true
  # ... etc

  # Homarr labels
  homarr.name: Jellyfin
  homarr.group: Media
  homarr.icon: https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/jellyfin.png
  homarr.href: https://flix.fig.systems
```

Redeploy the service:
```bash
docker compose up -d
```

Homarr will automatically add it to the dashboard!

### Manual

1. Click the "+" button in Homarr
2. Select "Add Service"
3. Fill in:
   - **Name:** Service name
   - **URL:** https://service.fig.systems
   - **Icon:** Choose from library or custom URL
   - **Category:** Group services (Media, Services, etc.)

## Integration with Services

### Jellyfin

Add to Jellyfin's `compose.yaml`:
```yaml
labels:
  homarr.name: Jellyfin
  homarr.group: Media
  homarr.icon: /icons/jellyfin.png
  homarr.widget.type: jellyfin
  homarr.widget.url: http://jellyfin:8096
  homarr.widget.key: ${JELLYFIN_API_KEY}
```

Shows: Currently playing, library stats

### Sonarr/Radarr

```yaml
labels:
  homarr.name: Sonarr
  homarr.group: Media Automation
  homarr.icon: /icons/sonarr.png
  homarr.widget.type: sonarr
  homarr.widget.url: http://sonarr:8989
  homarr.widget.key: ${SONARR_API_KEY}
```

Shows: Queue, calendar, missing episodes

### qBittorrent

```yaml
labels:
  homarr.name: qBittorrent
  homarr.group: Downloads
  homarr.icon: /icons/qbittorrent.png
  homarr.widget.type: qbittorrent
  homarr.widget.url: http://qbittorrent:8080
  homarr.widget.username: ${QBIT_USERNAME}
  homarr.widget.password: ${QBIT_PASSWORD}
```

Shows: Active torrents, download speed

## Available Widgets

### System Monitoring
- **CPU Usage** - Real-time CPU stats
- **Memory Usage** - RAM usage
- **Disk Space** - Storage capacity
- **Network** - Upload/download speeds

### Services
- **Jellyfin** - Media server stats
- **Sonarr** - TV show automation
- **Radarr** - Movie automation
- **Lidarr** - Music automation
- **Readarr** - Book automation
- **Prowlarr** - Indexer management
- **SABnzbd** - Usenet downloads
- **qBittorrent** - Torrent downloads
- **Overseerr/Jellyseerr** - Media requests

### Utilities
- **Weather** - Local weather forecast
- **Calendar** - Events and tasks
- **RSS Feeds** - News aggregator
- **Docker** - Container status
- **Speed Test** - Internet speed
- **Notes** - Sticky notes
- **Iframe** - Embed any website

## Customization

### Change Theme

1. Click settings icon (⚙️)
2. Go to "Appearance"
3. Choose color scheme
4. Save

### Reorganize Layout

1. Click edit mode (✏️)
2. Drag and drop services
3. Resize widgets
4. Click save

### Add Categories

1. Click "Add Category"
2. Name it (e.g., "Media", "Tools", "Infrastructure")
3. Drag services into categories
4. Collapse/expand as needed

### Custom Icons

**Option 1: Use Icon Library**
- Homarr includes icons from [Dashboard Icons](https://github.com/walkxcode/dashboard-icons)
- Search by service name

**Option 2: Custom URL**
```
https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/service.png
```

**Option 3: Local Icons**
- Place in `./icons/` directory
- Reference as `/icons/service.png`

## Recommended Dashboard Layout

```
┌─────────────────────────────────────────┐
│          🏠 Homelab Dashboard           │
├─────────────────────────────────────────┤
│  [System Stats] [Weather] [Calendar]    │
├─────────────────────────────────────────┤
│  📺 Media                               │
│  [Jellyfin] [Jellyseerr] [Immich]      │
├─────────────────────────────────────────┤
│  🤖 Media Automation                    │
│  [Sonarr] [Radarr] [qBittorrent]       │
├─────────────────────────────────────────┤
│  🛠️ Services                             │
│  [Linkwarden] [Vikunja] [FreshRSS]     │
├─────────────────────────────────────────┤
│  🔧 Infrastructure                      │
│  [Traefik] [LLDAP] [Tinyauth]          │
└─────────────────────────────────────────┘
```

## Add to All Services

To make all your services auto-discoverable, add these labels:

### Jellyfin
```yaml
homarr.name: Jellyfin
homarr.group: Media
homarr.icon: https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/jellyfin.png
```

### Jellyseerr
```yaml
homarr.name: Jellyseerr
homarr.group: Media
homarr.icon: https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/jellyseerr.png
```

### Immich
```yaml
homarr.name: Immich Photos
homarr.group: Media
homarr.icon: https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/immich.png
```

### Sonarr/Radarr/SABnzbd/qBittorrent
```yaml
homarr.name: [Service]
homarr.group: Automation
homarr.icon: https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/[service].png
```

### Linkwarden/Vikunja/etc.
```yaml
homarr.name: [Service]
homarr.group: Utilities
homarr.icon: https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/[service].png
```

## Mobile Access

Homarr is fully responsive. For best mobile experience:

1. Add to home screen (iOS/Android)
2. Works as PWA (Progressive Web App)
3. Touch-optimized interface

## Backup Configuration

### Backup
```bash
cd compose/services/homarr
tar -czf homarr-backup-$(date +%Y%m%d).tar.gz config/ data/
```

### Restore
```bash
cd compose/services/homarr
tar -xzf homarr-backup-YYYYMMDD.tar.gz
docker compose restart
```

## Troubleshooting

### Services not auto-discovered

Check Docker socket permission:
```bash
docker logs homarr
```

Verify labels on service:
```bash
docker inspect service-name | grep homarr
```

### Can't connect to services

Services must be on same Docker network or accessible via hostname.

Use container names, not `localhost`:
- ✅ `http://jellyfin:8096`
- ❌ `http://localhost:8096`

### Widgets not working

1. Check API keys are correct
2. Verify service URLs (use container names)
3. Check service is running: `docker ps`

## Alternatives Considered

| Dashboard | Auto-Discovery | Widgets | Complexity |
|-----------|---------------|---------|------------|
| **Homarr** | ✅ Excellent | ✅ Many | Low |
| Homepage | ✅ Good | ✅ Many | Low |
| Heimdall | ❌ Manual | ❌ Few | Very Low |
| Dashy | ⚠️ Limited | ✅ Some | Medium |
| Homer | ❌ Manual | ❌ None | Very Low |
| Organizr | ⚠️ Limited | ✅ Many | High |

**Homarr chosen for:** Best balance of features, auto-discovery, and ease of use.

## Resources

- [Official Docs](https://homarr.dev/docs)
- [GitHub](https://github.com/ajnart/homarr)
- [Discord Community](https://discord.gg/aCsmEV5RgA)
- [Icon Library](https://github.com/walkxcode/dashboard-icons)

## Tips

1. **Start Simple** - Add core services first, expand later
2. **Use Categories** - Group related services
3. **Enable Widgets** - Make dashboard informative
4. **Mobile First** - Test on phone/tablet
5. **Backup Config** - Save your layout regularly
