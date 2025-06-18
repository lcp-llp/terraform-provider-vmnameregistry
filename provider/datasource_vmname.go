package provider

import (
	"fmt"
	"io/ioutil"
	"net/http"

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
		},
	}
}

func dataSourceVmNameRead(d *schema.ResourceData, m interface{}) error {
	url := m.(string)
	vmName := d.Get("vm_name").(string)
	apiUrl := url + "?rowkey=" + vmName
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
	d.SetId(vmName)
	return nil
}
