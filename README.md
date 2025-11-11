# Homelab GitOps Configuration

This repository contains Docker Compose configurations for self-hosted home services.

## 💻 Hardware Specifications

- **Host**: Proxmox VE 9 (Debian 13)
  - CPU: AMD Ryzen 5 7600X (6 cores, 12 threads, up to 5.3 GHz)
  - GPU: NVIDIA GeForce GTX 1070 (8GB VRAM)
  - RAM: 32GB DDR5

- **VM**: AlmaLinux 9.6 (RHEL 9 compatible)
  - CPU: 8 vCPUs
  - RAM: 24GB
  - Storage: 500GB+ (expandable)
  - GPU: GTX 1070 (PCIe passthrough)

**Documentation:**
- [Complete Architecture Guide](docs/architecture.md) - Integration, networking, logging, GPU setup
- [AlmaLinux VM Setup](docs/setup/almalinux-vm.md) - Full installation and configuration guide

## 🏗️ Infrastructure

### Core Services (Port 80/443)
- **Traefik** - Reverse proxy with automatic Let's Encrypt SSL
- **LLDAP** - Lightweight LDAP server for user management
  - Admin: `edfig` (admin@edfig.dev)
  - Web UI: https://lldap.fig.systems
- **Tinyauth** - SSO authentication via Traefik forward auth
  - Connected to LLDAP for user authentication
  - Web UI: https://auth.fig.systems

## 📁 Directory Structure

```
compose/
├── core/           # Infrastructure services
│   ├── traefik/    # Reverse proxy & SSL
│   ├── lldap/      # LDAP user directory
│   └── tinyauth/   # SSO authentication
├── media/          # Media services
│   ├── frontend/   # Media frontends
│   │   ├── jellyfin/   # Media server (flix.fig.systems)
│   │   ├── jellyseer/  # Request management (requests.fig.systems)
│   │   └── immich/     # Photo management (photos.fig.systems)
│   └── automation/ # Media automation
│       ├── sonarr/     # TV show management
│       ├── radarr/     # Movie management
│       ├── sabnzbd/    # Usenet downloader
│       ├── qbittorrent/# Torrent client
│       ├── recyclarr/  # TRaSH Guides sync
│       └── profilarr/  # Profile manager (profilarr.fig.systems)
├── monitoring/      # Monitoring & logging
│   ├── logging/     # Centralized logging stack
│   │   ├── loki/        # Log aggregation (loki.fig.systems)
│   │   ├── promtail/    # Log collection agent
│   │   └── grafana/     # Log visualization (logs.fig.systems)
│   └── uptime/      # Uptime monitoring
│       └── uptime-kuma/ # Status & uptime monitoring (status.fig.systems)
└── services/       # Utility services
    ├── homarr/         # Dashboard (home.fig.systems)
    ├── backrest/       # Backup manager (backup.fig.systems)
    ├── static-sites/   # Static websites (Caddy)
    ├── karakeep/       # Bookmark manager with AI (links.fig.systems)
    ├── ollama/         # Local LLM server (ollama.fig.systems)
    ├── vikunja/        # Task management (tasks.fig.systems)
    ├── lubelogger/     # Vehicle tracker (garage.fig.systems)
    ├── calibre-web/    # Ebook library (books.fig.systems)
    ├── booklore/       # Book tracking (booklore.fig.systems)
    ├── FreshRSS/       # RSS reader (rss.fig.systems)
    ├── rsshub/         # RSS feed generator (rsshub.fig.systems)
    ├── microbin/       # Pastebin (paste.fig.systems)
    └── filebrowser/    # File manager (files.fig.systems)
```

## 🌐 Domains

Three domains are used with different purposes:

### fig.systems (Homelab Services)
Primary domain for all self-hosted homelab services:
- `*.fig.systems` - All services listed below

### edfig.dev (Professional/Public)
Professional and public-facing sites:
- `edfig.dev` / `www.edfig.dev` - Personal website/portfolio
- `blog.edfig.dev` - Technical blog

### figgy.foo (Experimental/Private)
Testing and experimental services:
- `figgy.foo` - Experimental lab (SSO protected)
- `*.figgy.foo` - Test instances of services

### Service URLs

| Service | URL | SSO Protected |
|---------|-----|---------------|
| Traefik Dashboard | traefik.fig.systems | ✅ |
| LLDAP | lldap.fig.systems | ✅ |
| Tinyauth | auth.fig.systems | ❌ |
| **Static Sites** | | |
| Personal Site | edfig.dev | ❌ |
| Blog | blog.edfig.dev | ❌ |
| Experimental Lab | figgy.foo | ✅ |
| **Monitoring** | | |
| Grafana (Logs) | logs.fig.systems | ❌* |
| Loki (API) | loki.fig.systems | ✅ |
| Uptime Kuma (Status) | status.fig.systems | ❌* |
| **Dashboard & Management** | | |
| Homarr | home.fig.systems | ✅ |
| Backrest | backup.fig.systems | ✅ |
| Jellyfin | flix.fig.systems | ❌* |
| Jellyseerr | requests.fig.systems | ✅ |
| Immich | photos.fig.systems | ❌* |
| Sonarr | sonarr.fig.systems | ✅ |
| Radarr | radarr.fig.systems | ✅ |
| SABnzbd | sabnzbd.fig.systems | ✅ |
| qBittorrent | qbt.fig.systems | ✅ |
| Profilarr | profilarr.fig.systems | ✅ |
| Karakeep | links.fig.systems | ✅ |
| Ollama (API) | ollama.fig.systems | ✅ |
| Vikunja | tasks.fig.systems | ✅ |
| LubeLogger | garage.fig.systems | ✅ |
| Calibre-web | books.fig.systems | ✅ |
| Booklore | booklore.fig.systems | ✅ |
| FreshRSS | rss.fig.systems | ✅ |
| RSSHub | rsshub.fig.systems | ❌* |
| MicroBin | paste.fig.systems | ❌* |
| File Browser | files.fig.systems | ✅ |

*Services marked with ❌* have their own authentication systems

## 📦 Media Folder Structure

The VM should have `/media` mounted at the root with this structure:

```
/media/
├── audiobooks/
├── books/
├── comics/
├── complete/      # Completed downloads
├── downloads/     # Active downloads
├── homemovies/
├── incomplete/    # Incomplete downloads
├── movies/
├── music/
├── photos/
└── tv/
```

## 🚀 Deployment

### Prerequisites

1. **DNS Configuration**: Point `*.fig.systems` and `*.edfig.dev` to your server IP
2. **Media Folders**: Ensure `/media` is mounted with the folder structure above
3. **Docker Network**: Create the homelab network

```bash
docker network create homelab
```

### Deployment Order

1. **Core Infrastructure** (must be first):
```bash
cd compose/core/traefik && docker compose up -d
cd compose/core/lldap && docker compose up -d
cd compose/core/tinyauth && docker compose up -d
```

2. **Configure LLDAP**:
   - Visit https://lldap.fig.systems
   - Login with admin credentials from `.env`
   - Create an observer user for tinyauth
   - Add regular users for authentication

3. **Update Passwords**:
   - Update `LLDAP_LDAP_USER_PASS` in `core/lldap/.env`
   - Update `LDAP_BIND_PASSWORD` in `core/tinyauth/.env` to match
   - Update `SESSION_SECRET` in `core/tinyauth/.env`
   - Update database passwords in service `.env` files

4. **Deploy Services**:
```bash
# Media frontend
cd compose/media/frontend/jellyfin && docker compose up -d
cd compose/media/frontend/jellyseer && docker compose up -d
cd compose/media/frontend/immich && docker compose up -d

# Media automation
cd compose/media/automation/sonarr && docker compose up -d
cd compose/media/automation/radarr && docker compose up -d
cd compose/media/automation/sabnzbd && docker compose up -d
cd compose/media/automation/qbittorrent && docker compose up -d

# Quality management (optional but recommended)
cd compose/media/automation/recyclarr && docker compose up -d
cd compose/media/automation/profilarr && docker compose up -d

# Utility services
cd compose/services/static-sites && docker compose up -d  # Static websites (edfig.dev, blog, figgy.foo)
cd compose/services/karakeep && docker compose up -d
cd compose/services/ollama && docker compose up -d
cd compose/services/vikunja && docker compose up -d
cd compose/services/homarr && docker compose up -d
cd compose/services/backrest && docker compose up -d

# Monitoring (optional but recommended)
cd compose/monitoring/logging && docker compose up -d
cd compose/monitoring/uptime && docker compose up -d
cd compose/services/lubelogger && docker compose up -d
cd compose/services/calibre-web && docker compose up -d
cd compose/services/booklore && docker compose up -d
cd compose/services/FreshRSS && docker compose up -d
cd compose/services/rsshub && docker compose up -d
cd compose/services/microbin && docker compose up -d
cd compose/services/filebrowser && docker compose up -d
```

## 🔐 Security Considerations

1. **Change Default Passwords**: All `.env` files contain placeholder passwords marked with `changeme_*`
2. **LLDAP Observer User**: Create a readonly user in LLDAP for tinyauth to bind
3. **SSL Certificates**: Traefik automatically obtains Let's Encrypt certificates
4. **Network Isolation**: Services use internal networks for database/cache communication
5. **SSO**: Most services are protected by tinyauth forward authentication

## 📝 Configuration Files

Each service has its own `.env` file where applicable. Key files to review:

- `core/lldap/.env` - LDAP configuration and admin credentials
- `core/tinyauth/.env` - LDAP connection and session settings
- `media/frontend/immich/.env` - Photo management configuration
- `services/karakeep/.env` - AI-powered bookmark manager
- `services/ollama/.env` - Local LLM configuration
- `services/microbin/.env` - Pastebin configuration

**Example Configuration Files:**
Several services include `.example` config files for reference:
- `media/automation/sonarr/config.xml.example`
- `media/automation/radarr/config.xml.example`
- `media/automation/sabnzbd/sabnzbd.ini.example`
- `media/automation/qbittorrent/qBittorrent.conf.example`
- `services/vikunja/config.yml.example`
- `services/FreshRSS/config.php.example`

Copy these to the appropriate location (usually `./config/`) and customize as needed.

## 🔧 Maintenance

### Viewing Logs
```bash
cd compose/[category]/[service]
docker compose logs -f
```

### Updating Services
```bash
cd compose/[category]/[service]
docker compose pull
docker compose up -d
```

### Backing Up Data
Important data locations:
- LLDAP: `compose/core/lldap/data/`
- Service configs: `compose/*/*/config/`
- Databases: `compose/*/*/db/` or `compose/*/*/pgdata/`
- Media: `/media/` (handle separately)

## 🐛 Troubleshooting

### Service won't start
1. Check logs: `docker compose logs`
2. Verify network exists: `docker network ls | grep homelab`
3. Check port conflicts: `docker ps -a`

### SSL certificate issues
1. Verify DNS points to your server
2. Check Traefik logs: `cd compose/core/traefik && docker compose logs`
3. Ensure ports 80 and 443 are open

### SSO not working
1. Verify tinyauth is running: `docker ps | grep tinyauth`
2. Check LLDAP connection in tinyauth logs
3. Verify LDAP bind credentials match in both services

### GPU not detected
1. Check GPU passthrough: `lspci | grep -i nvidia`
2. Verify drivers: `nvidia-smi`
3. Test in container: `docker exec ollama nvidia-smi`
4. See [AlmaLinux VM Setup](docs/setup/almalinux-vm.md) for GPU configuration

## 📊 Monitoring & Logging

### Centralized Logging (Loki + Promtail + Grafana)

All container logs are automatically collected and stored in Loki:

**Access Grafana**: https://logs.fig.systems

**Query examples:**
```logql
# View logs for specific service
{container="sonarr"}

# Filter by log level
{container="radarr"} |= "ERROR"

# Multiple services
{container=~"sonarr|radarr"}

# Search with JSON parsing
{container="karakeep"} |= "ollama" | json
```

**Retention**: 30 days (configurable in `compose/monitoring/logging/loki-config.yaml`)

### Uptime Monitoring (Uptime Kuma)

Monitor service availability and performance:

**Access Uptime Kuma**: https://status.fig.systems

**Features:**
- HTTP(s) monitoring for all web services
- Docker container health checks
- SSL certificate expiration alerts
- Public/private status pages
- 90+ notification integrations (Discord, Slack, Email, etc.)

### Service Integration

**How services integrate:**

```
Traefik (Reverse Proxy)
  ├─→ All services (SSL + routing)
  └─→ Let's Encrypt (certificates)

Tinyauth (SSO)
  ├─→ LLDAP (user authentication)
  └─→ Protected services (authorization)

Promtail (Log Collection)
  ├─→ Docker socket (all containers)
  └─→ Loki (log storage)

Loki (Log Storage)
  └─→ Grafana (visualization)

Karakeep (Bookmarks)
  ├─→ Ollama (AI tagging)
  ├─→ Meilisearch (search)
  └─→ Chrome (web archiving)

Sonarr/Radarr (Media Automation)
  ├─→ SABnzbd/qBittorrent (downloads)
  ├─→ Jellyfin (media library)
  └─→ Recyclarr/Profilarr (quality management)
```

See [Architecture Guide](docs/architecture.md) for complete integration details.

## 📄 License

This is a personal homelab configuration. Use at your own risk.

## 🤝 Contributing

This is a personal repository, but feel free to use it as a reference for your own homelab!
