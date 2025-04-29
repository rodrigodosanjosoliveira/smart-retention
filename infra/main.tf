provider "azurerm" {
  features {}
  subscription_id = "481ebe44-22ba-4099-b3b6-752fd9336192" # <- Substituir aqui pelo ID da sua assinatura
}

resource "azurerm_resource_group" "this" {
  name     = var.resource_group_name
  location = var.resource_group_location
}

resource "azurerm_container_registry" "this" {
  name                = var.acr_name
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  sku                 = "Basic"
  admin_enabled       = true
}

resource "azurerm_service_plan" "this" {
  name                = var.app_service_plan_name
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  os_type             = "Linux"
  sku_name            = "B1"
}


resource "azurerm_linux_web_app" "backend" {
  name                = var.backend_app_name
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  service_plan_id     = azurerm_service_plan.this.id

  site_config {
    application_stack {
      docker_image_name        = "${azurerm_container_registry.this.login_server}/smart-retention-backend:latest"
      docker_registry_url      = "https://${azurerm_container_registry.this.login_server}"
      docker_registry_username = azurerm_container_registry.this.admin_username
      docker_registry_password = azurerm_container_registry.this.admin_password
    }
    always_on        = true
    app_command_line = "./app"
  }

  app_settings = {
    WEBSITES_PORT = "8080"
  }
}

resource "azurerm_static_web_app" "frontend" {
  name                = var.frontend_app_name
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  sku_tier            = "Free"
  sku_size            = "Free"
}
