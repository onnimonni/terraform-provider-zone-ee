package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   false,
				Description: "Username for zone.ee API authentication",
				DefaultFunc: schema.EnvDefaultFunc("ZONE_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password for zone.ee API authentication",
				DefaultFunc: schema.EnvDefaultFunc("ZONE_PASSWORD", nil),
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.zone.ee/v2",
				Description: "API endpoint URL for zone.ee",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zonee_domain":            resourceDomain(),
			"zonee_domain_nameserver": resourceDomainNameserver(),
			"zonee_domain_dnssec":     resourceDomainDNSSEC(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure configures the provider and creates a client
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	apiURL := d.Get("api_url").(string)

	var diags diag.Diagnostics

	client := NewClient(apiURL, username, password)
	return client, diags
}
