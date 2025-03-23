terraform {
  required_providers {
    zone-ee = {
      source  = "local/zone-ee"
      version = "1.0.0"
    }
  }
}

provider "zone-ee" {
  # Your provider configuration here
}

resource "zone-ee_domain" "midwork_ee" {
  name      = "midwork.ee"
  autorenew = true
}

resource "zone-ee_domain_nameservers" "name_servers" {
  domain = zone-ee_domain.midwork_ee.name

  nameservers = [
    "houston.ns.cloudflare.com",
    "marissa.ns.cloudflare.com"
  ]
}