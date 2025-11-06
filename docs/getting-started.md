# Getting Started with Homelab

This guide will walk you through setting up your homelab from scratch.

## Prerequisites

### Hardware Requirements
- **Server/VM**: Linux server with Docker support
- **CPU**: 2+ cores recommended
- **RAM**: 8GB minimum, 16GB+ recommended
- **Storage**: 100GB+ for Docker containers and config
- **Optional GPU**: NVIDIA GPU for hardware transcoding (Jellyfin, Immich)

### Software Requirements
- **Operating System**: Ubuntu 22.04 or similar Linux distribution
- **Docker**: Version 24.0+
- **Docker Compose**: Version 2.20+
- **Git**: For cloning the repository
- **Domain Names**: `*.fig.systems` and `*.edfig.dev` (or your domains)

### Network Requirements
- **Ports**: 80 and 443 accessible from internet (for Let's Encrypt)
- **DNS**: Ability to create A records for your domains
- **Static IP**: Recommended for your homelab server

## Step 1: Prepare Your Server

### Install Docker and Docker Compose

```bash
# Update package index
sudo apt update

# Install dependencies
sudo apt install -y ca-certificates curl gnupg lsb-release

# Add Docker's official GPG key
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Set up the repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Add your user to docker group (logout and login after this)
sudo usermod -aG docker $USER

# Verify installation
docker --version
docker compose version
```

### Create Media Directory Structure

```bash
# Create media folders
sudo mkdir -p /media/{audiobooks,books,comics,complete,downloads,homemovies,incomplete,movies,music,photos,tv}

# Set ownership (replace with your username)
sudo chown -R $(whoami):$(whoami) /media

# Verify structure
tree -L 1 /media
```

## Step 2: Clone the Repository

```bash
# Clone the repository
cd ~
git clone https://github.com/efigueroa/homelab.git
cd homelab

# Checkout the main branch
git checkout main  # or your target branch
```

## Step 3: Configure DNS

You need to point your domains to your server's IP address.

### Option 1: Wildcard DNS (Recommended)

Add these A records to your DNS provider:

```
*.fig.systems    A    YOUR_SERVER_IP
*.edfig.dev      A    YOUR_SERVER_IP
```

### Option 2: Individual Records

Create A records for each service:

```
traefik.fig.systems     A    YOUR_SERVER_IP
lldap.fig.systems       A    YOUR_SERVER_IP
auth.fig.systems        A    YOUR_SERVER_IP
home.fig.systems        A    YOUR_SERVER_IP
backup.fig.systems      A    YOUR_SERVER_IP
flix.fig.systems        A    YOUR_SERVER_IP
photos.fig.systems      A    YOUR_SERVER_IP
# ... and so on for all services
```

### Verify DNS

Wait a few minutes for DNS propagation, then verify:

```bash
# Test DNS resolution
dig traefik.fig.systems +short
dig lldap.fig.systems +short

# Should return your server IP
```

## Step 4: Configure Environment Variables

Each service needs its environment variables configured with secure values.

### Generate Secure Secrets

Use these commands to generate secure values:

```bash
# For JWT secrets and session secrets (64 characters)
openssl rand -hex 32

# For passwords (32 alphanumeric characters)
openssl rand -base64 32 | tr -d /=+ | cut -c1-32

# For API keys (32 characters)
openssl rand -hex 16
```

### Update Core Services

**LLDAP** (`compose/core/lldap/.env`):
```bash
cd compose/core/lldap
nano .env

# Update these values:
LLDAP_LDAP_USER_PASS=<your-strong-password>
LLDAP_JWT_SECRET=<output-from-openssl-rand-hex-32>
```

**Tinyauth** (`compose/core/tinyauth/.env`):
```bash
cd ../tinyauth
nano .env

# Update these values (LDAP_BIND_PASSWORD must match LLDAP_LDAP_USER_PASS):
LDAP_BIND_PASSWORD=<same-as-LLDAP_LDAP_USER_PASS>
SESSION_SECRET=<output-from-openssl-rand-hex-32>
```

**Immich** (`compose/media/frontend/immich/.env`):
```bash
cd ../../media/frontend/immich
nano .env

# Update:
DB_PASSWORD=<output-from-openssl-rand-base64>
```

### Update All Other Services

Go through each service's `.env` file and replace all `changeme_*` values:

```bash
# Find all files that need updating
grep -r "changeme_" ~/homelab/compose

# Or update them individually
cd ~/homelab/compose/services/linkwarden
nano .env  # Update NEXTAUTH_SECRET, POSTGRES_PASSWORD, MEILI_MASTER_KEY

cd ../vikunja
nano .env  # Update VIKUNJA_DATABASE_PASSWORD, VIKUNJA_SERVICE_JWTSECRET, POSTGRES_PASSWORD
```

💡 **Tip**: Keep your secrets in a password manager!

See [Secrets Management Guide](./guides/secrets-management.md) for detailed instructions.

## Step 5: Create Docker Network

```bash
# Create the external homelab network
docker network create homelab

# Verify it was created
docker network ls | grep homelab
```

## Step 6: Deploy Services

Deploy services in order, starting with core infrastructure:

### Deploy Core Infrastructure

```bash
cd ~/homelab

# Deploy Traefik (reverse proxy)
cd compose/core/traefik
docker compose up -d

# Check logs to ensure it starts successfully
docker compose logs -f

# Wait for "Server configuration reloaded" message, then Ctrl+C
```

```bash
# Deploy LLDAP (user directory)
cd ../lldap
docker compose up -d
docker compose logs -f

# Access: https://lldap.fig.systems
# Default login: admin / <your LLDAP_LDAP_USER_PASS>
```

```bash
# Deploy Tinyauth (SSO)
cd ../tinyauth
docker compose up -d
docker compose logs -f

# Access: https://auth.fig.systems
```

### Create LLDAP Users

Before deploying other services, create your user in LLDAP:

1. Go to https://lldap.fig.systems
2. Login with admin credentials
3. Create your user:
   - Username: `edfig` (or your choice)
   - Email: `admin@edfig.dev`
   - Password: strong password
   - Add to `lldap_admin` group

### Deploy Media Services

```bash
cd ~/homelab/compose/media/frontend

# Jellyfin
cd jellyfin
docker compose up -d
# Access: https://flix.fig.systems

# Immich
cd ../immich
docker compose up -d
# Access: https://photos.fig.systems

# Jellyseerr
cd ../jellyseer
docker compose up -d
# Access: https://requests.fig.systems
```

```bash
# Media automation
cd ~/homelab/compose/media/automation

cd sonarr && docker compose up -d && cd ..
cd radarr && docker compose up -d && cd ..
cd sabnzbd && docker compose up -d && cd ..
cd qbittorrent && docker compose up -d && cd ..
```

### Deploy Utility Services

```bash
cd ~/homelab/compose/services

# Dashboard (start with this - it shows all your services!)
cd homarr && docker compose up -d && cd ..
# Access: https://home.fig.systems

# Backup manager
cd backrest && docker compose up -d && cd ..
# Access: https://backup.fig.systems

# Other services
cd linkwarden && docker compose up -d && cd ..
cd vikunja && docker compose up -d && cd ..
cd lubelogger && docker compose up -d && cd ..
cd calibre-web && docker compose up -d && cd ..
cd booklore && docker compose up -d && cd ..
cd FreshRSS && docker compose up -d && cd ..
cd rsshub && docker compose up -d && cd ..
cd microbin && docker compose up -d && cd ..
cd filebrowser && docker compose up -d && cd ..
```

### Quick Deploy All (Alternative)

If you've configured everything and want to deploy all at once:

```bash
cd ~/homelab

# Create a deployment script
cat > deploy-all.sh << 'SCRIPT'
#!/bin/bash
set -e

echo "Deploying homelab services..."

# Core
echo "==> Core Infrastructure"
cd compose/core/traefik && docker compose up -d && cd ../../..
sleep 5
cd compose/core/lldap && docker compose up -d && cd ../../..
sleep 5
cd compose/core/tinyauth && docker compose up -d && cd ../../..

# Media
echo "==> Media Services"
cd compose/media/frontend/immich && docker compose up -d && cd ../../../..
cd compose/media/frontend/jellyfin && docker compose up -d && cd ../../../..
cd compose/media/frontend/jellyseer && docker compose up -d && cd ../../../..
cd compose/media/automation/sonarr && docker compose up -d && cd ../../../..
cd compose/media/automation/radarr && docker compose up -d && cd ../../../..
cd compose/media/automation/sabnzbd && docker compose up -d && cd ../../../..
cd compose/media/automation/qbittorrent && docker compose up -d && cd ../../../..

# Utility
echo "==> Utility Services"
cd compose/services/homarr && docker compose up -d && cd ../..
cd compose/services/backrest && docker compose up -d && cd ../..
cd compose/services/linkwarden && docker compose up -d && cd ../..
cd compose/services/vikunja && docker compose up -d && cd ../..
cd compose/services/lubelogger && docker compose up -d && cd ../..
cd compose/services/calibre-web && docker compose up -d && cd ../..
cd compose/services/booklore && docker compose up -d && cd ../..
cd compose/services/FreshRSS && docker compose up -d && cd ../..
cd compose/services/rsshub && docker compose up -d && cd ../..
cd compose/services/microbin && docker compose up -d && cd ../..
cd compose/services/filebrowser && docker compose up -d && cd ../..

echo "==> Deployment Complete!"
echo "Access your dashboard at: https://home.fig.systems"
SCRIPT

chmod +x deploy-all.sh
./deploy-all.sh
```

## Step 7: Verify Deployment

### Check All Containers Are Running

```bash
# List all containers
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Check for any stopped containers
docker ps -a --filter "status=exited"
```

### Verify SSL Certificates

```bash
# Test SSL certificate
curl -I https://home.fig.systems

# Should show HTTP/2 200 and valid SSL cert
```

### Access Services

Visit your dashboard: **https://home.fig.systems**

This should show all your services with their status!

### Test SSO

1. Go to any SSO-protected service (e.g., https://tasks.fig.systems)
2. You should be redirected to https://auth.fig.systems
3. Login with your LLDAP credentials
4. You should be redirected back to the service

## Step 8: Initial Service Configuration

### Jellyfin Setup
1. Go to https://flix.fig.systems
2. Select language and create admin account
3. Add media libraries:
   - Movies: `/media/movies`
   - TV Shows: `/media/tv`
   - Music: `/media/music`
   - Photos: `/media/photos`

### Immich Setup
1. Go to https://photos.fig.systems
2. Create admin account
3. Upload some photos to test
4. Configure storage in Settings

### Sonarr/Radarr Setup
1. Go to https://sonarr.fig.systems and https://radarr.fig.systems
2. Complete initial setup wizard
3. Add indexers (for finding content)
4. Add download clients:
   - SABnzbd: http://sabnzbd:8080
   - qBittorrent: http://qbittorrent:8080
5. Configure root folders:
   - Sonarr: `/media/tv`
   - Radarr: `/media/movies`

### Jellyseerr Setup
1. Go to https://requests.fig.systems
2. Sign in with Jellyfin
3. Connect to Sonarr and Radarr
4. Configure user permissions

### Backrest Setup
1. Go to https://backup.fig.systems
2. Add Backblaze B2 repository (see [Backup Guide](./services/backup.md))
3. Create backup plan for Immich photos
4. Schedule automated backups

## Step 9: Optional Configurations

### Enable GPU Acceleration

If you have an NVIDIA GPU, see [GPU Setup Guide](./guides/gpu-setup.md).

### Configure Backups

See [Backup Operations Guide](./operations/backups.md).

### Add More Services

See [Adding Services Guide](./guides/adding-services.md).

## Next Steps

- ✅ [Set up automated backups](./operations/backups.md)
- ✅ [Configure monitoring](./operations/monitoring.md)
- ✅ [Review security settings](./guides/security.md)
- ✅ [Enable GPU acceleration](./guides/gpu-setup.md) (optional)
- ✅ [Configure media automation](./services/media-stack.md)

## Troubleshooting

If you encounter issues during setup, see:
- [Common Issues](./troubleshooting/common-issues.md)
- [FAQ](./troubleshooting/faq.md)
- [Debugging Guide](./troubleshooting/debugging.md)

## Quick Command Reference

```bash
# View all running containers
docker ps

# View logs for a service
cd compose/path/to/service
docker compose logs -f

# Restart a service
docker compose restart

# Stop a service
docker compose down

# Update and restart a service
docker compose pull
docker compose up -d

# View resource usage
docker stats
```

## Getting Help

- Check the [FAQ](./troubleshooting/faq.md)
- Review service-specific guides in [docs/services/](./services/)
- Check container logs for errors
- Verify DNS and SSL certificates

Welcome to your homelab! 🎉
