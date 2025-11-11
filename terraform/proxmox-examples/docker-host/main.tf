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

  # BIOS type - OVMF required for GPU passthrough
  bios = var.enable_gpu_passthrough ? "ovmf" : "seabios"

  # Machine type - q35 required for GPU passthrough
  machine = var.enable_gpu_passthrough ? "q35" : "pc"

  # CPU configuration
  cpu {
    cores = var.vm_cores
    type  = "host"  # Use host CPU type for best performance
  }

  # Memory configuration
  memory {
    dedicated = var.vm_memory
  }

  # EFI disk (required for OVMF BIOS when GPU passthrough is enabled)
  dynamic "efi_disk" {
    for_each = var.enable_gpu_passthrough ? [1] : []
    content {
      datastore_id = var.storage
      type         = "4m"
    }
  }

  # GPU passthrough configuration
  dynamic "hostpci" {
    for_each = var.enable_gpu_passthrough ? [1] : []
    content {
      device  = "hostpci0"
      mapping = var.gpu_pci_id
      pcie    = true
      rombar  = true
      xvga    = false
    }
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
  datastore_id = var.snippets_storage
  node_name    = var.proxmox_node

  source_raw {
    data = var.vm_os_type == "almalinux" ? local.cloud_init_almalinux : local.cloud_init_ubuntu

    file_name = "cloud-init-docker-${var.vm_name}.yaml"
  }
}

# Cloud-init configuration for Ubuntu
locals {
  cloud_init_ubuntu = <<-EOF
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
      ${var.mount_media_directories ? "- nfs-common" : ""}

    # Docker installation and NFS mount setup
    runcmd:
      # Install Docker
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
      - mkdir -p ${var.media_mount_path}/{audiobooks,books,comics,complete,downloads,homemovies,incomplete,movies,music,photos,tv}

      ${var.mount_media_directories ? "# Mount NFS shares from Proxmox host" : ""}
      ${var.mount_media_directories ? "- systemctl enable nfs-client.target" : ""}
      ${var.mount_media_directories ? "- systemctl start nfs-client.target" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/audiobooks ${var.media_mount_path}/audiobooks" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/books ${var.media_mount_path}/books" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/comics ${var.media_mount_path}/comics" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/complete ${var.media_mount_path}/complete" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/downloads ${var.media_mount_path}/downloads" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/homemovies ${var.media_mount_path}/homemovies" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/incomplete ${var.media_mount_path}/incomplete" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/movies ${var.media_mount_path}/movies" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/music ${var.media_mount_path}/music" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/photos ${var.media_mount_path}/photos" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/tv ${var.media_mount_path}/tv" : ""}

      - chown -R ${var.vm_username}:${var.vm_username} ${var.media_mount_path}
      - chmod -R 755 ${var.media_mount_path}

      ${var.clone_homelab_repo ? "- su - ${var.vm_username} -c 'cd ~ && git clone https://github.com/${var.github_username}/homelab.git'" : ""}

    ${var.mount_media_directories ? "# Make NFS mounts persistent" : ""}
    ${var.mount_media_directories ? "write_files:" : ""}
    ${var.mount_media_directories ? "  - path: /etc/fstab" : ""}
    ${var.mount_media_directories ? "    append: true" : ""}
    ${var.mount_media_directories ? "    content: |" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/audiobooks ${var.media_mount_path}/audiobooks nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/books ${var.media_mount_path}/books nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/comics ${var.media_mount_path}/comics nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/complete ${var.media_mount_path}/complete nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/downloads ${var.media_mount_path}/downloads nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/homemovies ${var.media_mount_path}/homemovies nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/incomplete ${var.media_mount_path}/incomplete nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/movies ${var.media_mount_path}/movies nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/music ${var.media_mount_path}/music nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/photos ${var.media_mount_path}/photos nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/tv ${var.media_mount_path}/tv nfs defaults 0 0" : ""}

    # Set timezone
    timezone: ${var.vm_timezone}

    # Reboot after setup
    power_state:
      mode: reboot
      condition: true
  EOF

  cloud_init_almalinux = <<-EOF
    #cloud-config
    hostname: ${var.vm_name}
    manage_etc_hosts: true

    # Install Docker and dependencies
    package_update: true
    package_upgrade: true

    packages:
      - curl
      - ca-certificates
      - git
      - vim
      - htop
      - net-tools
      ${var.mount_media_directories ? "- nfs-utils" : ""}

    # Docker installation and NFS mount setup
    runcmd:
      # Install Docker
      - dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
      - dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
      - systemctl enable docker
      - systemctl start docker
      - usermod -aG docker ${var.vm_username}
      - docker network create homelab || true

      # Create media directories
      - mkdir -p ${var.media_mount_path}/{audiobooks,books,comics,complete,downloads,homemovies,incomplete,movies,music,photos,tv}

      ${var.mount_media_directories ? "# Mount NFS shares from Proxmox host" : ""}
      ${var.mount_media_directories ? "- systemctl enable nfs-client.target" : ""}
      ${var.mount_media_directories ? "- systemctl start nfs-client.target" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/audiobooks ${var.media_mount_path}/audiobooks" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/books ${var.media_mount_path}/books" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/comics ${var.media_mount_path}/comics" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/complete ${var.media_mount_path}/complete" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/downloads ${var.media_mount_path}/downloads" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/homemovies ${var.media_mount_path}/homemovies" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/incomplete ${var.media_mount_path}/incomplete" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/movies ${var.media_mount_path}/movies" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/music ${var.media_mount_path}/music" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/photos ${var.media_mount_path}/photos" : ""}
      ${var.mount_media_directories ? "- mount -t nfs ${var.proxmox_host_ip}:${var.media_source_path}/tv ${var.media_mount_path}/tv" : ""}

      - chown -R ${var.vm_username}:${var.vm_username} ${var.media_mount_path}
      - chmod -R 755 ${var.media_mount_path}

      ${var.clone_homelab_repo ? "- su - ${var.vm_username} -c 'cd ~ && git clone https://github.com/${var.github_username}/homelab.git'" : ""}

    ${var.mount_media_directories ? "# Make NFS mounts persistent" : ""}
    ${var.mount_media_directories ? "write_files:" : ""}
    ${var.mount_media_directories ? "  - path: /etc/fstab" : ""}
    ${var.mount_media_directories ? "    append: true" : ""}
    ${var.mount_media_directories ? "    content: |" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/audiobooks ${var.media_mount_path}/audiobooks nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/books ${var.media_mount_path}/books nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/comics ${var.media_mount_path}/comics nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/complete ${var.media_mount_path}/complete nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/downloads ${var.media_mount_path}/downloads nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/homemovies ${var.media_mount_path}/homemovies nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/incomplete ${var.media_mount_path}/incomplete nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/movies ${var.media_mount_path}/movies nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/music ${var.media_mount_path}/music nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/photos ${var.media_mount_path}/photos nfs defaults 0 0" : ""}
    ${var.mount_media_directories ? "      ${var.proxmox_host_ip}:${var.media_source_path}/tv ${var.media_mount_path}/tv nfs defaults 0 0" : ""}

    # Set timezone
    timezone: ${var.vm_timezone}

    # Reboot after setup
    power_state:
      mode: reboot
      condition: true
  EOF
}
