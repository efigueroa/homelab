output "vm_id" {
  description = "VM ID"
  value       = proxmox_virtual_environment_vm.docker_host.vm_id
}

output "vm_name" {
  description = "VM name"
  value       = proxmox_virtual_environment_vm.docker_host.name
}

output "vm_ipv4_address" {
  description = "VM IPv4 address"
  value       = try(proxmox_virtual_environment_vm.docker_host.ipv4_addresses[1][0], "DHCP - check Proxmox UI")
}

output "vm_mac_address" {
  description = "VM MAC address"
  value       = proxmox_virtual_environment_vm.docker_host.mac_addresses[0]
}

output "ssh_command" {
  description = "SSH command to connect to VM"
  value       = "ssh ${var.vm_username}@${try(proxmox_virtual_environment_vm.docker_host.ipv4_addresses[1][0], "DHCP-ADDRESS")}"
}

output "docker_status_command" {
  description = "Command to check Docker status"
  value       = "ssh ${var.vm_username}@${try(proxmox_virtual_environment_vm.docker_host.ipv4_addresses[1][0], "DHCP-ADDRESS")} 'docker ps'"
}
