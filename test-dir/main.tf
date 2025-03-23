terraform {
  required_providers {
    zonee = {
      source  = "local/zonee"
      version = "1.0.0"
    }
  }
}

provider "zonee" {
  # Your provider configuration here
}

resource "zonee_domain" "midwork_ee" {
  name      = "midwork.ee"
  autorenew = true
}