terraform {
  required_providers {
    zone-ee = {
      source  = "local/zone-ee"
      version = "0.1.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

variable "cloudflare_account_id" {
	description = "Cloudflare account id. Zones currently need this to be set"
	type = string
	sensitive = true 
}

variable "domain" {
	description = "The domain name to manage"
	type = string
}

resource "cloudflare_zone" "zone" {
	account_id	= var.cloudflare_account_id
	zone 		    = var.domain
	plan       	= "free"
	type       	= "full"
	jump_start 	= false # Don't scan the existing DNS records
}

resource "zone-ee_domain_nameservers" "name_servers" {
  domain = var.domain
  nameservers = cloudflare_zone.zone.name_servers
}

data "zone-ee_domain_dnssec" "domain_dnssec" {
  domain = var.domain
}

# Throw an error if DNSSEC is not disabled
# Cloudflare can't activate the name servers when DNSSEC is still activated
resource "terraform_data" "dnssec_check" {
  lifecycle {
    precondition {
      condition     = !data.zone-ee_domain_dnssec.domain_dnssec.enabled
      error_message = "🚨 ${var.domain} has DNSSEC enabled. Disable it here https://my.zone.eu/dashboard/en/${var.domain}/domain/dnssec?domain=${var.domain}"
    }
  }
}
