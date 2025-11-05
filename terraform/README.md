# OpenTofu Infrastructure as Code for Proxmox

This directory contains OpenTofu (Terraform) configurations for managing Proxmox infrastructure.

## What is OpenTofu?

OpenTofu is an open-source fork of Terraform, providing Infrastructure as Code (IaC) capabilities. It allows you to:

- 📝 **Define infrastructure as code** - Version control your infrastructure
- 🔄 **Automate provisioning** - Create VMs/containers with one command
- 🎯 **Consistency** - Same config = same result every time
- 🔍 **Plan changes** - Preview changes before applying
- 🗑️ **Easy cleanup** - Destroy infrastructure when done

## Why OpenTofu over Terraform?

- ✅ **Truly Open Source** - MPL 2.0 license (vs. Terraform's BSL)
- ✅ **Community Driven** - Not controlled by single company
- ✅ **Terraform Compatible** - Drop-in replacement
- ✅ **Active Development** - Regular updates and features

## Quick Start

### 1. Install OpenTofu

**Linux/macOS:**
```bash
# Install via package manager
curl --proto '=https' --tlsv1.2 -fsSL https://get.opentofu.org/install-opentofu.sh | sh

# Or via Homebrew (macOS/Linux)
brew install opentofu
```

**Verify installation:**
```bash
tofu version
```

### 2. Configure Proxmox API

**Create API Token in Proxmox:**
1. Login to Proxmox web UI
2. Datacenter → Permissions → API Tokens
3. Add new token:
   - User: `root@pam`
   - Token ID: `terraform`
   - Privilege Separation: Unchecked (for full access)
4. Save the token ID and secret!

**Set environment variables:**
```bash
export PM_API_URL="https://proxmox.local:8006/api2/json"
export PM_API_TOKEN_ID="root@pam!terraform"
export PM_API_TOKEN_SECRET="your-secret-here"

# Verify SSL (optional, set to false for self-signed certs)
export PM_TLS_INSECURE=true
```

### 3. Choose Your Use Case

We provide examples for common scenarios:

| Example | Description | Best For |
|---------|-------------|----------|
| [single-vm](./proxmox-examples/single-vm/) | Simple Ubuntu VM | Learning, testing |
| [docker-host](./proxmox-examples/docker-host/) | VM for Docker containers | Production homelab |
| [lxc-containers](./proxmox-examples/lxc-containers/) | Lightweight LXC containers | Resource efficiency |
| [multi-node](./proxmox-examples/multi-node/) | Multiple VMs/services | Complex deployments |
| [cloud-init](./proxmox-examples/cloud-init/) | Cloud-init automation | Production VMs |

## Directory Structure

```
terraform/
├── README.md                    # This file
├── proxmox-examples/
│   ├── single-vm/              # Simple VM example
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── terraform.tfvars
│   ├── docker-host/            # Docker host VM
│   ├── lxc-containers/         # LXC container examples
│   ├── multi-node/             # Multiple VM deployment
│   └── cloud-init/             # Cloud-init examples
└── modules/                    # Reusable modules (future)
```

## Basic Workflow

### Initialize

```bash
cd proxmox-examples/single-vm
tofu init
```

### Plan

```bash
tofu plan
```

Preview changes before applying.

### Apply

```bash
tofu apply
```

Review plan and type `yes` to proceed.

### Destroy

```bash
tofu destroy
```

Removes all managed resources.

## Common Commands

```bash
# Initialize and download providers
tofu init

# Validate configuration syntax
tofu validate

# Format code to standard style
tofu fmt

# Preview changes
tofu plan

# Apply changes
tofu apply

# Apply without confirmation (careful!)
tofu apply -auto-approve

# Show current state
tofu show

# List all resources
tofu state list

# Destroy specific resource
tofu destroy -target=proxmox_vm_qemu.vm

# Destroy everything
tofu destroy
```

## Provider Configuration

### Proxmox Provider

```hcl
terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "~> 0.50"
    }
  }
}

provider "proxmox" {
  endpoint = var.pm_api_url
  api_token = "${var.pm_token_id}!${var.pm_token_secret}"
  insecure = true  # For self-signed certs

  ssh {
    agent = true
  }
}
```

## Best Practices

### 1. Use Variables

Don't hardcode values:
```hcl
# Bad
target_node = "pve"

# Good
target_node = var.proxmox_node
```

### 2. Use terraform.tfvars

Store configuration separately:
```hcl
# terraform.tfvars
proxmox_node = "pve"
vm_name = "docker-host"
vm_cores = 4
vm_memory = 8192
```

### 3. Version Control

**Commit:**
- ✅ `*.tf` files
- ✅ `*.tfvars` (if no secrets)
- ✅ `.terraform.lock.hcl`

**DO NOT commit:**
- ❌ `terraform.tfstate`
- ❌ `terraform.tfstate.backup`
- ❌ `.terraform/` directory
- ❌ Secrets/passwords

Use `.gitignore`:
```
.terraform/
*.tfstate
*.tfstate.backup
*.tfvars  # If contains secrets
```

### 4. Use Modules

For reusable components:
```hcl
module "docker_vm" {
  source = "./modules/docker-host"

  vm_name = "docker-01"
  cores   = 4
  memory  = 8192
}
```

### 5. State Management

**Local State (default):**
- Simple, single-user
- State in `terraform.tfstate`

**Remote State (recommended for teams):**
```hcl
terraform {
  backend "s3" {
    bucket = "my-terraform-state"
    key    = "proxmox/terraform.tfstate"
    region = "us-east-1"
  }
}
```

## Example Use Cases

### Homelab Docker Host

Provision a VM optimized for Docker:
- 4-8 CPU cores
- 8-16GB RAM
- 50GB+ disk
- Ubuntu Server 24.04
- Docker pre-installed via cloud-init

See: `proxmox-examples/docker-host/`

### Development Environment

Multiple VMs for testing:
- Web server VM
- Database VM
- Application VM
- All networked together

See: `proxmox-examples/multi-node/`

### LXC Containers

Lightweight containers for services:
- Lower overhead than VMs
- Fast startup
- Resource efficient

See: `proxmox-examples/lxc-containers/`

## Proxmox Provider Resources

### Virtual Machines (QEMU)

```hcl
resource "proxmox_vm_qemu" "vm" {
  name        = "my-vm"
  target_node = "pve"

  clone      = "ubuntu-cloud-template"  # Template to clone
  cores      = 2
  memory     = 2048

  disk {
    size    = "20G"
    storage = "local-lvm"
  }

  network {
    model  = "virtio"
    bridge = "vmbr0"
  }
}
```

### LXC Containers

```hcl
resource "proxmox_lxc" "container" {
  hostname    = "my-container"
  target_node = "pve"

  ostemplate  = "local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.gz"
  cores       = 1
  memory      = 512

  rootfs {
    storage = "local-lvm"
    size    = "8G"
  }

  network {
    name   = "eth0"
    bridge = "vmbr0"
    ip     = "dhcp"
  }
}
```

### Cloud-Init

```hcl
resource "proxmox_vm_qemu" "cloudinit_vm" {
  # ... basic config ...

  ciuser     = "ubuntu"
  cipassword = var.vm_password
  sshkeys    = file("~/.ssh/id_rsa.pub")

  ipconfig0  = "ip=dhcp"
}
```

## Troubleshooting

### SSL Certificate Errors

```bash
export PM_TLS_INSECURE=true
```

Or add to provider:
```hcl
provider "proxmox" {
  insecure = true
}
```

### API Permission Errors

Ensure API token has necessary permissions:
```bash
# In Proxmox shell
pveum acl modify / -token 'root@pam!terraform' -role Administrator
```

### VM Clone Errors

Ensure template exists:
```bash
# List VMs
qm list

# Check template flag
qm config 9000
```

### Timeout Errors

Increase timeout:
```hcl
resource "proxmox_vm_qemu" "vm" {
  # ...
  timeout_create = "30m"
  timeout_clone  = "30m"
}
```

## Migration from Terraform

OpenTofu is a drop-in replacement:

```bash
# Rename binary
alias tofu=terraform

# Or replace commands
terraform → tofu
```

State files are compatible - no conversion needed!

## Advanced Topics

### Custom Cloud Images

1. Download cloud image
2. Create VM template
3. Use cloud-init for customization

See: `proxmox-examples/cloud-init/`

### Network Configuration

```hcl
# VLAN tagging
network {
  model = "virtio"
  bridge = "vmbr0"
  tag = 100  # VLAN 100
}

# Multiple NICs
network {
  model = "virtio"
  bridge = "vmbr0"
}
network {
  model = "virtio"
  bridge = "vmbr1"
}
```

### Storage Options

```hcl
# Local LVM
disk {
  storage = "local-lvm"
  size    = "50G"
  type    = "scsi"
}

# NFS/CIFS
disk {
  storage = "nfs-storage"
  size    = "100G"
}

# Multiple disks
disk {
  slot = 0
  size = "50G"
  storage = "local-lvm"
}
disk {
  slot = 1
  size = "100G"
  storage = "data"
}
```

## Recommended Resources

### Providers

- **[bpg/proxmox](https://registry.terraform.io/providers/bpg/proxmox)** - Most feature-complete (recommended)
- **[Telmate/proxmox](https://registry.terraform.io/providers/Telmate/proxmox)** - Legacy, still works

### Learning

- [OpenTofu Docs](https://opentofu.org/docs/)
- [Proxmox Provider Docs](https://registry.terraform.io/providers/bpg/proxmox/latest/docs)
- [Terraform/OpenTofu Tutorial](https://developer.hashicorp.com/terraform/tutorials)

### Tools

- **[tflint](https://github.com/terraform-linters/tflint)** - Linting
- **[terraform-docs](https://github.com/terraform-docs/terraform-docs)** - Generate docs
- **[infracost](https://www.infracost.io/)** - Cost estimation
- **[terragrunt](https://terragrunt.gruntwork.io/)** - Wrapper for DRY configs

## Next Steps

1. **Start Simple:** Try `proxmox-examples/single-vm/`
2. **Learn Basics:** Get familiar with plan/apply/destroy
3. **Expand:** Try docker-host or multi-node
4. **Customize:** Adapt examples to your needs
5. **Automate:** Integrate with CI/CD

## Getting Help

- Check example READMEs in each directory
- Review Proxmox provider docs
- OpenTofu community Discord
- Ask in r/Proxmox or r/selfhosted

Happy Infrastructure as Code! 🚀
