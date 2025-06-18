package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVmNames() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmNamesRead,
		Schema: map[string]*schema.Schema{
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment to filter VM names (e.g., dev, prod, preprod, devtest)",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Azure region/location to further filter VM names (e.g., uksouth)",
			},
			"vm_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of VM names returned by the API.",
			},
			"statuses": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A map of VM name to status, if details are available.",
			},
		},
	}
}

func dataSourceVmNamesRead(d *schema.ResourceData, m interface{}) error {
	url := m.(string)
	environment := d.Get("environment").(string)
	location, hasLocation := d.GetOk("location")

	apiUrl := fmt.Sprintf("%s?environment=%s", url, environment)
	if hasLocation {
		apiUrl += "&location=" + location.(string)
	}
	resp, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	vmNames := strings.Split(strings.TrimSpace(string(body)), ",")
	var filtered []string
	statuses := make(map[string]string)
	for _, n := range vmNames {
		name := strings.TrimSpace(n)
		if name == "" {
			continue
		}
		filtered = append(filtered, name)
		// Fetch status for each VM name
		detailsUrl := fmt.Sprintf("%s?environment=%s&rowkey=%s&details=true", url, environment, name)
		detailsResp, err := http.Get(detailsUrl)
		if err != nil {
			statuses[name] = "error"
			continue
		}
		defer detailsResp.Body.Close()
		if detailsResp.StatusCode == 200 {
			var result map[string]interface{}
			b, _ := ioutil.ReadAll(detailsResp.Body)
			if err := json.Unmarshal(b, &result); err == nil {
				if status, ok := result["Status"].(string); ok {
					statuses[name] = status
				}
			}
		}
	}
	d.SetId(environment + "-" + fmt.Sprintf("%d", len(filtered)))
	d.Set("vm_names", filtered)
	d.Set("statuses", statuses)
	return nil
}
