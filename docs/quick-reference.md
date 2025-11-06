# Quick Reference Guide

Fast reference for common tasks and commands.

## Service URLs

All services accessible via:
- Primary domain: `*.fig.systems`
- Secondary domain: `*.edfig.dev`

### Core Services
```
https://traefik.fig.systems  # Reverse proxy dashboard
https://lldap.fig.systems    # User directory
https://auth.fig.systems     # SSO authentication
```

### Dashboard & Management
```
https://home.fig.systems     # Homarr dashboard (START HERE!)
https://backup.fig.systems   # Backrest backup manager
```

### Media Services
```
https://flix.fig.systems      # Jellyfin media server
https://photos.fig.systems    # Immich photo library
https://requests.fig.systems  # Jellyseerr media requests
https://sonarr.fig.systems    # TV show automation
https://radarr.fig.systems    # Movie automation
https://sabnzbd.fig.systems   # Usenet downloader
https://qbt.fig.systems       # qBittorrent client
```

### Utility Services
```
https://links.fig.systems     # Linkwarden bookmarks
https://tasks.fig.systems     # Vikunja task management
https://garage.fig.systems    # LubeLogger vehicle tracking
https://books.fig.systems     # Calibre-web ebook library
https://booklore.fig.systems  # Book tracking
https://rss.fig.systems       # FreshRSS reader
https://files.fig.systems     # File Browser
```

## Common Commands

### Docker Compose

```bash
# Start service
cd ~/homelab/compose/path/to/service
docker compose up -d

# View logs
docker compose logs -f

# Restart service
docker compose restart

# Stop service
docker compose down

# Update and restart
docker compose pull
docker compose up -d

# Rebuild service
docker compose up -d --force-recreate
```

### Docker Management

```bash
# List all containers
docker ps

# List all containers (including stopped)
docker ps -a

# View logs
docker logs <container_name>
docker logs -f <container_name>  # Follow logs

# Execute command in container
docker exec -it <container_name> bash

# View resource usage
docker stats

# Remove stopped containers
docker container prune

# Remove unused images
docker image prune -a

# Remove unused volumes (CAREFUL!)
docker volume prune

# Complete cleanup
docker system prune -a --volumes
```

### Service Management

```bash
# Start all core services
cd ~/homelab/compose/core
for dir in traefik lldap tinyauth; do
  cd $dir && docker compose up -d && cd ..
done

# Stop all services
cd ~/homelab
find compose -name "compose.yaml" -execdir docker compose down \;

# Restart single service
cd ~/homelab/compose/services/servicename
docker compose restart

# View all running containers
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

### System Checks

```bash
# Check all containers
docker ps --format "table {{.Names}}\t{{.Status}}"

# Check network
docker network inspect homelab

# Check disk usage
docker system df
df -h

# Check logs for errors
docker compose logs --tail=100 | grep -i error

# Test DNS resolution
dig home.fig.systems +short

# Test SSL
curl -I https://home.fig.systems
```

## Secret Generation

```bash
# JWT/Session secrets (64 char)
openssl rand -hex 32

# Database passwords (32 char alphanumeric)
openssl rand -base64 32 | tr -d /=+ | cut -c1-32

# API keys (32 char hex)
openssl rand -hex 16

# Find what needs updating
grep -r "changeme_" ~/homelab/compose
```

## Troubleshooting

### Service Won't Start

```bash
# Check logs
docker compose logs

# Check container status
docker compose ps

# Check for port conflicts
sudo netstat -tulpn | grep :80
sudo netstat -tulpn | grep :443

# Recreate container
docker compose down
docker compose up -d
```

### SSL Certificate Issues

```bash
# Check Traefik logs
docker logs traefik | grep -i certificate

# Check Let's Encrypt logs
docker logs traefik | grep -i letsencrypt

# Verify DNS
dig home.fig.systems +short

# Test port 80 accessibility
curl -I http://home.fig.systems
```

### SSO Not Working

```bash
# Check LLDAP
docker logs lldap

# Check Tinyauth
docker logs tinyauth

# Verify passwords match
grep LLDAP_LDAP_USER_PASS ~/homelab/compose/core/lldap/.env
grep LDAP_BIND_PASSWORD ~/homelab/compose/core/tinyauth/.env

# Test LDAP connection
docker exec tinyauth nc -zv lldap 3890
```

### Database Connection Failures

```bash
# Check database container
docker ps | grep postgres

# View database logs
docker logs <db_container_name>

# Test connection from app container
docker exec <app_container> nc -zv <db_container> 5432

# Verify password in .env
cat .env | grep POSTGRES_PASSWORD
```

## File Locations

### Configuration
```
~/homelab/compose/               # All services
~/homelab/compose/core/          # Core infrastructure
~/homelab/compose/media/         # Media services
~/homelab/compose/services/      # Utility services
```

### Service Data
```
compose/<service>/config/        # Service configuration
compose/<service>/data/          # Service data
compose/<service>/db/            # Database files
compose/<service>/.env           # Environment variables
```

### Media Files
```
/media/movies/                   # Movies
/media/tv/                       # TV shows
/media/music/                    # Music
/media/photos/                   # Photos
/media/books/                    # Books
/media/downloads/                # Active downloads
/media/complete/                 # Completed downloads
```

### Logs
```
docker logs <container_name>     # Container logs
compose/<service>/logs/          # Service-specific logs (if configured)
/var/lib/docker/volumes/         # Volume data
```

## Network

### Create Network
```bash
docker network create homelab
```

### Inspect Network
```bash
docker network inspect homelab
```

### Connect Container to Network
```bash
docker network connect homelab <container_name>
```

## GPU (NVIDIA GTX 1070)

### Check GPU Status
```bash
nvidia-smi
```

### Test GPU in Docker
```bash
docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi
```

### Monitor GPU Usage
```bash
watch -n 1 nvidia-smi
```

### Check GPU in Container
```bash
docker exec jellyfin nvidia-smi
docker exec immich_machine_learning nvidia-smi
```

## Backup

### Backup Configuration Files
```bash
cd ~/homelab
tar czf homelab-config-$(date +%Y%m%d).tar.gz \
  $(find compose -name ".env") \
  $(find compose -name "compose.yaml")
```

### Backup Service Data
```bash
# Example: Backup Immich
cd ~/homelab/compose/media/frontend/immich
tar czf immich-backup-$(date +%Y%m%d).tar.gz upload/ config/
```

### Restore Configuration
```bash
tar xzf homelab-config-YYYYMMDD.tar.gz
```

## Updates

### Update Single Service
```bash
cd ~/homelab/compose/path/to/service
docker compose pull
docker compose up -d
```

### Update All Services
```bash
cd ~/homelab
for dir in $(find compose -name "compose.yaml" -exec dirname {} \;); do
  echo "Updating $dir"
  cd $dir
  docker compose pull
  docker compose up -d
  cd ~/homelab
done
```

### Update Docker
```bash
sudo apt update
sudo apt upgrade docker-ce docker-ce-cli containerd.io
```

## Performance

### Check Resource Usage
```bash
# Overall system
htop

# Docker containers
docker stats

# Disk usage
df -h
docker system df

# Network usage
iftop
```

### Clean Up Disk Space
```bash
# Docker cleanup
docker system prune -a

# Remove old logs
sudo journalctl --vacuum-time=7d

# Find large files
du -h /media | sort -rh | head -20
```

## DNS Configuration

### Cloudflare Example
```
Type: A
Name: *
Content: YOUR_SERVER_IP
Proxy: Off (disable for Let's Encrypt)
TTL: Auto
```

### Local DNS (Pi-hole/hosts file)
```
192.168.1.100  home.fig.systems
192.168.1.100  flix.fig.systems
192.168.1.100  photos.fig.systems
# ... etc
```

## Environment Variables

### List All Services with Secrets
```bash
find ~/homelab/compose -name ".env" -exec echo {} \;
```

### Check for Unconfigured Secrets
```bash
grep -r "changeme_" ~/homelab/compose | wc -l
# Should be 0
```

### Backup All .env Files
```bash
cd ~/homelab
tar czf env-files-$(date +%Y%m%d).tar.gz $(find compose -name ".env")
gpg -c env-files-$(date +%Y%m%d).tar.gz
```

## Monitoring

### Service Health
```bash
# Check all containers are running
docker ps --format "{{.Names}}: {{.Status}}" | grep -v "Up"

# Check for restarts
docker ps --format "{{.Names}}: {{.Status}}" | grep "Restarting"

# Check logs for errors
docker compose logs --tail=100 | grep -i error
```

### SSL Certificate Expiry
```bash
# Check cert expiry
echo | openssl s_client -servername home.fig.systems -connect home.fig.systems:443 2>/dev/null | openssl x509 -noout -dates
```

### Disk Space
```bash
# Overall
df -h

# Docker
docker system df

# Media
du -sh /media/*
```

## Common File Paths

```bash
# Core services
~/homelab/compose/core/traefik/
~/homelab/compose/core/lldap/
~/homelab/compose/core/tinyauth/

# Media
~/homelab/compose/media/frontend/jellyfin/
~/homelab/compose/media/frontend/immich/
~/homelab/compose/media/automation/sonarr/

# Utilities
~/homelab/compose/services/homarr/
~/homelab/compose/services/backrest/
~/homelab/compose/services/linkwarden/

# Documentation
~/homelab/docs/
~/homelab/README.md
```

## Port Reference

```
80   - HTTP (Traefik)
443  - HTTPS (Traefik)
3890 - LLDAP
6881 - qBittorrent (TCP/UDP)
8096 - Jellyfin
2283 - Immich
```

## Default Credentials

⚠️ **Change these immediately after first login!**

### qBittorrent
```
Username: admin
Password: adminadmin
```

### Microbin
```
Check compose/services/microbin/.env
MICROBIN_ADMIN_USERNAME
MICROBIN_ADMIN_PASSWORD
```

### All Other Services
Use SSO (LLDAP) or create admin account on first visit.

## Quick Deployment

### Deploy Everything
```bash
cd ~/homelab
chmod +x deploy-all.sh
./deploy-all.sh
```

### Deploy Core Only
```bash
cd ~/homelab/compose/core/traefik && docker compose up -d
cd ../lldap && docker compose up -d
cd ../tinyauth && docker compose up -d
```

### Deploy Media Stack
```bash
cd ~/homelab/compose/media/frontend
for dir in */; do cd "$dir" && docker compose up -d && cd ..; done

cd ~/homelab/compose/media/automation
for dir in */; do cd "$dir" && docker compose up -d && cd ..; done
```

## Emergency Procedures

### Stop All Services
```bash
cd ~/homelab
find compose -name "compose.yaml" -execdir docker compose down \;
```

### Remove All Containers (Nuclear Option)
```bash
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)
```

### Reset Network
```bash
docker network rm homelab
docker network create homelab
```

### Reset Service
```bash
cd ~/homelab/compose/path/to/service
docker compose down -v  # REMOVES VOLUMES!
docker compose up -d
```

---

**For detailed guides, see the [docs folder](./README.md).**
