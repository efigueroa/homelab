# Service Template

This template provides a starting point for adding new services to your homelab.

## Quick Start

1. **Copy this template:**
   ```bash
   cp -r templates/service-template compose/[category]/[service-name]
   cd compose/[category]/[service-name]
   ```

2. **Choose the correct category:**
   - `compose/core/` - Infrastructure services (reverse proxy, auth, etc.)
   - `compose/media/` - Media-related services (streaming, automation)
   - `compose/services/` - Utility services (bookmarks, tasks, etc.)

3. **Update compose.yaml:**
   - Replace `service-name` with actual service name
   - Update `image:` with the correct Docker image
   - Configure environment variables
   - Update port numbers
   - Update Traefik domain (replace `service.fig.systems`)
   - Add any required volumes or dependencies

4. **Create .env file (if needed):**
   ```bash
   cp .env.example .env
   # Edit .env and set real values
   ```

5. **Update README.md:**
   - Add service to main README.md service table
   - Include service URL and description
   - Document any special configuration

6. **Test deployment:**
   ```bash
   docker compose config  # Validate syntax
   docker compose up -d   # Deploy
   docker compose logs -f # Check logs
   ```

## Template Features

### Included by Default
- ✅ Traefik integration with SSL/TLS
- ✅ Dual domain support (fig.systems + edfig.dev)
- ✅ SSO middleware support (commented)
- ✅ Network configuration (homelab external)
- ✅ Standard environment variables (PUID, PGID, TZ)
- ✅ Restart policy
- ✅ Health check template

### Optional Components (Commented)
- Database service (PostgreSQL example)
- Redis cache
- Internal network for multi-container setups
- Port exposure (prefer Traefik)
- Named volumes
- Health checks

## Common Patterns

### Simple Single-Container Service
```yaml
services:
  app:
    image: app:latest
    volumes:
      - ./config:/config
    networks:
      - homelab
    labels:
      # Traefik config
```

### Multi-Container with Database
```yaml
services:
  app:
    image: app:latest
    depends_on:
      - database
    networks:
      - homelab
      - app_internal

  database:
    image: postgres:16-alpine
    networks:
      - app_internal

networks:
  homelab:
    external: true
  app_internal:
    driver: bridge
```

### Service with Media Access
```yaml
services:
  app:
    image: app:latest
    volumes:
      - ./config:/config
      - /media/movies:/movies:ro
      - /media/books:/books:ro
```

## Checklist

Before submitting a PR with a new service:

- [ ] Service name is descriptive and lowercase with hyphens
- [ ] Docker image is from a trusted source
- [ ] All placeholder passwords use `changeme_*` format
- [ ] Traefik labels are complete (router, entrypoint, tls, rule)
- [ ] Both domains configured (fig.systems + edfig.dev)
- [ ] SSO middleware decision made (enabled/disabled with comment)
- [ ] Networks properly configured (external: true)
- [ ] Health check added (if applicable)
- [ ] Service added to README.md
- [ ] Documentation header in compose.yaml
- [ ] .env.example provided (if using env_file)
- [ ] Tested locally before committing

## Domain Selection

Choose a subdomain that makes sense:

**Common Patterns:**
- Service name: `servicename.fig.systems` (most common)
- Function-based: `monitor.fig.systems`, `backup.fig.systems`
- Alternative names: `flix.fig.systems` (Jellyfin), `requests.fig.systems` (Jellyseerr)

**Reserved Domains:**
- `auth.fig.systems` - Tinyauth
- `lldap.fig.systems` - LLDAP
- `traefik.fig.systems` - Traefik dashboard
- See README.md for complete list

## Network Configuration

### Single Container
```yaml
networks:
  homelab:
    external: true
```

### Multi-Container (with internal network)
```yaml
networks:
  homelab:
    external: true  # For Traefik access
  service_internal:
    name: service_internal
    driver: bridge  # For inter-container communication
```

### Traefik Network Selection
If using multiple networks, specify which Traefik should use:
```yaml
labels:
  traefik.docker.network: homelab
```

## Volume Patterns

### Configuration Only
```yaml
volumes:
  - ./config:/config
```

### With Data Storage
```yaml
volumes:
  - ./config:/config
  - ./data:/data
```

### With Media Access (read-only recommended)
```yaml
volumes:
  - ./config:/config
  - /media/movies:/movies:ro
  - /media/tv:/tv:ro
```

### With Database
```yaml
volumes:
  - ./db:/var/lib/postgresql/data
```

## Troubleshooting

### Service won't start
```bash
# Check logs
docker compose logs app

# Validate compose syntax
docker compose config

# Check network exists
docker network ls | grep homelab
```

### Can't access via domain
```bash
# Check Traefik is running
docker ps | grep traefik

# Check Traefik logs
docker logs traefik

# Verify DNS points to server
dig service.fig.systems

# Check SSL certificate
curl -I https://service.fig.systems
```

### Permission errors
```bash
# Check PUID/PGID match your user
id

# Fix ownership
sudo chown -R 1000:1000 ./config ./data
```

## Examples

See these services for reference:
- **Simple:** `compose/services/filebrowser/`
- **With database:** `compose/services/vikunja/`
- **Multi-container:** `compose/media/frontend/immich/`
- **Media service:** `compose/media/frontend/jellyfin/`

## Resources

- [Docker Compose Docs](https://docs.docker.com/compose/)
- [Traefik Docker Provider](https://doc.traefik.io/traefik/providers/docker/)
- [LinuxServer.io Images](https://fleet.linuxserver.io/)
- [Awesome Selfhosted](https://github.com/awesome-selfhosted/awesome-selfhosted)
