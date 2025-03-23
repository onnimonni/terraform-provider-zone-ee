package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DNSZoneResponse represents the actual API response structure
type DNSZoneResponse struct {
	ResourceURL   string `json:"resource_url"`
	Identificator string `json:"identificator"`
	Active        bool   `json:"active"`
	IPv6          bool   `json:"ipv6"`
	Domain        bool   `json:"domain"`
	DNSSEC        bool   `json:"dnssec"`
}

func dataSourceDomainDNSSEC() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainDNSSECRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether DNSSEC is enabled",
			},
		},
	}
}

func dataSourceDomainDNSSECRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	domainName := d.Get("domain").(string)

	// GET /dns/{service_name}
	resp, err := client.doRequest("GET", fmt.Sprintf("/dns/%s", domainName), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var dnsZones []DNSZoneResponse
	if err := parseResponse(resp, &dnsZones); err != nil {
		return diag.FromErr(err)
	}

	if len(dnsZones) == 0 {
		return diag.Errorf("No DNS zone found for domain %s", domainName)
	}

	dnsZone := dnsZones[0]
	d.SetId(domainName)
	d.Set("domain", dnsZone.Identificator)
	d.Set("enabled", dnsZone.DNSSEC)

	return diags
}
