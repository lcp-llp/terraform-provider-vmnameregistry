package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVmName() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmNameRead,
		Schema: map[string]*schema.Schema{
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
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location/region extracted from the VM name.",
			},
		},
	}
}

func dataSourceVmNameRead(d *schema.ResourceData, m interface{}) error {
	url := m.(string)
	vmName := d.Get("vm_name").(string)

	// Call API with details=true to get status and business unit
	apiUrl := url + "?rowkey=" + vmName + "&details=true"
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

	// Parse the JSON response to extract details
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		// Set computed attributes from API response
		if status, ok := result["Status"].(string); ok {
			d.Set("status", status)
		}
		if businessUnit, ok := result["BusinessUnit"].(string); ok {
			d.Set("business_unit", businessUnit)
		}
	}

	// Extract location from VM name pattern (e.g., lcpdevuks-0001 -> uks)
	// VM names follow pattern: lcp{env}{location}-{number}
	location := extractLocationFromVmName(vmName)
	if location != "" {
		d.Set("location", location)
	}

	d.SetId(vmName)
	return nil
}

// Helper function to extract location from VM name
func extractLocationFromVmName(vmName string) string {
	// VM names follow pattern: lcp{env}{location}-{number}
	// e.g., lcpdevuks-0001, lcpprduks-0002, etc.
	parts := strings.Split(vmName, "-")
	if len(parts) < 2 {
		return ""
	}

	// Remove the "lcp" prefix and environment code to get location
	prefix := parts[0]    // e.g., "lcpdevuks"
	if len(prefix) <= 3 { // Should be at least "lcp" + something
		return ""
	}

	// Remove "lcp" prefix
	withoutLcp := prefix[3:] // e.g., "devuks"

	// Common environment codes to remove
	envCodes := []string{"dev", "prd", "ppd", "dvt"}
	for _, envCode := range envCodes {
		if strings.HasPrefix(withoutLcp, envCode) {
			return withoutLcp[len(envCode):] // Return the location part
		}
	}

	return ""
}
