package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DomainNameserver represents a domain nameserver
type DomainNameserver struct {
	Hostname string `json:"hostname"`
	IPv4     string `json:"ipv4,omitempty"`
	IPv6     string `json:"ipv6,omitempty"`
}

func resourceDomainNameserver() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainNameserverCreate,
		ReadContext:   resourceDomainNameserverRead,
		UpdateContext: resourceDomainNameserverUpdate,
		DeleteContext: resourceDomainNameserverDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Nameserver hostname",
			},
			"ipv4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 address for the nameserver (glue record)",
			},
			"ipv6": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv6 address for the nameserver (glue record)",
			},
		},
	}
}

func resourceDomainNameserverCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	domainName := d.Get("domain").(string)
	hostname := d.Get("hostname").(string)

	nameserver := DomainNameserver{
		Hostname: hostname,
	}

	if v, ok := d.GetOk("ipv4"); ok {
		nameserver.IPv4 = v.(string)
	}

	if v, ok := d.GetOk("ipv6"); ok {
		nameserver.IPv6 = v.(string)
	}

	// POST /domain/{service_name}/nameserver
	_, err := client.doRequest("POST", fmt.Sprintf("/domain/%s/nameserver", domainName), []DomainNameserver{nameserver})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", domainName, hostname))

	return resourceDomainNameserverRead(ctx, d, m)
}

func resourceDomainNameserverRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	domainName := d.Get("domain").(string)
	hostname := d.Get("hostname").(string)

	// GET /domain/{service_name}/nameserver/{hostname}
	resp, err := client.doRequest("GET", fmt.Sprintf("/domain/%s/nameserver/%s", domainName, hostname), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var nameservers []DomainNameserver
	if err := parseResponse(resp, &nameservers); err != nil {
		return diag.FromErr(err)
	}

	if len(nameservers) == 0 {
		d.SetId("")
		return diags
	}

	nameserver := nameservers[0]
	d.Set("hostname", nameserver.Hostname)
	d.Set("ipv4", nameserver.IPv4)
	d.Set("ipv6", nameserver.IPv6)

	return diags
}

func resourceDomainNameserverUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	domainName := d.Get("domain").(string)
	hostname := d.Get("hostname").(string)

	nameserver := DomainNameserver{
		Hostname: hostname,
	}

	if v, ok := d.GetOk("ipv4"); ok {
		nameserver.IPv4 = v.(string)
	}

	if v, ok := d.GetOk("ipv6"); ok {
		nameserver.IPv6 = v.(string)
	}

	// PUT /domain/{service_name}/nameserver/{hostname}
	_, err := client.doRequest("PUT", fmt.Sprintf("/domain/%s/nameserver/%s", domainName, hostname), nameserver)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDomainNameserverRead(ctx, d, m)
}

func resourceDomainNameserverDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	domainName := d.Get("domain").(string)
	hostname := d.Get("hostname").(string)

	// DELETE /domain/{service_name}/nameserver/{hostname}
	_, err := client.doRequest("DELETE", fmt.Sprintf("/domain/%s/nameserver/%s", domainName, hostname), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
