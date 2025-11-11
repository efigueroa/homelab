# AlmaLinux 9.6 VM Setup Guide

Complete setup guide for the homelab VM on AlmaLinux 9.6 running on Proxmox VE 9.

## Hardware Context

- **Host**: Proxmox VE 9 (Debian 13 based)
  - CPU: AMD Ryzen 5 7600X (6C/12T, 5.3 GHz boost)
  - GPU: NVIDIA GTX 1070 (8GB VRAM)
  - RAM: 32GB DDR5

- **VM Allocation**:
  - OS: AlmaLinux 9.6 (RHEL 9 compatible)
  - CPU: 8 vCPUs
  - RAM: 24GB
  - Disk: 500GB+ (expandable)
  - GPU: GTX 1070 (PCIe passthrough)

## Proxmox VM Creation

### 1. Create VM

```bash
# On Proxmox host
qm create 100 \
  --name homelab \
  --memory 24576 \
  --cores 8 \
  --cpu host \
  --sockets 1 \
  --net0 virtio,bridge=vmbr0 \
  --scsi0 local-lvm:500 \
  --ostype l26 \
  --boot order=scsi0

# Attach AlmaLinux ISO
qm set 100 --ide2 local:iso/AlmaLinux-9.6-x86_64-dvd.iso,media=cdrom

# Enable UEFI
qm set 100 --bios ovmf --efidisk0 local-lvm:1
```

### 2. GPU Passthrough

**Find GPU PCI address:**
```bash
lspci | grep -i nvidia
# Example output: 01:00.0 VGA compatible controller: NVIDIA Corporation GP104 [GeForce GTX 1070]
```

**Enable IOMMU in Proxmox:**

Edit `/etc/default/grub`:
```bash
# For AMD CPU (Ryzen 5 7600X)
GRUB_CMDLINE_LINUX_DEFAULT="quiet amd_iommu=on iommu=pt"
```

Update GRUB and reboot:
```bash
update-grub
reboot
```

**Verify IOMMU:**
```bash
dmesg | grep -e DMAR -e IOMMU
# Should show IOMMU enabled
```

**Add GPU to VM:**

Edit `/etc/pve/qemu-server/100.conf`:
```
hostpci0: 0000:01:00,pcie=1,x-vga=1
```

Or via command:
```bash
qm set 100 --hostpci0 0000:01:00,pcie=1,x-vga=1
```

**Blacklist GPU on host:**

Edit `/etc/modprobe.d/blacklist-nvidia.conf`:
```
blacklist nouveau
blacklist nvidia
blacklist nvidia_drm
blacklist nvidia_modeset
blacklist nvidia_uvm
```

Update initramfs:
```bash
update-initramfs -u
reboot
```

## AlmaLinux Installation

### 1. Install AlmaLinux 9.6

Start VM and follow installer:
1. **Language**: English (US)
2. **Installation Destination**: Use all space, automatic partitioning
3. **Network**: Enable and set hostname to `homelab.fig.systems`
4. **Software Selection**: Minimal Install
5. **Root Password**: Set strong password
6. **User Creation**: Create admin user (e.g., `homelab`)

### 2. Post-Installation Configuration

```bash
# SSH into VM
ssh homelab@<vm-ip>

# Update system
sudo dnf update -y

# Install essential tools
sudo dnf install -y \
  vim \
  git \
  curl \
  wget \
  htop \
  ncdu \
  tree \
  tmux \
  bind-utils \
  net-tools \
  firewalld

# Enable and configure firewall
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### 3. Configure Static IP (Optional)

```bash
# Find connection name
nmcli connection show

# Set static IP (example: 192.168.1.100)
sudo nmcli connection modify "System eth0" \
  ipv4.addresses 192.168.1.100/24 \
  ipv4.gateway 192.168.1.1 \
  ipv4.dns "1.1.1.1,8.8.8.8" \
  ipv4.method manual

# Restart network
sudo nmcli connection down "System eth0"
sudo nmcli connection up "System eth0"
```

## Docker Installation

### 1. Install Docker Engine

```bash
# Remove old versions
sudo dnf remove docker \
  docker-client \
  docker-client-latest \
  docker-common \
  docker-latest \
  docker-latest-logrotate \
  docker-logrotate \
  docker-engine

# Add Docker repository
sudo dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

# Install Docker
sudo dnf install -y \
  docker-ce \
  docker-ce-cli \
  containerd.io \
  docker-buildx-plugin \
  docker-compose-plugin

# Start Docker
sudo systemctl enable --now docker

# Verify
sudo docker run hello-world
```

### 2. Configure Docker

**Add user to docker group:**
```bash
sudo usermod -aG docker $USER
newgrp docker

# Verify (no sudo needed)
docker ps
```

**Configure Docker daemon:**

Create `/etc/docker/daemon.json`:
```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "storage-driver": "overlay2",
  "features": {
    "buildkit": true
  }
}
```

Restart Docker:
```bash
sudo systemctl restart docker
```

## NVIDIA GPU Setup

### 1. Install NVIDIA Drivers

```bash
# Add EPEL repository
sudo dnf install -y epel-release

# Add NVIDIA repository
sudo dnf config-manager --add-repo \
  https://developer.download.nvidia.com/compute/cuda/repos/rhel9/x86_64/cuda-rhel9.repo

# Install drivers
sudo dnf install -y \
  nvidia-driver \
  nvidia-driver-cuda \
  nvidia-settings \
  nvidia-persistenced

# Reboot to load drivers
sudo reboot
```

### 2. Verify GPU

```bash
# Check driver version
nvidia-smi

# Expected output:
# +-----------------------------------------------------------------------------+
# | NVIDIA-SMI 535.xx.xx    Driver Version: 535.xx.xx    CUDA Version: 12.2    |
# |-------------------------------+----------------------+----------------------+
# | GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
# |   0  GeForce GTX 1070    Off  | 00000000:01:00.0 Off |                  N/A |
# +-------------------------------+----------------------+----------------------+
```

### 3. Install NVIDIA Container Toolkit

```bash
# Add NVIDIA Container Toolkit repository
sudo dnf config-manager --add-repo \
  https://nvidia.github.io/libnvidia-container/stable/rpm/nvidia-container-toolkit.repo

# Install toolkit
sudo dnf install -y nvidia-container-toolkit

# Configure Docker to use nvidia runtime
sudo nvidia-ctk runtime configure --runtime=docker

# Restart Docker
sudo systemctl restart docker

# Test GPU in container
docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi
```

## Storage Setup

### 1. Create Media Directory

```bash
# Create media directory structure
sudo mkdir -p /media/{tv,movies,music,photos,books,audiobooks,comics,homemovies}
sudo mkdir -p /media/{downloads,complete,incomplete}

# Set ownership
sudo chown -R $USER:$USER /media

# Set permissions
chmod -R 755 /media
```

### 2. Mount Additional Storage (Optional)

If using separate disk for media:

```bash
# Find disk
lsblk

# Format disk (example: /dev/sdb)
sudo mkfs.ext4 /dev/sdb

# Get UUID
sudo blkid /dev/sdb

# Add to /etc/fstab
echo "UUID=<uuid> /media ext4 defaults,nofail 0 2" | sudo tee -a /etc/fstab

# Mount
sudo mount -a
```

## Homelab Repository Setup

### 1. Clone Repository

```bash
# Create workspace
mkdir -p ~/homelab
cd ~/homelab

# Clone repository
git clone https://github.com/efigueroa/homelab.git .

# Or if using SSH
git clone git@github.com:efigueroa/homelab.git .
```

### 2. Create Docker Network

```bash
# Create homelab network
docker network create homelab

# Verify
docker network ls | grep homelab
```

### 3. Configure Environment Variables

```bash
# Generate secrets for all services
cd ~/homelab

# LLDAP
cd compose/core/lldap
openssl rand -hex 32 > /tmp/lldap_jwt_secret
openssl rand -base64 32 | tr -d /=+ | cut -c1-32 > /tmp/lldap_pass
# Update .env with generated secrets

# Tinyauth
cd ../tinyauth
openssl rand -hex 32 > /tmp/tinyauth_session
# Update .env (LDAP_BIND_PASSWORD must match LLDAP)

# Continue for all services...
```

See [`docs/guides/secrets-management.md`](../guides/secrets-management.md) for complete guide.

## SELinux Configuration

AlmaLinux uses SELinux by default. Configure for Docker:

```bash
# Check SELinux status
getenforce
# Should show: Enforcing

# Allow Docker to access bind mounts
sudo setsebool -P container_manage_cgroup on

# If you encounter permission issues:
# Option 1: Add SELinux context to directories
sudo chcon -R -t container_file_t ~/homelab/compose
sudo chcon -R -t container_file_t /media

# Option 2: Use :Z flag in docker volumes (auto-relabels)
# Example: ./data:/data:Z

# Option 3: Set SELinux to permissive (not recommended)
# sudo setenforce 0
```

## System Tuning

### 1. Increase File Limits

```bash
# Add to /etc/security/limits.conf
echo "* soft nofile 65536" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65536" | sudo tee -a /etc/security/limits.conf

# Add to /etc/sysctl.conf
echo "fs.file-max = 65536" | sudo tee -a /etc/sysctl.conf
echo "fs.inotify.max_user_watches = 524288" | sudo tee -a /etc/sysctl.conf

# Apply
sudo sysctl -p
```

### 2. Optimize for Media Server

```bash
# Network tuning
echo "net.core.rmem_max = 134217728" | sudo tee -a /etc/sysctl.conf
echo "net.core.wmem_max = 134217728" | sudo tee -a /etc/sysctl.conf
echo "net.ipv4.tcp_rmem = 4096 87380 67108864" | sudo tee -a /etc/sysctl.conf
echo "net.ipv4.tcp_wmem = 4096 65536 67108864" | sudo tee -a /etc/sysctl.conf

# Apply
sudo sysctl -p
```

### 3. CPU Governor (Ryzen 5 7600X)

```bash
# Install cpupower
sudo dnf install -y kernel-tools

# Set to performance mode
sudo cpupower frequency-set -g performance

# Make permanent
echo "performance" | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor
```

## Deployment

### 1. Deploy Core Services

```bash
cd ~/homelab

# Create network
docker network create homelab

# Deploy Traefik
cd compose/core/traefik
docker compose up -d

# Deploy LLDAP
cd ../lldap
docker compose up -d

# Wait for LLDAP to be ready (30 seconds)
sleep 30

# Deploy Tinyauth
cd ../tinyauth
docker compose up -d
```

### 2. Configure LLDAP

```bash
# Access LLDAP web UI
# https://lldap.fig.systems

# 1. Login with admin credentials from .env
# 2. Create observer user for tinyauth
# 3. Create regular users
```

### 3. Deploy Monitoring

```bash
cd ~/homelab

# Deploy logging stack
cd compose/monitoring/logging
docker compose up -d

# Deploy uptime monitoring
cd ../uptime
docker compose up -d
```

### 4. Deploy Services

See [`README.md`](../../README.md) for complete deployment order.

## Verification

### 1. Check All Services

```bash
# List all running containers
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Check networks
docker network ls

# Check volumes
docker volume ls
```

### 2. Test GPU Access

```bash
# Test in Jellyfin
docker exec jellyfin nvidia-smi

# Test in Ollama
docker exec ollama nvidia-smi

# Test in Immich
docker exec immich-machine-learning nvidia-smi
```

### 3. Test Logging

```bash
# Check Promtail is collecting logs
docker logs promtail | grep "clients configured"

# Access Grafana
# https://logs.fig.systems

# Query logs
# {container="traefik"}
```

### 4. Test SSL

```bash
# Check certificate
curl -vI https://sonarr.fig.systems 2>&1 | grep -i "subject:"

# Should show valid Let's Encrypt certificate
```

## Backup Strategy

### 1. VM Snapshots (Proxmox)

```bash
# On Proxmox host
# Create snapshot before major changes
qm snapshot 100 pre-update-$(date +%Y%m%d)

# List snapshots
qm listsnapshot 100

# Restore snapshot
qm rollback 100 <snapshot-name>
```

### 2. Configuration Backup

```bash
# On VM
cd ~/homelab

# Backup all configs (excludes data directories)
tar czf homelab-config-$(date +%Y%m%d).tar.gz \
  --exclude='*/data' \
  --exclude='*/db' \
  --exclude='*/pgdata' \
  --exclude='*/config' \
  --exclude='*/models' \
  --exclude='*_data' \
  compose/

# Backup to external storage
scp homelab-config-*.tar.gz user@backup-server:/backups/
```

### 3. Automated Backups with Backrest

Backrest service is included and configured. See:
- `compose/services/backrest/`
- Access: https://backup.fig.systems

## Maintenance

### Weekly

```bash
# Update containers
cd ~/homelab
find compose -name "compose.yaml" -type f | while read compose; do
  dir=$(dirname "$compose")
  echo "Updating $dir"
  cd "$dir"
  docker compose pull
  docker compose up -d
  cd ~/homelab
done

# Clean up old images
docker image prune -a -f

# Check disk space
df -h
ncdu /media
```

### Monthly

```bash
# Update AlmaLinux
sudo dnf update -y

# Update NVIDIA drivers (if available)
sudo dnf update nvidia-driver* -y

# Reboot if kernel updated
sudo reboot
```

## Troubleshooting

### Services Won't Start

```bash
# Check SELinux denials
sudo ausearch -m avc -ts recent

# If SELinux is blocking:
sudo setsebool -P container_manage_cgroup on

# Or relabel directories
sudo restorecon -Rv ~/homelab/compose
```

### GPU Not Detected

```bash
# Check GPU is passed through
lspci | grep -i nvidia

# Check drivers loaded
lsmod | grep nvidia

# Reinstall drivers
sudo dnf reinstall nvidia-driver* -y
sudo reboot
```

### Network Issues

```bash
# Check firewall
sudo firewall-cmd --list-all

# Add ports if needed
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --reload

# Check Docker network
docker network inspect homelab
```

### Permission Denied Errors

```bash
# Check ownership
ls -la ~/homelab/compose/*/

# Fix ownership
sudo chown -R $USER:$USER ~/homelab

# Check SELinux context
ls -Z ~/homelab/compose

# Fix SELinux labels
sudo chcon -R -t container_file_t ~/homelab/compose
```

## Performance Monitoring

### System Stats

```bash
# CPU usage
htop

# GPU usage
watch -n 1 nvidia-smi

# Disk I/O
iostat -x 1

# Network
iftop

# Per-container stats
docker stats
```

### Resource Limits

Example container resource limits:

```yaml
# In compose.yaml
deploy:
  resources:
    limits:
      cpus: '2.0'
      memory: 4G
    reservations:
      cpus: '1.0'
      memory: 2G
```

## Security Hardening

### 1. Disable Root SSH

```bash
# Edit /etc/ssh/sshd_config
sudo sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# Restart SSH
sudo systemctl restart sshd
```

### 2. Configure Fail2Ban

```bash
# Install
sudo dnf install -y fail2ban

# Configure
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local

# Edit /etc/fail2ban/jail.local
# [sshd]
# enabled = true
# maxretry = 3
# bantime = 3600

# Start
sudo systemctl enable --now fail2ban
```

### 3. Automatic Updates

```bash
# Install dnf-automatic
sudo dnf install -y dnf-automatic

# Configure /etc/dnf/automatic.conf
# apply_updates = yes

# Enable
sudo systemctl enable --now dnf-automatic.timer
```

## Next Steps

1. ✅ VM created and AlmaLinux installed
2. ✅ Docker and NVIDIA drivers configured
3. ✅ Homelab repository cloned
4. ✅ Network and storage configured
5. ⬜ Deploy core services
6. ⬜ Configure SSO
7. ⬜ Deploy all services
8. ⬜ Configure backups
9. ⬜ Set up monitoring

---

**System ready for deployment!** 🚀