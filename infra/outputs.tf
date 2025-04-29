output "backend_app_url" {
  value = "https://${var.backend_app_name}.azurewebsites.net"
}

output "frontend_static_site_url" {
  value = azurerm_static_web_app.frontend.default_host_name
}

output "acr_login_server" {
  value = azurerm_container_registry.this.login_server
}
