# Docker Host VM with OpenTofu

This configuration creates a VM optimized for running Docker containers in your homelab with support for GPU passthrough and NFS media mounts.

## What This Creates

- ✅ Ubuntu or AlmaLinux VM (from cloud template)
- ✅ Docker & Docker Compose installed
- ✅ Homelab network created
- ✅ /media directories structure
- ✅ SSH key authentication
- ✅ Automatic updates enabled
- ✅ Optional GPU passthrough (NVIDIA GTX 1070)
- ✅ Optional NFS mounts from Proxmox host

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

**Or create AlmaLinux 9.6 Cloud Template:**

```bash
# SSH to Proxmox server
ssh root@proxmox.local

# Download AlmaLinux cloud image
wget https://repo.almalinux.org/almalinux/9/cloud/x86_64/images/AlmaLinux-9-GenericCloud-latest.x86_64.qcow2

# Create VM
qm create 9001 --name almalinux-cloud-template --memory 2048 --net0 virtio,bridge=vmbr0

# Import disk
qm importdisk 9001 AlmaLinux-9-GenericCloud-latest.x86_64.qcow2 local-lvm

# Attach disk
qm set 9001 --scsihw virtio-scsi-pci --scsi0 local-lvm:vm-9001-disk-0

# Add cloud-init drive
qm set 9001 --ide2 local-lvm:cloudinit

# Set boot disk
qm set 9001 --boot c --bootdisk scsi0

# Add serial console
qm set 9001 --serial0 socket --vga serial0

# Convert to template
qm template 9001

# Cleanup
rm AlmaLinux-9-GenericCloud-latest.x86_64.qcow2
```

### 2. (Optional) Enable GPU Passthrough

**For NVIDIA GTX 1070 on AMD Ryzen CPU:**

```bash
# On Proxmox host, edit GRUB config
nano /etc/default/grub

# Add to GRUB_CMDLINE_LINUX_DEFAULT:
GRUB_CMDLINE_LINUX_DEFAULT="quiet amd_iommu=on iommu=pt"

# Update GRUB
update-grub

# Load required kernel modules
nano /etc/modules

# Add these lines:
vfio
vfio_iommu_type1
vfio_pci
vfio_virqfd

# Blacklist NVIDIA drivers on host
nano /etc/modprobe.d/blacklist.conf

# Add:
blacklist nouveau
blacklist nvidia
blacklist nvidiafb
blacklist nvidia_drm

# Update initramfs
update-initramfs -u -k all

# Reboot Proxmox host
reboot

# After reboot, verify IOMMU is enabled:
dmesg | grep -e DMAR -e IOMMU

# Find GPU PCI ID:
lspci | grep -i nvidia
# Output example: 01:00.0 VGA compatible controller: NVIDIA Corporation GP104 [GeForce GTX 1070]
# Use: 0000:01:00 (note the format)
```

### 3. (Optional) Configure NFS Server on Proxmox

**Export media directories from Proxmox host:**

```bash
# On Proxmox host
# Install NFS server
apt update
apt install nfs-kernel-server -y

# Create /etc/exports entry
nano /etc/exports

# Add (replace 192.168.1.0/24 with your network):
/data/media/audiobooks 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/books 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/comics 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/complete 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/downloads 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/homemovies 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/incomplete 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/movies 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/music 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/photos 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)
/data/media/tv 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)

# Export NFS shares
exportfs -ra

# Enable and start NFS server
systemctl enable nfs-server
systemctl start nfs-server

# Verify exports
showmount -e localhost
```

### 4. Create API Token

In Proxmox UI:
1. Datacenter → Permissions → API Tokens
2. Add → User: `root@pam`, Token ID: `terraform`
3. Uncheck "Privilege Separation"
4. Save the secret!

### 5. Install OpenTofu

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
- `vm_os_type` - Choose "ubuntu", "almalinux", or "debian"
- `template_vm_id` - Use 9001 for AlmaLinux template
- `enable_gpu_passthrough` - Set to true for GPU support
- `gpu_pci_id` - Your GPU PCI ID (find with `lspci`)
- `mount_media_directories` - Set to true for NFS mounts
- `proxmox_host_ip` - IP for NFS server (Proxmox host)
- `media_source_path` - Path on Proxmox host (default: /data/media)

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

### GPU Passthrough Configuration

**Enable NVIDIA GTX 1070 for Jellyfin, Ollama, Immich:**

```hcl
# Must complete Proxmox host GPU passthrough setup first
enable_gpu_passthrough = true
gpu_pci_id = "0000:01:00"  # Find with: lspci | grep -i nvidia

# Use AlmaLinux for better GPU support
vm_os_type = "almalinux"
template_vm_id = 9001

# Allocate sufficient resources
vm_cores = 8
vm_memory = 24576  # 24GB
```

**Verify GPU in VM after deployment:**

```bash
ssh ubuntu@<VM-IP>

# Install NVIDIA drivers (AlmaLinux)
sudo dnf install -y epel-release
sudo dnf config-manager --add-repo https://developer.download.nvidia.com/compute/cuda/repos/rhel9/x86_64/cuda-rhel9.repo
sudo dnf install -y nvidia-driver nvidia-container-toolkit

# Verify
nvidia-smi
docker run --rm --gpus all nvidia/cuda:12.3.0-base-ubuntu22.04 nvidia-smi
```

### NFS Media Mounts Configuration

**Mount Proxmox host media directories to VM:**

```hcl
# Enable NFS mounts from Proxmox host
mount_media_directories = true

# Proxmox host IP (not API URL)
proxmox_host_ip = "192.168.1.100"

# Source path on Proxmox host
media_source_path = "/data/media"

# Mount point in VM
media_mount_path = "/media"
```

**After deployment, verify mounts:**

```bash
ssh ubuntu@<VM-IP>

# Check mounts
df -h | grep /media
ls -la /media

# Expected directories:
# /media/audiobooks, /media/books, /media/comics,
# /media/complete, /media/downloads, /media/homemovies,
# /media/incomplete, /media/movies, /media/music,
# /media/photos, /media/tv
```

### Operating System Selection

**AlmaLinux 9.6 (Recommended for GPU):**

```hcl
vm_os_type = "almalinux"
template_vm_id = 9001
vm_username = "almalinux"  # Default AlmaLinux user
```

**Ubuntu 22.04 LTS:**

```hcl
vm_os_type = "ubuntu"
template_vm_id = 9000
vm_username = "ubuntu"
```

**Key differences:**
- AlmaLinux: Better RHEL ecosystem, SELinux, dnf package manager
- Ubuntu: Wider community support, apt package manager
- Both support Docker, GPU passthrough, and NFS mounts

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

# Check media directories and NFS mounts
ls -la /media
df -h | grep /media

# If GPU passthrough is enabled
nvidia-smi
lspci | grep -i nvidia

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

### GPU Not Detected in VM

**Verify IOMMU is enabled:**
```bash
# On Proxmox host
dmesg | grep -e DMAR -e IOMMU
# Should show: "IOMMU enabled"
```

**Check GPU is available:**
```bash
# On Proxmox host
lspci | grep -i nvidia
lspci -n -s 01:00

# Verify it's not being used by host
lsmod | grep nvidia
# Should be empty (blacklisted)
```

**In VM, install drivers:**
```bash
# AlmaLinux
sudo dnf install -y epel-release
sudo dnf config-manager --add-repo https://developer.download.nvidia.com/compute/cuda/repos/rhel9/x86_64/cuda-rhel9.repo
sudo dnf install -y nvidia-driver

# Ubuntu
sudo apt install -y nvidia-driver-535
sudo reboot

# Verify
nvidia-smi
```

### NFS Mounts Not Working

**On Proxmox host, verify NFS server:**
```bash
systemctl status nfs-server
showmount -e localhost
# Should list all /data/media/* exports
```

**In VM, test manual mount:**
```bash
# Install NFS client if missing
sudo apt install nfs-common  # Ubuntu
sudo dnf install nfs-utils   # AlmaLinux

# Test mount
sudo mount -t nfs 192.168.1.100:/data/media/movies /mnt
ls /mnt
sudo umount /mnt
```

**Check /etc/fstab in VM:**
```bash
cat /etc/fstab | grep nfs
# Should show all media directory mounts
```

**Firewall issues:**
```bash
# On Proxmox host, allow NFS
ufw allow from 192.168.1.0/24 to any port nfs
# Or disable firewall temporarily to test:
systemctl stop ufw
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
