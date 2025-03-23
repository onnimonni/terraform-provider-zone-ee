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
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Api key for zone.ee API authentication",
				DefaultFunc: schema.EnvDefaultFunc("ZONE_API_KEY", nil),
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.zone.eu/v2",
				Description: "API endpoint URL for zone.ee",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zone-ee_domain":             resourceDomain(),
			"zone-ee_domain_nameservers": resourceDomainNameservers(),
			"zone-ee_domain_dnssec":      resourceDomainDNSSEC(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure configures the provider and creates a client
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	api_key := d.Get("api_key").(string)
	apiURL := d.Get("api_url").(string)

	var diags diag.Diagnostics

	client := NewClient(apiURL, username, api_key)
	return client, diags
}
