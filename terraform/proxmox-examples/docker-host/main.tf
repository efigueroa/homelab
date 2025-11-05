terraform {
  required_version = ">= 1.6"

  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "~> 0.50"
    }
  }
}

provider "proxmox" {
  endpoint = var.pm_api_url

  api_token = var.pm_api_token_secret != "" ? "${var.pm_api_token_id}=${var.pm_api_token_secret}" : null

  # For self-signed certificates
  insecure = var.pm_tls_insecure

  ssh {
    agent = true
  }
}

resource "proxmox_virtual_environment_vm" "docker_host" {
  name        = var.vm_name
  description = "Docker host for homelab services - Managed by OpenTofu"
  node_name   = var.proxmox_node

  # Clone from template (must exist in Proxmox)
  clone {
    vm_id = var.template_vm_id
    full  = true
  }

  # CPU configuration
  cpu {
    cores = var.vm_cores
    type  = "host"  # Use host CPU type for best performance
  }

  # Memory configuration
  memory {
    dedicated = var.vm_memory
  }

  # Network interface
  network_device {
    bridge = var.network_bridge
    model  = "virtio"
  }

  # Disk configuration
  disk {
    datastore_id = var.storage
    size         = var.disk_size
    interface    = "scsi0"
    discard      = "on"  # Enable TRIM for SSDs
    iothread     = true
  }

  # Cloud-init configuration
  initialization {
    ip_config {
      ipv4 {
        address = var.vm_ip_address == "dhcp" ? "dhcp" : "${var.vm_ip_address}/${var.vm_ip_netmask}"
        gateway = var.vm_gateway
      }
    }

    user_account {
      username = var.vm_username
      keys     = var.vm_ssh_keys
      password = var.vm_password
    }

    user_data_file_id = proxmox_virtual_environment_file.cloud_init_user_data.id
  }

  # Start VM on boot
  on_boot = true

  # Tags for organization
  tags = ["terraform", "docker", "homelab"]
}

# Cloud-init user data for Docker installation
resource "proxmox_virtual_environment_file" "cloud_init_user_data" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = var.proxmox_node

  source_raw {
    data = <<-EOF
      #cloud-config
      hostname: ${var.vm_name}
      manage_etc_hosts: true

      # Install Docker and dependencies
      package_update: true
      package_upgrade: true

      packages:
        - apt-transport-https
        - ca-certificates
        - curl
        - gnupg
        - lsb-release
        - git
        - vim
        - htop
        - net-tools

      # Add Docker's official GPG key and repository
      runcmd:
        - mkdir -p /etc/apt/keyrings
        - curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        - chmod a+r /etc/apt/keyrings/docker.gpg
        - echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
        - apt-get update
        - apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        - systemctl enable docker
        - systemctl start docker
        - usermod -aG docker ${var.vm_username}
        - docker network create homelab || true

      # Create media directories
      write_files:
        - path: /usr/local/bin/setup-media-dirs
          permissions: '0755'
          content: |
            #!/bin/bash
            mkdir -p /media/{audiobooks,books,comics,complete,downloads,homemovies,incomplete,movies,music,photos,tv}
            chown -R ${var.vm_username}:${var.vm_username} /media
            chmod -R 755 /media

      # Run setup script
      runcmd:
        - /usr/local/bin/setup-media-dirs

      # Optional: Clone homelab repo
      ${var.clone_homelab_repo ? "- su - ${var.vm_username} -c 'cd ~ && git clone https://github.com/${var.github_username}/homelab.git'" : "# Homelab repo cloning disabled"}

      # Set timezone
      timezone: ${var.vm_timezone}

      # Reboot after setup
      power_state:
        mode: reboot
        condition: true
    EOF

    file_name = "cloud-init-docker-${var.vm_name}.yaml"
  }
}
