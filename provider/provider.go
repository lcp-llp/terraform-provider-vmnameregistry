package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base URL for the registry API.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vmnameregistry_vmname": resourceVmName(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"vmnameregistry_vmname":  dataSourceVmName(),
			"vmnameregistry_vmnames": dataSourceVmNames(),
		},
		ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return d.Get("url").(string), nil
		},
	}
}
