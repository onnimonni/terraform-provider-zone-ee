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
}

func resourceDomainNameservers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainNameserversCreate,
		ReadContext:   resourceDomainNameserversRead,
		UpdateContext: resourceDomainNameserversUpdate,
		DeleteContext: resourceDomainNameserversDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name",
			},
			"nameservers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of nameserver hostnames",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDomainNameserversCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	domainName := d.Get("domain").(string)
	nameserversRaw := d.Get("nameservers").([]interface{})

	// Convert nameservers to proper type
	var nameservers []DomainNameserver
	for _, ns := range nameserversRaw {
		nameserver := DomainNameserver{
			Hostname: ns.(string),
		}
		nameservers = append(nameservers, nameserver)
	}

	// First, delete all existing nameservers
	// GET existing nameservers
	resp, err := client.doRequest("GET", fmt.Sprintf("/domain/%s/nameserver", domainName), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var existingNameservers []DomainNameserver
	if err := parseResponse(resp, &existingNameservers); err != nil {
		return diag.FromErr(err)
	}

	// DELETE each existing nameserver
	for _, ns := range existingNameservers {
		_, err := client.doRequest("DELETE", fmt.Sprintf("/domain/%s/nameserver/%s", domainName, ns.Hostname), nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Then create new nameservers
	// POST /domain/{service_name}/nameserver
	_, err = client.doRequest("POST", fmt.Sprintf("/domain/%s/nameserver", domainName), nameservers)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domainName)

	return resourceDomainNameserversRead(ctx, d, m)
}

func resourceDomainNameserversRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	domainName := d.Get("domain").(string)

	// GET /domain/{service_name}/nameserver
	resp, err := client.doRequest("GET", fmt.Sprintf("/domain/%s/nameserver", domainName), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var nameservers []DomainNameserver
	if err := parseResponse(resp, &nameservers); err != nil {
		return diag.FromErr(err)
	}

	// Extract just the hostnames
	nsHostnames := make([]string, 0, len(nameservers))
	for _, ns := range nameservers {
		nsHostnames = append(nsHostnames, ns.Hostname)
	}

	if err := d.Set("nameservers", nsHostnames); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDomainNameserversUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// For updates, we'll use the same approach as create - delete all and recreate
	return resourceDomainNameserversCreate(ctx, d, m)
}

func resourceDomainNameserversDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	domainName := d.Get("domain").(string)

	// GET existing nameservers
	resp, err := client.doRequest("GET", fmt.Sprintf("/domain/%s/nameserver", domainName), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var nameservers []DomainNameserver
	if err := parseResponse(resp, &nameservers); err != nil {
		return diag.FromErr(err)
	}

	// DELETE each nameserver
	for _, ns := range nameservers {
		_, err := client.doRequest("DELETE", fmt.Sprintf("/domain/%s/nameserver/%s", domainName, ns.Hostname), nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return diags
}
