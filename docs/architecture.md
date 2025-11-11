# Homelab Architecture & Integration

Complete integration guide for the homelab setup on AlmaLinux 9.6.

## 🖥️ Hardware Specifications

### Host System
- **Hypervisor**: Proxmox VE 9 (Debian 13 based)
- **CPU**: AMD Ryzen 5 7600X (6 cores, 12 threads, up to 5.3 GHz)
- **GPU**: NVIDIA GeForce GTX 1070 (8GB VRAM, 1920 CUDA cores)
- **RAM**: 32GB DDR5

### VM Configuration
- **OS**: AlmaLinux 9.6 (RHEL 9 compatible)
- **CPU**: 8 vCPUs (allocated from host)
- **RAM**: 24GB (leaving 8GB for host)
- **Storage**: 500GB+ (adjust based on media library size)
- **GPU**: GTX 1070 (PCIe passthrough from Proxmox)

## 🏗️ Architecture Overview

### Network Architecture

```
Internet
    ↓
[Router/Firewall]
    ↓ (Port 80/443)
[Traefik Reverse Proxy]
    ↓
┌──────────────────────────────────────┐
│         homelab network              │
│  (Docker bridge - 172.18.0.0/16)    │
│                                      │
│  ┌─────────────┐  ┌──────────────┐  │
│  │ Core        │  │ Media        │  │
│  │ - Traefik   │  │ - Jellyfin   │  │
│  │ - LLDAP     │  │ - Sonarr     │  │
│  │ - Tinyauth  │  │ - Radarr     │  │
│  └─────────────┘  └──────────────┘  │
│                                      │
│  ┌─────────────┐  ┌──────────────┐  │
│  │ Services    │  │ Monitoring   │  │
│  │ - Karakeep  │  │ - Loki       │  │
│  │ - Ollama    │  │ - Promtail   │  │
│  │ - Vikunja   │  │ - Grafana    │  │
│  └─────────────┘  └──────────────┘  │
└──────────────────────────────────────┘
           ↓
    [Promtail Agent]
           ↓
      [Loki Storage]
```

### Service Internal Networks

Services with databases use isolated internal networks:

```
karakeep
├── homelab (external traffic)
└── karakeep_internal
    ├── karakeep (app)
    ├── karakeep-chrome (browser)
    └── karakeep-meilisearch (search)

vikunja
├── homelab (external traffic)
└── vikunja_internal
    ├── vikunja (app)
    └── vikunja-db (postgres)

monitoring/logging
├── homelab (external traffic)
└── logging_internal
    ├── loki (storage)
    ├── promtail (collector)
    └── grafana (UI)
```

## 🔐 Security Architecture

### Authentication Flow

```
User Request
    ↓
[Traefik] → Check route rules
    ↓
[Tinyauth Middleware] → Forward Auth
    ↓
[LLDAP] → Verify credentials
    ↓
[Backend Service] → Authorized access
```

### SSL/TLS

- **Certificate Provider**: Let's Encrypt
- **Challenge Type**: HTTP-01 (ports 80/443)
- **Automatic Renewal**: Via Traefik
- **Domains**:
  - Primary: `*.fig.systems`
  - Fallback: `*.edfig.dev`

### SSO Protection

**Protected Services** (require authentication):
- Traefik Dashboard
- LLDAP
- Sonarr, Radarr, SABnzbd, qBittorrent
- Profilarr, Recyclarr (monitoring)
- Homarr, Backrest
- Karakeep, Vikunja, LubeLogger
- Calibre-web, Booklore, FreshRSS, File Browser
- Loki API, Ollama API

**Unprotected Services** (own authentication):
- Tinyauth (SSO provider itself)
- Jellyfin (own user system)
- Jellyseerr (linked to Jellyfin)
- Immich (own user system)
- RSSHub (public feed generator)
- MicroBin (public pastebin)
- Grafana (own authentication)
- Uptime Kuma (own authentication)

## 📊 Logging Architecture

### Centralized Logging with Loki

All services forward logs to Loki via Promtail:

```
[Docker Container] → stdout/stderr
         ↓
[Docker Socket] → /var/run/docker.sock
         ↓
[Promtail] → Scrapes logs via Docker API
         ↓
[Loki] → Stores and indexes logs
         ↓
[Grafana] → Query and visualize
```

### Log Labels

Promtail automatically adds labels to all logs:
- `container`: Container name
- `compose_project`: Docker Compose project
- `compose_service`: Service name from compose
- `image`: Docker image name
- `stream`: stdout or stderr

### Log Retention

- **Default**: 30 days
- **Storage**: `compose/monitoring/logging/loki-data/`
- **Automatic cleanup**: Enabled via Loki compactor

### Querying Logs

**View all logs for a service:**
```logql
{container="sonarr"}
```

**Filter by log level:**
```logql
{container="radarr"} |= "ERROR"
```

**Multiple services:**
```logql
{container=~"sonarr|radarr"}
```

**Time range with filters:**
```logql
{container="karakeep"} |= "ollama" | json
```

## 🌐 Network Configuration

### Docker Networks

**homelab** (external bridge):
- Type: External bridge network
- Subnet: Auto-assigned by Docker
- Purpose: Inter-service communication + Traefik routing
- Create: `docker network create homelab`

**Service-specific internal networks**:
- `karakeep_internal`: Karakeep + Chrome + Meilisearch
- `vikunja_internal`: Vikunja + PostgreSQL
- `logging_internal`: Loki + Promtail + Grafana
- etc.

### Port Mappings

**External Ports** (exposed to host):
- `80/tcp`: HTTP (Traefik) - redirects to HTTPS
- `443/tcp`: HTTPS (Traefik)
- `6881/tcp+udp`: BitTorrent (qBittorrent)

**No other ports exposed** - all access via Traefik reverse proxy.

## 🔧 Traefik Integration

### Standard Traefik Labels

All services use consistent Traefik labels:

```yaml
labels:
  # Enable Traefik
  traefik.enable: true
  traefik.docker.network: homelab

  # Router configuration
  traefik.http.routers.<service>.rule: Host(`<service>.fig.systems`) || Host(`<service>.edfig.dev`)
  traefik.http.routers.<service>.entrypoints: websecure
  traefik.http.routers.<service>.tls.certresolver: letsencrypt

  # Service configuration (backend port)
  traefik.http.services.<service>.loadbalancer.server.port: <port>

  # SSO middleware (if protected)
  traefik.http.routers.<service>.middlewares: tinyauth

  # Homarr auto-discovery
  homarr.name: <Service Name>
  homarr.group: <Category>
  homarr.icon: mdi:<icon-name>
```

### Middleware

**tinyauth** - Forward authentication:
```yaml
# Defined in traefik/compose.yaml
middlewares:
  tinyauth:
    forwardAuth:
      address: http://tinyauth:8080
      trustForwardHeader: true
```

## 💾 Volume Management

### Volume Types

**Bind Mounts** (host directories):
```yaml
volumes:
  - ./data:/data          # Service data
  - ./config:/config      # Configuration files
  - /media:/media         # Media library (shared)
```

**Named Volumes** (Docker-managed):
```yaml
volumes:
  - loki-data:/loki       # Loki storage
  - postgres-data:/var/lib/postgresql/data
```

### Media Directory Structure

```
/media/
├── tv/              # TV shows (Sonarr → Jellyfin)
├── movies/          # Movies (Radarr → Jellyfin)
├── music/           # Music
├── photos/          # Photos (Immich)
├── books/           # Ebooks (Calibre-web)
├── audiobooks/      # Audiobooks
├── comics/          # Comics
├── homemovies/      # Home videos
├── downloads/       # Active downloads (SABnzbd/qBittorrent)
├── complete/        # Completed downloads
└── incomplete/      # In-progress downloads
```

### Backup Strategy

**Important directories to backup:**
```
compose/core/lldap/data/              # User directory
compose/core/traefik/letsencrypt/     # SSL certificates
compose/services/*/config/            # Service configurations
compose/services/*/data/              # Service data
compose/monitoring/logging/loki-data/ # Logs (optional)
/media/                               # Media library
```

**Excluded from backups:**
```
compose/services/*/db/                # Databases (backup via dump)
compose/monitoring/logging/loki-data/ # Logs (can be recreated)
/media/downloads/                     # Temporary downloads
/media/incomplete/                    # Incomplete downloads
```

## 🎮 GPU Acceleration

### NVIDIA GTX 1070 Configuration

**GPU Passthrough (Proxmox → VM):**

1. **Proxmox host** (`/etc/pve/nodes/<node>/qemu-server/<vmid>.conf`):
```
hostpci0: 0000:01:00,pcie=1,x-vga=1
```

2. **VM (AlmaLinux)** - Install NVIDIA drivers:
```bash
# Add NVIDIA repository
sudo dnf config-manager --add-repo https://developer.download.nvidia.com/compute/cuda/repos/rhel9/x86_64/cuda-rhel9.repo

# Install drivers
sudo dnf install nvidia-driver nvidia-settings

# Verify
nvidia-smi
```

3. **Docker** - Install NVIDIA Container Toolkit:
```bash
# Add NVIDIA Container Toolkit repo
sudo dnf config-manager --add-repo https://nvidia.github.io/libnvidia-container/stable/rpm/nvidia-container-toolkit.repo

# Install toolkit
sudo dnf install nvidia-container-toolkit

# Configure Docker
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker

# Verify
docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi
```

### Services Using GPU

**Jellyfin** (Hardware transcoding):
```yaml
# Uncomment in compose.yaml
devices:
  - /dev/dri:/dev/dri  # For NVENC/NVDEC
environment:
  - NVIDIA_VISIBLE_DEVICES=all
  - NVIDIA_DRIVER_CAPABILITIES=all
```

**Immich** (AI features):
```yaml
# Already configured
deploy:
  resources:
    reservations:
      devices:
        - driver: nvidia
          count: 1
          capabilities: [gpu]
```

**Ollama** (LLM inference):
```yaml
# Uncomment in compose.yaml
deploy:
  resources:
    reservations:
      devices:
        - driver: nvidia
          count: 1
          capabilities: [gpu]
```

### GPU Performance Tuning

**For Ryzen 5 7600X + GTX 1070:**

- **Jellyfin**: Can transcode 4-6 simultaneous 4K → 1080p streams
- **Ollama**:
  - 3B models: 40-60 tokens/sec
  - 7B models: 20-35 tokens/sec
  - 13B models: 10-15 tokens/sec (quantized)
- **Immich**: AI tagging ~5-10 images/sec

## 🚀 Resource Allocation

### CPU Allocation (Ryzen 5 7600X - 6C/12T)

**High Priority** (4-6 cores):
- Jellyfin (transcoding)
- Sonarr/Radarr (media processing)
- Ollama (when running)

**Medium Priority** (2-4 cores):
- Immich (AI processing)
- Karakeep (bookmark processing)
- SABnzbd/qBittorrent (downloads)

**Low Priority** (1-2 cores):
- Traefik, LLDAP, Tinyauth
- Monitoring services
- Other utilities

### RAM Allocation (32GB Total, 24GB VM)

**Recommended allocation:**

```
Host (Proxmox): 8GB
VM Total: 24GB breakdown:
  ├── System: 4GB (AlmaLinux base)
  ├── Docker: 2GB (daemon overhead)
  ├── Jellyfin: 2-4GB (transcoding buffers)
  ├── Immich: 2-3GB (ML models + database)
  ├── Sonarr/Radarr: 1GB each
  ├── Ollama: 4-6GB (when running models)
  ├── Databases: 2-3GB total
  ├── Monitoring: 2GB (Loki + Grafana)
  └── Other services: 4-5GB
```

### Disk Space Planning

**System:** 100GB
**Docker:** 50GB (images + containers)
**Service Data:** 50GB (configs, databases, logs)
**Media Library:** Remaining space (expandable)

**Recommended VM disk:**
- Minimum: 500GB (200GB system + 300GB media)
- Recommended: 1TB+ (allows room for growth)

## 🔄 Service Dependencies

### Startup Order

**Critical order for initial deployment:**

1. **Networks**: `docker network create homelab`
2. **Core** (must start first):
   - Traefik (reverse proxy)
   - LLDAP (user directory)
   - Tinyauth (SSO provider)
3. **Monitoring** (optional but recommended):
   - Loki + Promtail + Grafana
   - Uptime Kuma
4. **Media Automation**:
   - Sonarr, Radarr
   - SABnzbd, qBittorrent
   - Recyclarr, Profilarr
5. **Media Frontend**:
   - Jellyfin
   - Jellyseer
   - Immich
6. **Services**:
   - Karakeep, Ollama (AI features)
   - Vikunja, Homarr
   - All other services

### Service Integration Map

```
Traefik
 ├─→ All services (reverse proxy)
 └─→ Let's Encrypt (SSL)

Tinyauth
 ├─→ LLDAP (authentication backend)
 └─→ All SSO-protected services

LLDAP
 └─→ User database for SSO

Promtail
 ├─→ Docker socket (log collection)
 └─→ Loki (log forwarding)

Loki
 └─→ Grafana (log visualization)

Karakeep
 ├─→ Ollama (AI tagging)
 ├─→ Meilisearch (search)
 └─→ Chrome (web archiving)

Jellyseer
 ├─→ Jellyfin (media info)
 ├─→ Sonarr (TV requests)
 └─→ Radarr (movie requests)

Sonarr/Radarr
 ├─→ SABnzbd/qBittorrent (downloads)
 ├─→ Jellyfin (media library)
 └─→ Recyclarr/Profilarr (quality profiles)

Homarr
 └─→ All services (dashboard auto-discovery)
```

## 🐛 Troubleshooting

### Check Service Health

```bash
# All services status
cd ~/homelab
docker ps -a --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Logs for specific service
docker logs <service-name> --tail 100 -f

# Logs via Loki/Grafana
# Go to https://logs.fig.systems
# Query: {container="<service-name>"}
```

### Network Issues

```bash
# Check homelab network exists
docker network ls | grep homelab

# Inspect network
docker network inspect homelab

# Test service connectivity
docker exec <service-a> ping <service-b>
docker exec karakeep curl http://ollama:11434
```

### GPU Not Detected

```bash
# Check GPU in VM
nvidia-smi

# Check Docker can access GPU
docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi

# Check service GPU allocation
docker exec jellyfin nvidia-smi
docker exec ollama nvidia-smi
```

### SSL Certificate Issues

```bash
# Check Traefik logs
docker logs traefik | grep -i certificate

# Force certificate renewal
docker exec traefik rm -rf /letsencrypt/acme.json
docker restart traefik

# Verify DNS
dig +short sonarr.fig.systems
```

### SSO Not Working

```bash
# Check Tinyauth status
docker logs tinyauth

# Check LLDAP connection
docker exec tinyauth nc -zv lldap 3890
docker exec tinyauth nc -zv lldap 17170

# Verify credentials match
grep LDAP_BIND_PASSWORD compose/core/tinyauth/.env
grep LLDAP_LDAP_USER_PASS compose/core/lldap/.env
```

## 📈 Monitoring Best Practices

### Key Metrics to Monitor

**System Level:**
- CPU usage per container
- Memory usage per container
- Disk I/O
- Network throughput
- GPU utilization (for Jellyfin/Ollama/Immich)

**Application Level:**
- Traefik request rate
- Failed authentication attempts
- Jellyfin concurrent streams
- Download speeds (SABnzbd/qBittorrent)
- Sonarr/Radarr queue size

### Uptime Kuma Monitoring

Configure monitors for:
- **HTTP(s)**: All web services (200 status check)
- **TCP**: Database ports (PostgreSQL, etc.)
- **Docker**: Container health (via Docker socket)
- **SSL**: Certificate expiration (30-day warning)

### Log Monitoring

Set up Loki alerts for:
- ERROR level logs
- Authentication failures
- Service crashes
- Disk space warnings

## 🔧 Maintenance Tasks

### Daily
- Check Uptime Kuma dashboard
- Review any critical alerts

### Weekly
- Check disk space: `df -h`
- Review failed downloads in Sonarr/Radarr
- Check Loki logs for errors

### Monthly
- Update all containers: `docker compose pull && docker compose up -d`
- Review and clean old Docker images: `docker image prune -a`
- Backup configurations
- Check SSL certificate renewal

### Quarterly
- Review and update documentation
- Clean up old media (if needed)
- Review and adjust quality profiles
- Update Recyclarr configurations

## 📚 Additional Resources

- [Traefik Documentation](https://doc.traefik.io/traefik/)
- [Docker Compose Best Practices](https://docs.docker.com/compose/production/)
- [Loki LogQL Guide](https://grafana.com/docs/loki/latest/logql/)
- [NVIDIA Container Toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/)
- [Proxmox GPU Passthrough](https://pve.proxmox.com/wiki/PCI_Passthrough)
- [AlmaLinux Documentation](https://wiki.almalinux.org/)

---

**System Ready!** 🚀
