# Docker Host VM with OpenTofu

This configuration creates a VM optimized for running Docker containers in your homelab.

## What This Creates

- ✅ Ubuntu VM (from cloud template)
- ✅ Docker & Docker Compose installed
- ✅ Homelab network created
- ✅ /media directories structure
- ✅ SSH key authentication
- ✅ Automatic updates enabled

## Prerequisites

### 1. Create Ubuntu Cloud Template

First, create a cloud-init enabled template in Proxmox:

```bash
# SSH to Proxmox server
ssh root@proxmox.local

# Download Ubuntu cloud image
wget https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img

# Create VM
qm create 9000 --name ubuntu-cloud-template --memory 2048 --net0 virtio,bridge=vmbr0

# Import disk
qm importdisk 9000 jammy-server-cloudimg-amd64.img local-lvm

# Attach disk
qm set 9000 --scsihw virtio-scsi-pci --scsi0 local-lvm:vm-9000-disk-0

# Add cloud-init drive
qm set 9000 --ide2 local-lvm:cloudinit

# Set boot disk
qm set 9000 --boot c --bootdisk scsi0

# Add serial console
qm set 9000 --serial0 socket --vga serial0

# Convert to template
qm template 9000

# Cleanup
rm jammy-server-cloudimg-amd64.img
```

### 2. Create API Token

In Proxmox UI:
1. Datacenter → Permissions → API Tokens
2. Add → User: `root@pam`, Token ID: `terraform`
3. Uncheck "Privilege Separation"
4. Save the secret!

### 3. Install OpenTofu

```bash
# Linux/macOS
curl --proto '=https' --tlsv1.2 -fsSL https://get.opentofu.org/install-opentofu.sh | sh

# Verify
tofu version
```

## Quick Start

### 1. Configure Variables

```bash
cd terraform/proxmox-examples/docker-host

# Copy example config
cp terraform.tfvars.example terraform.tfvars

# Edit with your values
nano terraform.tfvars
```

**Required changes:**
- `pm_api_token_secret` - Your Proxmox API secret
- `vm_ssh_keys` - Your SSH public key
- `vm_password` - Set a secure password

**Optional changes:**
- `vm_name` - Change VM name
- `vm_cores` / `vm_memory` - Adjust resources
- `vm_ip_address` - Set static IP (or keep DHCP)

### 2. Initialize

```bash
tofu init
```

### 3. Plan

```bash
tofu plan
```

Review what will be created.

### 4. Apply

```bash
tofu apply
```

Type `yes` to confirm.

### 5. Connect

```bash
# Get SSH command from output
tofu output ssh_command

# Or manually
ssh ubuntu@<VM-IP>

# Verify Docker
docker --version
docker ps
docker network ls | grep homelab
```

## Configuration Options

### Resource Sizing

**Light workload (1-5 containers):**
```hcl
vm_cores  = 2
vm_memory = 4096
disk_size = "30"
```

**Medium workload (5-15 containers):**
```hcl
vm_cores  = 4
vm_memory = 8192
disk_size = "50"
```

**Heavy workload (15+ containers):**
```hcl
vm_cores  = 8
vm_memory = 16384
disk_size = "100"
```

### Network Configuration

**DHCP (easiest):**
```hcl
vm_ip_address = "dhcp"
```

**Static IP:**
```hcl
vm_ip_address = "192.168.1.100"
vm_ip_netmask = 24
vm_gateway    = "192.168.1.1"
```

### Multiple SSH Keys

```hcl
vm_ssh_keys = [
  "ssh-rsa AAAAB3... user1@laptop",
  "ssh-rsa AAAAB3... user2@desktop"
]
```

## Post-Deployment

### Deploy Homelab Services

```bash
# SSH to VM
ssh ubuntu@<VM-IP>

# Clone homelab repo (if not auto-cloned)
git clone https://github.com/efigueroa/homelab.git
cd homelab

# Deploy services
cd compose/core/traefik
docker compose up -d

cd ../lldap
docker compose up -d

# Continue with other services...
```

### Verify Setup

```bash
# Check Docker
docker --version
docker compose version

# Check network
docker network ls | grep homelab

# Check media directories
ls -la /media

# Check system resources
htop
df -h
```

## Managing the VM

### View State

```bash
tofu show
tofu state list
```

### Update VM

1. Edit `terraform.tfvars`:
```hcl
vm_cores  = 8  # Increase from 4
vm_memory = 16384  # Increase from 8192
```

2. Apply changes:
```bash
tofu plan
tofu apply
```

**Note:** Some changes require VM restart.

### Destroy VM

```bash
# Backup any data first!
tofu destroy
```

Type `yes` to confirm deletion.

## Troubleshooting

### Template Not Found

Error: `template with ID 9000 not found`

**Solution:** Create cloud template (see Prerequisites)

### API Permission Error

Error: `permission denied`

**Solution:** Check API token permissions:
```bash
pveum acl modify / -token 'root@pam!terraform' -role Administrator
```

### Cloud-Init Not Working

**Check cloud-init status:**
```bash
ssh ubuntu@<VM-IP>
sudo cloud-init status
sudo cat /var/log/cloud-init-output.log
```

### Docker Not Installed

**Manual installation:**
```bash
ssh ubuntu@<VM-IP>
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker ubuntu
```

### VM Won't Start

**Check Proxmox logs:**
```bash
# On Proxmox server
qm status <VM-ID>
tail -f /var/log/pve/tasks/active
```

## Advanced Usage

### Multiple VMs

Create `docker-host-02.tfvars`:
```hcl
vm_name = "docker-host-02"
vm_ip_address = "192.168.1.101"
```

Deploy:
```bash
tofu apply -var-file="docker-host-02.tfvars"
```

### Custom Cloud-Init

Edit `main.tf` to add custom cloud-init sections:
```yaml
users:
  - name: myuser
    groups: sudo, docker
    shell: /bin/bash
    sudo: ALL=(ALL) NOPASSWD:ALL

packages:
  - zsh
  - tmux
  - neovim

runcmd:
  - sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

### Attach Additional Disk

Add to `main.tf`:
```hcl
disk {
  datastore_id = var.storage
  size         = 200
  interface    = "scsi1"
}
```

Mount in cloud-init:
```yaml
mounts:
  - ["/dev/sdb", "/mnt/data", "ext4", "defaults", "0", "0"]
```

## Cost Analysis

**Resource Usage:**
- 4 cores, 8GB RAM, 50GB disk
- Running 24/7

**Homelab Cost:** $0 (uses existing hardware)

**If in cloud (comparison):**
- AWS: ~$50-100/month
- DigitalOcean: ~$40/month
- Linode: ~$40/month

**Homelab ROI:** Pays for itself in ~2-3 months!

## Security Hardening

### Enable Firewall

Add to cloud-init:
```yaml
runcmd:
  - ufw default deny incoming
  - ufw default allow outgoing
  - ufw allow ssh
  - ufw allow 80/tcp
  - ufw allow 443/tcp
  - ufw --force enable
```

### Disable Password Authentication

After SSH key setup:
```yaml
ssh_pwauth: false
```

### Automatic Updates

Already enabled in cloud-init. Verify:
```bash
sudo systemctl status unattended-upgrades
```

## Next Steps

1. ✅ Deploy core services (Traefik, LLDAP, Tinyauth)
2. ✅ Configure SSL certificates
3. ✅ Deploy media services
4. ✅ Set up backups (Restic)
5. ✅ Add monitoring (Prometheus/Grafana)

## Resources

- [OpenTofu Docs](https://opentofu.org/docs/)
- [Proxmox Provider](https://registry.terraform.io/providers/bpg/proxmox/latest/docs)
- [Cloud-Init Docs](https://cloudinit.readthedocs.io/)
- [Docker Docs](https://docs.docker.com/)
