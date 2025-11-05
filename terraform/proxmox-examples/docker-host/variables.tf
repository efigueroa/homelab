variable "pm_api_url" {
  description = "Proxmox API URL"
  type        = string
  default     = "https://proxmox.local:8006"
}

variable "pm_api_token_id" {
  description = "Proxmox API token ID (format: user@realm!tokenid)"
  type        = string
  default     = "root@pam!terraform"
}

variable "pm_api_token_secret" {
  description = "Proxmox API token secret"
  type        = string
  sensitive   = true
}

variable "pm_tls_insecure" {
  description = "Disable TLS verification for self-signed certificates"
  type        = bool
  default     = true
}

variable "proxmox_node" {
  description = "Proxmox node name"
  type        = string
  default     = "pve"
}

variable "vm_name" {
  description = "VM name"
  type        = string
  default     = "docker-host"
}

variable "template_vm_id" {
  description = "Template VM ID to clone from"
  type        = number
  default     = 9000
}

variable "vm_cores" {
  description = "Number of CPU cores"
  type        = number
  default     = 4
}

variable "vm_memory" {
  description = "Memory in MB"
  type        = number
  default     = 8192
}

variable "disk_size" {
  description = "Disk size (e.g., 50G, 100G)"
  type        = string
  default     = "50"
}

variable "storage" {
  description = "Storage pool name"
  type        = string
  default     = "local-lvm"
}

variable "network_bridge" {
  description = "Network bridge"
  type        = string
  default     = "vmbr0"
}

variable "vm_ip_address" {
  description = "Static IP address or 'dhcp'"
  type        = string
  default     = "dhcp"
}

variable "vm_ip_netmask" {
  description = "Network netmask (CIDR notation, e.g., 24)"
  type        = number
  default     = 24
}

variable "vm_gateway" {
  description = "Network gateway"
  type        = string
  default     = "192.168.1.1"
}

variable "vm_username" {
  description = "VM username"
  type        = string
  default     = "ubuntu"
}

variable "vm_password" {
  description = "VM user password"
  type        = string
  sensitive   = true
}

variable "vm_ssh_keys" {
  description = "List of SSH public keys"
  type        = list(string)
  default     = []
}

variable "vm_timezone" {
  description = "VM timezone"
  type        = string
  default     = "America/Los_Angeles"
}

variable "clone_homelab_repo" {
  description = "Clone homelab repository on first boot"
  type        = bool
  default     = false
}

variable "github_username" {
  description = "GitHub username for homelab repo"
  type        = string
  default     = "efigueroa"
}
