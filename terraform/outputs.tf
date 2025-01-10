output "service_endpoint" {
  value = format("http://%s", module.alb.dns_name)
}