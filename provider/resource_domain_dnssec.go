package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DNSSecKey represents a DNSSEC key
type DNSSecKey struct {
	Flags     int    `json:"flags"`
	Algorithm int    `json:"algorithm"`
	PublicKey string `json:"public_key"`
}

// DNSZone represents a DNS zone with DNSSEC settings
type DNSZone struct {
	Name       string      `json:"name"`
	DNSSECKeys []DNSSecKey `json:"dnssec_keys"`
	HasDNSSEC  bool        `json:"has_dnssec"`
}

func resourceDomainDNSSEC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainDNSSECCreate,
		ReadContext:   resourceDomainDNSSECRead,
		UpdateContext: resourceDomainDNSSECUpdate,
		DeleteContext: resourceDomainDNSSECDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether DNSSEC is enabled",
			},
			"key": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flags": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "DNSSEC key flags (usually 256 for ZSK or 257 for KSK)",
							ValidateFunc: validation.IntInSlice([]int{256, 257}),
						},
						"algorithm": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "DNSSEC algorithm number",
						},
						"public_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DNSSEC public key",
						},
					},
				},
			},
		},
	}
}

func resourceDomainDNSSECCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_ = m.(*Client)
	_ = d.Get("domain").(string)

	// For DNSSEC, we'll use the PUT method as it's an update operation to the DNS zone
	return resourceDomainDNSSECUpdate(ctx, d, m)
}

func resourceDomainDNSSECRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		d.SetId("")
		return diags
	}

	dnsZone := dnsZones[0]
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
		d.Set("key", keys)
	}

	return diags
}

func resourceDomainDNSSECUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	domainName := d.Get("domain").(string)
	enabled := d.Get("enabled").(bool)

	dnsZone := DNSZone{
		Name:      domainName,
		HasDNSSEC: enabled,
	}

	if enabled {
		keySet := d.Get("key").(*schema.Set)
		if keySet.Len() > 0 {
			keys := make([]DNSSecKey, 0, keySet.Len())

			for _, v := range keySet.List() {
				key := v.(map[string]interface{})
				dnsSecKey := DNSSecKey{
					Flags:     key["flags"].(int),
					Algorithm: key["algorithm"].(int),
					PublicKey: key["public_key"].(string),
				}
				keys = append(keys, dnsSecKey)
			}

			dnsZone.DNSSECKeys = keys
		}
	}

	// PUT /dns/{service_name}
	_, err := client.doRequest("PUT", fmt.Sprintf("/dns/%s", domainName), dnsZone)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domainName)

	return resourceDomainDNSSECRead(ctx, d, m)
}

func resourceDomainDNSSECDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Deleting DNSSEC is just disabling it
	d.Set("enabled", false)
	return resourceDomainDNSSECUpdate(ctx, d, m)
}
