landingzoneFake = {
  backend_type        = "azurerm"
  global_settings_key = "caf_foundations"
  level               = "level2"
  key                 = "networking_hub"
  tfstates = {
    caf_foundations = {
      level   = "lower"
      tfstate = "caf_foundations.tfstate"
    }
    launchpad = {
      level   = "lower"
      tfstate = "caf_foundations.tfstate"
    }
  }
}
