package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVmName() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmNameRead,
		Schema: map[string]*schema.Schema{
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment for the VM name (e.g., dev, prod, preprod, devtest)",
			},
			"vm_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VM name to look up.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the VM name.",
			},
			"business_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The business unit of the VM name.",
			},
			"cloud_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud platform of the VM name (aws, azure, gcp, onprem).",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location/region of the VM name.",
			},
		},
	}
}

func dataSourceVmNameRead(d *schema.ResourceData, m interface{}) error {
	url := m.(string)
	vmName := d.Get("vm_name").(string)
	environment := d.Get("environment").(string)

	// Call API with details=true to get all fields
	apiUrl := fmt.Sprintf("%s?environment=%s&rowkey=%s&details=true", url, environment, vmName)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return fmt.Errorf("VM name not found")
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if status, ok := result["Status"].(string); ok {
			d.Set("status", status)
		}
		if businessUnit, ok := result["BusinessUnit"].(string); ok {
			d.Set("business_unit", businessUnit)
		}
		if cloudPlatform, ok := result["CloudPlatform"].(string); ok {
			d.Set("cloud_platform", cloudPlatform)
		}
		// Use Location from API response; fall back to pattern extraction
		if location, ok := result["Location"].(string); ok && location != "" {
			d.Set("location", location)
		} else {
			location := extractLocationFromVmName(vmName)
			if location != "" {
				d.Set("location", location)
			}
		}
	}

	d.SetId(vmName)
	return nil
}
