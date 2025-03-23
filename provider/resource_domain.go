package provider

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Domain represents a domain resource
type Domain struct {
	Name        string    `json:"name"`
	Expires     time.Time `json:"expires"`
	Autorenew   bool      `json:"autorenew"`
	IsRenewable bool      `json:"is_renewable"`
	IsDelegated bool      `json:"is_delegated"`
	Status      string    `json:"status"`
}

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name",
			},
			"autorenew": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to automatically renew the domain",
			},
			"expires": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain expiration date",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain status",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_ = m.(*Client)

	// Domain registration requires creating an order
	// This would use the /order/domain endpoint
	// The actual implementation depends on the API structure

	// For now, we'll just set the ID to the domain name since we don't see the full API
	d.SetId(d.Get("name").(string))

	return resourceDomainRead(ctx, d, m)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	domainName := d.Id()

	// GET /domain/{service_name}
	resp, err := client.doRequest("GET", fmt.Sprintf("/domain/%s", domainName), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var domains []Domain
	if err := parseResponse(resp, &domains); err != nil {
		return diag.FromErr(err)
	}

	if len(domains) == 0 {
		d.SetId("")
		return diags
	}

	domain := domains[0]
	d.Set("name", domain.Name)
	d.Set("autorenew", domain.Autorenew)
	d.Set("expires", domain.Expires.Format(time.RFC3339))
	d.Set("status", domain.Status)

	return diags
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	domainName := d.Id()

	// Only certain properties can be updated
	if d.HasChange("autorenew") {
		// PUT /domain/{service_name}
		updateBody := map[string]interface{}{
			"autorenew": d.Get("autorenew").(bool),
		}

		resp, err := client.doRequest("PUT", fmt.Sprintf("/domain/%s", domainName), updateBody)
		if err != nil {
			return diag.FromErr(err)
		}

		// Check response
		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			return diag.Errorf("Failed to update domain: %s", string(body))
		}
	}

	return resourceDomainRead(ctx, d, m)
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Domains typically can't be deleted via API, only allowed to expire
	// For now, we'll just remove it from Terraform state

	var diags diag.Diagnostics
	return diags
}
