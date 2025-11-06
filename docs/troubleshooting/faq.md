# Frequently Asked Questions (FAQ)

Common questions and answers about the homelab setup.

## General Questions

### Q: What is this homelab setup?

**A:** This is a GitOps-based infrastructure for self-hosting services using Docker Compose. It includes:
- 20+ pre-configured services (media, productivity, utilities)
- Automatic SSL/TLS with Let's Encrypt via Traefik
- Single Sign-On (SSO) with LLDAP and Tinyauth
- Automated backups with Backrest
- Service discovery dashboard with Homarr

### Q: What are the minimum hardware requirements?

**A:**
- **CPU**: 2+ cores (4+ recommended)
- **RAM**: 8GB minimum (16GB+ recommended)
- **Storage**: 100GB for containers, additional space for media
- **Network**: Static IP recommended, ports 80 and 443 accessible
- **GPU** (Optional): NVIDIA GPU for hardware transcoding

### Q: Do I need my own domain name?

**A:** Yes, you need at least one domain (two configured by default: `fig.systems` and `edfig.dev`). You can:
- Register a domain from any registrar
- Update all compose files to use your domain
- Configure wildcard DNS (`*.yourdomain.com`)

### Q: Can I run this on a Raspberry Pi?

**A:** Partially. ARM64 architecture is supported by most services, but:
- Performance will be limited
- No GPU acceleration available
- Some services may not have ARM images
- 8GB RAM minimum recommended (Pi 4 or Pi 5)

### Q: How much does this cost to run?

**A:**
- **Server**: $0 (if using existing hardware) or $5-20/month (VPS)
- **Domain**: $10-15/year
- **Backblaze B2**: ~$0.60/month for 100GB photos
- **Electricity**: Varies by hardware and location
- **Total**: $15-30/year minimum

## Setup Questions

### Q: Why won't my services start?

**A:** Common causes:
1. **Environment variables not set**: Check for `changeme_*` in `.env` files
2. **Ports already in use**: Check if 80/443 are available
3. **Network not created**: Run `docker network create homelab`
4. **DNS not configured**: Services need valid DNS records
5. **Insufficient resources**: Check RAM and disk space

**Debug:**
```bash
cd compose/path/to/service
docker compose logs
docker compose ps
```

### Q: How do I know if everything is working?

**A:** Check these indicators:
1. **All containers running**: `docker ps` shows all services
2. **SSL certificates valid**: Visit https://home.fig.systems (no cert errors)
3. **Dashboard accessible**: Homarr shows all services
4. **SSO working**: Can login to protected services
5. **No errors in logs**: `docker compose logs` shows no critical errors

### Q: What order should I deploy services?

**A:** Follow this order:
1. **Core**: Traefik → LLDAP → Tinyauth
2. **Configure**: Create LLDAP users
3. **Media**: Jellyfin → Immich → Jellyseerr → Sonarr → Radarr → Downloaders
4. **Utility**: Homarr → Backrest → Everything else

### Q: Do I need to configure all 20 services?

**A:** No! Deploy only what you need:
- **Core** (required): Traefik, LLDAP, Tinyauth
- **Media** (optional): Jellyfin, Immich, Sonarr, Radarr
- **Utility** (pick what you want): Homarr, Backrest, Linkwarden, Vikunja, etc.

## Configuration Questions

### Q: What secrets do I need to change?

**A:** Search for `changeme_*` in all `.env` files:
```bash
grep -r "changeme_" compose/
```

Critical secrets:
- **LLDAP_LDAP_USER_PASS**: Admin password for LLDAP
- **LLDAP_JWT_SECRET**: 64-character hex string
- **SESSION_SECRET**: 64-character hex string for Tinyauth
- **DB_PASSWORD**: Database passwords (Immich, Vikunja, Linkwarden)
- **NEXTAUTH_SECRET**: NextAuth secret for Linkwarden
- **VIKUNJA_SERVICE_JWTSECRET**: JWT secret for Vikunja

### Q: How do I generate secure secrets?

**A:** Use these commands:

```bash
# 64-character hex (for JWT secrets, session secrets)
openssl rand -hex 32

# 32-character password (for databases)
openssl rand -base64 32 | tr -d /=+ | cut -c1-32

# 32-character hex (for API keys)
openssl rand -hex 16
```

See [Secrets Management Guide](../guides/secrets-management.md) for details.

### Q: Can I change the domains from fig.systems to my own?

**A:** Yes! You need to:
1. Find and replace in all `compose.yaml` files:
   ```bash
   find compose -name "compose.yaml" -exec sed -i 's/fig\.systems/yourdomain.com/g' {} \;
   find compose -name "compose.yaml" -exec sed -i 's/edfig\.dev/yourotherdomain.com/g' {} \;
   ```
2. Update DNS records to point to your server
3. Update `.env` files with new URLs (e.g., `NEXTAUTH_URL`, `VIKUNJA_SERVICE_PUBLICURL`)

### Q: Do all passwords need to match?

**A:** No, but some must match:
- **LLDAP_LDAP_USER_PASS** must equal **LDAP_BIND_PASSWORD** (in tinyauth)
- **VIKUNJA_DATABASE_PASSWORD** must equal **POSTGRES_PASSWORD** (in vikunja)
- **Linkwarden POSTGRES_PASSWORD** is used in DATABASE_URL

All other passwords should be unique!

## SSL/TLS Questions

### Q: Why am I getting SSL certificate errors?

**A:** Common causes:
1. **DNS not configured**: Ensure domains point to your server
2. **Ports not accessible**: Let's Encrypt needs port 80 for HTTP challenge
3. **Rate limiting**: Let's Encrypt has rate limits (5 certs per domain/week)
4. **First startup**: Certs take a few minutes to generate

**Debug:**
```bash
docker logs traefik | grep -i error
docker logs traefik | grep -i certificate
```

### Q: How long do SSL certificates last?

**A:** Let's Encrypt certificates:
- Valid for 90 days
- Traefik auto-renews at 30 days before expiration
- Renewals happen automatically in the background

### Q: Can I use my own SSL certificates?

**A:** Yes, but it requires modifying Traefik configuration. The default Let's Encrypt setup is recommended.

## SSO Questions

### Q: What is SSO and do I need it?

**A:** SSO (Single Sign-On) lets you log in once and access all services:
- **LLDAP**: Stores users and passwords
- **Tinyauth**: Authenticates users before allowing service access
- **Benefits**: One login for all services, centralized user management
- **Optional**: Some services can work without SSO (have their own auth)

### Q: Why can't I log into SSO-protected services?

**A:** Check:
1. **LLDAP is running**: `docker ps | grep lldap`
2. **Tinyauth is running**: `docker ps | grep tinyauth`
3. **User exists in LLDAP**: Go to https://lldap.fig.systems and verify
4. **Passwords match**: LDAP_BIND_PASSWORD = LLDAP_LDAP_USER_PASS
5. **User in correct group**: Check user is in `lldap_admin` group

**Debug:**
```bash
cd compose/core/tinyauth
docker compose logs -f
```

### Q: Can I disable SSO for a service?

**A:** Yes! Comment out the middleware line in compose.yaml:
```yaml
# traefik.http.routers.servicename.middlewares: tinyauth
```

Then restart the service:
```bash
docker compose up -d
```

### Q: How do I reset my LLDAP admin password?

**A:**
1. Stop LLDAP: `cd compose/core/lldap && docker compose down`
2. Update `LLDAP_LDAP_USER_PASS` in `.env`
3. Remove the database: `rm -rf data/`
4. Restart: `docker compose up -d`
5. Recreate users in LLDAP UI

⚠️ **Warning**: This deletes all users!

## Service-Specific Questions

### Q: Jellyfin shows "Playback Error" - what's wrong?

**A:** Common causes:
1. **Media file corrupt**: Test file with VLC
2. **Permissions**: Check file ownership (`ls -la /media/movies`)
3. **Codec not supported**: Enable transcoding or use different file
4. **GPU not configured**: If using GPU, verify NVIDIA Container Toolkit

### Q: Immich won't upload photos - why?

**A:** Check:
1. **Database connected**: `docker logs immich_postgres`
2. **Upload directory writable**: Check permissions on `./upload`
3. **Disk space**: `df -h`
4. **File size limits**: Check browser console for errors

### Q: Why isn't Homarr showing my services?

**A:** Homarr needs:
1. **Docker socket access**: Volume mount `/var/run/docker.sock`
2. **Labels on services**: Each service needs `homarr.name` label
3. **Same network**: Homarr must be on `homelab` network
4. **Time to detect**: Refresh page or wait 30 seconds

### Q: Backrest shows "Repository not initialized" - what do I do?

**A:**
1. Go to https://backup.fig.systems
2. Click "Add Repository"
3. Configure Backblaze B2 settings
4. Click "Initialize Repository"

See [Backup Guide](../services/backup.md) for detailed setup.

### Q: Sonarr/Radarr can't find anything - help!

**A:**
1. **Add indexers**: Settings → Indexers → Add indexer
2. **Configure download client**: Settings → Download Clients → Add
3. **Set root folder**: Series/Movies → Add Root Folder → `/media/tv` or `/media/movies`
4. **Test indexers**: Settings → Indexers → Test

### Q: qBittorrent shows "Unauthorized" - what's the password?

**A:** Default credentials:
- Username: `admin`
- Password: `adminadmin`

⚠️ **Change this immediately** in qBittorrent settings!

## Media Questions

### Q: Where should I put my media files?

**A:** Use the /media directory structure:
- Movies: `/media/movies/Movie Name (Year)/movie.mkv`
- TV: `/media/tv/Show Name/Season 01/episode.mkv`
- Music: `/media/music/Artist/Album/song.flac`
- Photos: `/media/photos/` (any structure)
- Books: `/media/books/` (any structure)

### Q: How do I add more media storage?

**A:**
1. Mount additional drive to `/media2` (or any path)
2. Update compose files to include new volume:
   ```yaml
   volumes:
     - /media:/media:ro
     - /media2:/media2:ro  # Add this
   ```
3. Restart service: `docker compose up -d`
4. Add new library in service UI

### Q: Can Sonarr/Radarr automatically download shows/movies?

**A:** Yes! That's their purpose:
1. Add indexers (for searching)
2. Add download client (SABnzbd or qBittorrent)
3. Add a series/movie
4. Enable monitoring
5. Sonarr/Radarr will search, download, and organize automatically

### Q: How do I enable hardware transcoding in Jellyfin?

**A:** See [GPU Setup Guide](../guides/gpu-setup.md) for full instructions.

Quick steps:
1. Install NVIDIA Container Toolkit on host
2. Uncomment GPU sections in `jellyfin/compose.yaml`
3. Restart Jellyfin
4. Enable in Jellyfin: Dashboard → Playback → Hardware Acceleration → NVIDIA NVENC

## Network Questions

### Q: Can I access services only from my local network?

**A:** Yes, don't expose ports 80/443 to internet:
1. Use firewall to block external access
2. Use local DNS (Pi-hole, AdGuard Home)
3. Point domains to local IP (192.168.x.x)
4. Use self-signed certs or no HTTPS

**Or** use Traefik's IP allowlist middleware.

### Q: Can I use a VPN with these services?

**A:** Yes, options:
1. **VPN on download clients**: Add VPN container for qBittorrent/SABnzbd
2. **VPN to access homelab**: Use WireGuard/Tailscale to access from anywhere
3. **VPN for entire server**: All traffic goes through VPN (not recommended)

### Q: Why can't I access services from outside my network?

**A:** Check:
1. **Port forwarding**: Ports 80 and 443 forwarded to homelab server
2. **Firewall**: Allow ports 80/443 through firewall
3. **DNS**: Domains point to your public IP
4. **ISP**: Some ISPs block ports 80/443 (use CloudFlare Tunnel)

## Backup Questions

### Q: What should I backup?

**A:** Priority order:
1. **High**: Immich photos (`compose/media/frontend/immich/upload`)
2. **High**: Configuration files (all `.env` files, compose files)
3. **Medium**: Service data directories (`./config`, `./data` in each service)
4. **Low**: Media files (usually have source elsewhere)

### Q: How do I restore from backup?

**A:** See [Backup Operations Guide](../operations/backups.md).

Quick steps:
1. Install fresh homelab setup
2. Restore `.env` files and configs
3. Use Backrest to restore data
4. Restart services

### Q: Does Backrest backup everything automatically?

**A:** Only what you configure:
- Default: Immich photos and homelab configs
- Add more paths in `backrest/compose.yaml` volumes
- Create backup plans in Backrest UI for each path

## Performance Questions

### Q: Services are running slow - how do I optimize?

**A:**
1. **Check resources**: `docker stats` - are you out of RAM/CPU?
2. **Reduce services**: Stop unused services
3. **Use SSD**: Move Docker to SSD storage
4. **Add RAM**: Minimum 8GB, 16GB+ recommended
5. **Enable GPU**: For Jellyfin and Immich

### Q: Docker is using too much disk space - what do I do?

**A:**
```bash
# Check Docker disk usage
docker system df

# Clean up
docker system prune -a --volumes

# WARNING: This removes all stopped containers and unused volumes!
```

Better approach - clean specific services:
```bash
cd compose/path/to/service
docker compose down
docker volume rm $(docker volume ls -q | grep servicename)
docker compose up -d
```

### Q: How do I limit RAM/CPU for a service?

**A:** Add resource limits to compose.yaml:
```yaml
services:
  servicename:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 4G
        reservations:
          memory: 2G
```

## Update Questions

### Q: How do I update a service?

**A:**
```bash
cd compose/path/to/service
docker compose pull
docker compose up -d
```

See [Updates Guide](../operations/updates.md) for details.

### Q: How often should I update?

**A:**
- **Security updates**: Weekly
- **Feature updates**: Monthly
- **Major versions**: When stable

Use Watchtower for automatic updates (optional).

### Q: Will updating break my configuration?

**A:** Usually no, but:
- Always backup before major updates
- Check release notes for breaking changes
- Test in staging environment if critical

## Security Questions

### Q: Is this setup secure?

**A:** Reasonably secure with best practices:
- ✅ SSL/TLS encryption
- ✅ SSO authentication
- ✅ Secrets in environment files
- ⚠️ Some services exposed to internet
- ⚠️ Depends on keeping services updated

See [Security Guide](../guides/security.md) for hardening.

### Q: Should I expose my homelab to the internet?

**A:** Depends on your risk tolerance:
- **Yes**: Convenient access from anywhere, Let's Encrypt works
- **No**: More secure, requires VPN for external access
- **Hybrid**: Expose only essential services, use VPN for sensitive ones

### Q: What if someone gets my LLDAP password?

**A:** They can access all SSO-protected services. Mitigations:
- Use strong, unique passwords
- Enable 2FA where supported
- Review LLDAP access logs
- Use fail2ban to block brute force
- Consider VPN-only access

## Troubleshooting

For specific error messages and debugging, see:
- [Common Issues](./common-issues.md)
- [Debugging Guide](./debugging.md)

Still stuck? Check:
1. Service logs: `docker compose logs`
2. Traefik logs: `docker logs traefik`
3. Container status: `docker ps -a`
4. Network connectivity: `docker network inspect homelab`
