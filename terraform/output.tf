output "load_balancer_public_ip" {
  description = "Public IP address of load balancer"
  value = tolist(tolist(yandex_lb_network_load_balancer.wp_lb.listener).0.external_address_spec).0.address
}

output "database_host_fqdn" {
  description = "DB hostname"
  value = local.dbhosts
}

output "database_cluster_fqdn" {
  description = "DB cluster FQDN"
  value = local.dbclusterfqdn
}

output "database_user" {
  description = "DB user"
  value = local.dbuser
}

output "database_password" {
  description = "DB password"
  value = local.dbpassword
  sensitive = true
}

output "database_name" {
  description = "DB name"
  value = local.dbname
}

output "vm_linux_public_ip_address" {
  description = "Virtual machine IP"
  value = yandex_compute_instance.wp-app-a[0].network_interface[0].nat_ip_address
}