# Terraform Provider for Zone.ee

This Terraform provider allows you to manage resources on Zone.ee, an Estonian hosting provider offering domain registration, web hosting, cloud servers, and DNS services.

## Overview

The Zone.ee Terraform provider enables infrastructure as code management for various Zone.ee resources, including:

- Domains
- Domain nameservers
- DNS records
- DNSSEC settings
- And more...

## Documentation

For detailed information about the Zone.ee API used by this provider, see the official documentation:

https://api.zone.eu/v2

## Usage

### Authentication

The provider requires your Zone.ee username and API key for authentication. [You can generate an API key in your Zone.ee account management panel](https://help.zone.eu/en/kb/zone-api-en/).

Configure authentication using environment variables:

```bash
export ZONE_USERNAME="your-username"
export ZONE_API_KEY="your-api-key"
```

Or via provider configuration:

```hcl
provider "zone-ee" {
    username = "your-username"
    api_key = "your-api-key"
}
```

## Development
This provider is still under development. Contributions are welcome!

### Building the provider for local development
```sh
$ brew install go terraform
$ go build -o terraform-provider-zone-ee
$ mkdir -p ~/.terraform.d/plugins/registry.terraform.io/local/zone-ee/1.0.0/$(go env GOOS)_$(go env GOARCH)
$ cp terraform-provider-zone-ee ~/.terraform.d/plugins/registry.terraform.io/local/zone-ee/1.0.0/$(go env GOOS)_$(go env GOARCH)/
```

## License

MIT License