variable "tfstate_container_name" {}
variable "tfstate_key" {}
variable "tfstate_resource_group_name" {}

variable "landingzone" {
  default = {
    backend_type        = "azurerm"
    global_settings_key = "launchpad"

  
