package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flags": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNSSEC key flags (usually 256 for ZSK or 257 for KSK)",
						},
						"algorithm": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNSSEC algorithm number",
						},
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNSSEC public key",
						},
					},
				},
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

	var dnsZones []DNSZone
	if err := parseResponse(resp, &dnsZones); err != nil {
		return diag.FromErr(err)
	}

	if len(dnsZones) == 0 {
		return diag.Errorf("No DNS zone found for domain %s", domainName)
	}

	dnsZone := dnsZones[0]
	d.SetId(domainName)
	d.Set("domain", dnsZone.Name)
	d.Set("enabled", dnsZone.HasDNSSEC)

	if dnsZone.HasDNSSEC && len(dnsZone.DNSSECKeys) > 0 {
		keys := make([]map[string]interface{}, 0, len(dnsZone.DNSSECKeys))
		for _, k := range dnsZone.DNSSECKeys {
			key := map[string]interface{}{
				"flags":      k.Flags,
				"algorithm":  k.Algorithm,
				"public_key": k.PublicKey,
			}
			keys = append(keys, key)
		}
		d.Set("keys", keys)
	}

	return diags
}
