# Common Issues and Solutions

This guide covers the most common problems you might encounter and how to fix them.

## Table of Contents
- [Service Won't Start](#service-wont-start)
- [SSL/TLS Certificate Errors](#ssltls-certificate-errors)
- [SSO Authentication Issues](#sso-authentication-issues)
- [Access Issues](#access-issues)
- [Performance Problems](#performance-problems)
- [Database Errors](#database-errors)
- [Network Issues](#network-issues)
- [GPU Problems](#gpu-problems)

## Service Won't Start

### Symptom
Container exits immediately or shows "Exited (1)" status.

### Diagnosis
```bash
cd ~/homelab/compose/path/to/service

# Check container status
docker compose ps

# View logs
docker compose logs

# Check for specific errors
docker compose logs | grep -i error
```

### Common Causes and Fixes

#### 1. Environment Variables Not Set

**Error in logs:**
```
Error: POSTGRES_PASSWORD is not set
Error: required environment variable 'XXX' is missing
```

**Fix:**
```bash
# Check .env file exists
ls -la .env

# Check for changeme_ values
grep "changeme_" .env

# Update with proper secrets (see secrets guide)
nano .env

# Restart
docker compose up -d
```

#### 2. Port Already in Use

**Error in logs:**
```
Error: bind: address already in use
Error: failed to bind to port 80: address already in use
```

**Fix:**
```bash
# Find what's using the port
sudo netstat -tulpn | grep :80
sudo netstat -tulpn | grep :443

# Stop conflicting service
sudo systemctl stop apache2  # Example
sudo systemctl stop nginx    # Example

# Or change port in compose.yaml
```

#### 3. Network Not Created

**Error in logs:**
```
network homelab declared as external, but could not be found
```

**Fix:**
```bash
# Create network
docker network create homelab

# Verify
docker network ls | grep homelab

# Restart service
docker compose up -d
```

#### 4. Volume Permission Issues

**Error in logs:**
```
Permission denied: '/config'
mkdir: cannot create directory '/data': Permission denied
```

**Fix:**
```bash
# Check directory ownership
ls -la ./config ./data

# Fix ownership (replace 1000:1000 with your UID:GID)
sudo chown -R 1000:1000 ./config ./data

# Restart
docker compose up -d
```

#### 5. Dependency Not Running

**Error in logs:**
```
Failed to connect to database
Connection refused: postgres:5432
```

**Fix:**
```bash
# Start dependency first
cd ~/homelab/compose/path/to/dependency
docker compose up -d

# Wait for it to be healthy
docker compose logs -f

# Then start the service
cd ~/homelab/compose/path/to/service
docker compose up -d
```

## SSL/TLS Certificate Errors

### Symptom
Browser shows "Your connection is not private" or "NET::ERR_CERT_AUTHORITY_INVALID"

### Diagnosis
```bash
# Check Traefik logs
docker logs traefik | grep -i certificate
docker logs traefik | grep -i letsencrypt
docker logs traefik | grep -i error

# Test certificate
echo | openssl s_client -servername home.fig.systems -connect home.fig.systems:443 2>/dev/null | openssl x509 -noout -dates
```

### Common Causes and Fixes

#### 1. DNS Not Configured

**Fix:**
```bash
# Test DNS resolution
dig home.fig.systems +short

# Should return your server's IP
# If not, configure DNS A records:
# *.fig.systems -> YOUR_SERVER_IP
```

#### 2. Port 80 Not Accessible

Let's Encrypt needs port 80 for HTTP-01 challenge.

**Fix:**
```bash
# Test from external network
curl -I http://home.fig.systems

# Check firewall
sudo ufw status
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Check port forwarding on router
# Ensure ports 80 and 443 are forwarded to server
```

#### 3. Rate Limiting

Let's Encrypt has limits: 5 certificates per domain per week.

**Fix:**
```bash
# Check Traefik logs for rate limit errors
docker logs traefik | grep -i "rate limit"

# Wait for rate limit to reset (1 week)
# Or use Let's Encrypt staging environment for testing

# Enable staging in traefik/compose.yaml:
# - --certificatesresolvers.letsencrypt.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory
```

#### 4. First Startup - Certificates Not Yet Generated

**Fix:**
```bash
# Wait 2-5 minutes for certificate generation
docker logs traefik -f

# Look for:
# "Certificate obtained for domain"
```

#### 5. Certificate Expired

Traefik should auto-renew, but if manual renewal needed:

**Fix:**
```bash
# Remove old certificates
cd ~/homelab/compose/core/traefik
rm -rf ./acme.json

# Restart Traefik
docker compose restart

# Wait for new certificates
docker logs traefik -f
```

## SSO Authentication Issues

### Symptom
- Can't login to SSO-protected services
- Redirected to auth page but login fails
- "Invalid credentials" error

### Diagnosis
```bash
# Check LLDAP is running
docker ps | grep lldap

# Check Tinyauth is running
docker ps | grep tinyauth

# View logs
docker logs lldap
docker logs tinyauth
```

### Common Causes and Fixes

#### 1. Password Mismatch

LDAP_BIND_PASSWORD must match LLDAP_LDAP_USER_PASS.

**Fix:**
```bash
# Check both passwords
grep LLDAP_LDAP_USER_PASS ~/homelab/compose/core/lldap/.env
grep LDAP_BIND_PASSWORD ~/homelab/compose/core/tinyauth/.env

# They must be EXACTLY the same!

# If different, update tinyauth/.env
cd ~/homelab/compose/core/tinyauth
nano .env
# Set LDAP_BIND_PASSWORD to match LLDAP_LDAP_USER_PASS

# Restart Tinyauth
docker compose restart
```

#### 2. User Doesn't Exist in LLDAP

**Fix:**
```bash
# Access LLDAP web UI
# Go to: https://lldap.fig.systems

# Login with admin credentials
# Username: admin
# Password: <your LLDAP_LDAP_USER_PASS>

# Create user:
# - Click "Create user"
# - Set username, email, password
# - Add to "lldap_admin" group

# Try logging in again
```

#### 3. LLDAP or Tinyauth Not Running

**Fix:**
```bash
# Start LLDAP
cd ~/homelab/compose/core/lldap
docker compose up -d

# Wait for it to be ready
docker compose logs -f

# Start Tinyauth
cd ~/homelab/compose/core/tinyauth
docker compose up -d
docker compose logs -f
```

#### 4. Network Issue Between Tinyauth and LLDAP

**Fix:**
```bash
# Test connection
docker exec tinyauth nc -zv lldap 3890

# Should show: Connection to lldap 3890 port [tcp/*] succeeded!

# If not, check both are on homelab network
docker network inspect homelab
```

## Access Issues

### Symptom
- Can't access service from browser
- Connection timeout
- "This site can't be reached"

### Diagnosis
```bash
# Test from server
curl -I https://home.fig.systems

# Test DNS
dig home.fig.systems +short

# Check container is running
docker ps | grep servicename

# Check Traefik routing
docker logs traefik | grep servicename
```

### Common Causes and Fixes

#### 1. Service Not Running

**Fix:**
```bash
cd ~/homelab/compose/path/to/service
docker compose up -d
docker compose logs -f
```

#### 2. Traefik Not Running

**Fix:**
```bash
cd ~/homelab/compose/core/traefik
docker compose up -d
docker compose logs -f
```

#### 3. DNS Not Resolving

**Fix:**
```bash
# Check DNS
dig service.fig.systems +short

# Should return your server IP
# If not, add/update DNS A record
```

#### 4. Firewall Blocking

**Fix:**
```bash
# Check firewall
sudo ufw status

# Allow if needed
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

#### 5. Wrong Traefik Labels

**Fix:**
```bash
# Check compose.yaml has correct labels
cd ~/homelab/compose/path/to/service
cat compose.yaml | grep -A 10 "labels:"

# Should have:
# traefik.enable: true
# traefik.http.routers.servicename.rule: Host(`service.fig.systems`)
# etc.
```

## Performance Problems

### Symptom
- Services running slowly
- High CPU/RAM usage
- System unresponsive

### Diagnosis
```bash
# Overall system
htop

# Docker resources
docker stats

# Disk usage
df -h
docker system df
```

### Common Causes and Fixes

#### 1. Insufficient RAM

**Fix:**
```bash
# Check RAM usage
free -h

# If low, either:
# 1. Add more RAM
# 2. Stop unused services
# 3. Add resource limits to compose files

# Example resource limit:
deploy:
  resources:
    limits:
      memory: 2G
    reservations:
      memory: 1G
```

#### 2. Disk Full

**Fix:**
```bash
# Check disk usage
df -h

# Clean Docker
docker system prune -a

# Remove old logs
sudo journalctl --vacuum-time=7d

# Check media folder
du -sh /media/*
```

#### 3. Too Many Services Running

**Fix:**
```bash
# Stop unused services
cd ~/homelab/compose/services/unused-service
docker compose down

# Or deploy only what you need
```

#### 4. Database Not Optimized

**Fix:**
```bash
# For postgres services, add to .env:
POSTGRES_INITDB_ARGS=--data-checksums

# Increase shared buffers (if enough RAM):
# Edit compose.yaml, add to postgres:
command: postgres -c shared_buffers=256MB -c max_connections=200
```

## Database Errors

### Symptom
- "Connection refused" to database
- "Authentication failed for user"
- "Database does not exist"

### Diagnosis
```bash
# Check database container
docker ps | grep postgres

# View database logs
docker logs <postgres_container_name>

# Test connection from app
docker exec <app_container> nc -zv <db_container> 5432
```

### Common Causes and Fixes

#### 1. Password Mismatch

**Fix:**
```bash
# Check passwords match in .env
cat .env | grep PASSWORD

# For example, in Vikunja:
# VIKUNJA_DATABASE_PASSWORD and POSTGRES_PASSWORD must match!

# Update if needed
nano .env
docker compose down
docker compose up -d
```

#### 2. Database Not Initialized

**Fix:**
```bash
# Remove database and reinitialize
docker compose down
rm -rf ./db/  # CAREFUL: This deletes all data!
docker compose up -d
```

#### 3. Database Still Starting

**Fix:**
```bash
# Wait for database to be ready
docker logs <postgres_container> -f

# Look for "database system is ready to accept connections"

# Then restart app
docker compose restart <app_service>
```

## Network Issues

### Symptom
- Containers can't communicate
- "Connection refused" between services

### Diagnosis
```bash
# Inspect network
docker network inspect homelab

# Test connectivity
docker exec container1 ping container2
docker exec container1 nc -zv container2 PORT
```

### Common Causes and Fixes

#### 1. Containers Not on Same Network

**Fix:**
```bash
# Check compose.yaml has networks section
networks:
  homelab:
    external: true

# Ensure service is using the network
services:
  servicename:
    networks:
      - homelab
```

#### 2. Network Doesn't Exist

**Fix:**
```bash
docker network create homelab
docker compose up -d
```

#### 3. DNS Resolution Between Containers

**Fix:**
```bash
# Use container name, not localhost
# Wrong: http://localhost:5432
# Right:  http://postgres:5432

# Or use service name from compose.yaml
```

## GPU Problems

### Symptom
- "No hardware acceleration available"
- GPU not detected in container
- "Failed to open GPU"

### Diagnosis
```bash
# Check GPU on host
nvidia-smi

# Check GPU in container
docker exec jellyfin nvidia-smi

# Check Docker GPU runtime
docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi
```

### Common Causes and Fixes

#### 1. NVIDIA Container Toolkit Not Installed

**Fix:**
```bash
# Install toolkit
sudo apt install nvidia-container-toolkit

# Configure runtime
sudo nvidia-ctk runtime configure --runtime=docker

# Restart Docker
sudo systemctl restart docker
```

#### 2. Runtime Not Specified in Compose

**Fix:**
```bash
# Edit compose.yaml
nano compose.yaml

# Uncomment:
runtime: nvidia
deploy:
  resources:
    reservations:
      devices:
        - driver: nvidia
          count: all
          capabilities: [gpu]

# Restart
docker compose up -d
```

#### 3. GPU Already in Use

**Fix:**
```bash
# Check processes using GPU
nvidia-smi

# Kill process if needed
sudo kill <PID>

# Restart service
docker compose restart
```

#### 4. GPU Not Passed Through to VM (Proxmox)

**Fix:**
```bash
# From Proxmox host, check GPU passthrough
lspci | grep -i nvidia

# From VM, check GPU visible
lspci | grep -i nvidia

# If not visible, reconfigure passthrough (see GPU guide)
```

## Getting More Help

If your issue isn't listed here:

1. **Check service-specific logs**:
   ```bash
   cd ~/homelab/compose/path/to/service
   docker compose logs --tail=200
   ```

2. **Search container logs for errors**:
   ```bash
   docker compose logs | grep -i error
   docker compose logs | grep -i fail
   ```

3. **Check FAQ**: See [FAQ](./faq.md)

4. **Debugging Guide**: See [Debugging Guide](./debugging.md)

5. **Service Documentation**: Check service's official documentation

---

**Most issues can be solved by checking logs and environment variables!**
