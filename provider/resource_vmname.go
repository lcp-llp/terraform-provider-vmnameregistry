package provider

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVmName() *schema.Resource {
	return &schema.Resource{
		Create: resourceVmNameCreate,
		Read:   resourceVmNameRead,
		Update: resourceVmNameUpdate,
		Delete: resourceVmNameDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVmNameImport,
		},
		Schema: map[string]*schema.Schema{
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment for the VM name (e.g., dev, prod, preprod, devtest)",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Azure region/location for the VM name (e.g., uksouth)",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Deployed",
				Description: "The status of the VM name (Deployed, Reserved, Available)",
			},
			"business_unit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The business unit for the VM name",
			},
			"vm_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated VM name returned by the API.",
			},
		},
	}
}

func resourceVmNameCreate(d *schema.ResourceData, m interface{}) error {
	url := m.(string)
	environment := d.Get("environment").(string)
	location := d.Get("location").(string)
	status := d.Get("status").(string)
	businessUnit := d.Get("business_unit").(string)

	apiUrl := fmt.Sprintf("%s?environment=%s&location=%s&status=%s&businessunit=%s", url, environment, location, status, businessUnit)
	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
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
	vmName := string(body)
	if idx := bytes.IndexByte(body, ','); idx > 0 {
		vmName = string(body[:idx])
	}
	d.SetId(vmName)
	d.Set("vm_name", vmName)
	return resourceVmNameRead(d, m)
}

func resourceVmNameRead(d *schema.ResourceData, m interface{}) error {
	url := m.(string)
	vmName := d.Id()
	environment := d.Get("environment").(string)
	apiUrl := fmt.Sprintf("%s?environment=%s&rowkey=%s", url, environment, vmName)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	d.Set("vm_name", string(body))
	return nil
}

func resourceVmNameUpdate(d *schema.ResourceData, m interface{}) error {
	url := m.(string)
	vmName := d.Id()
	status := d.Get("status").(string)
	environment := d.Get("environment").(string)
	businessUnit := d.Get("business_unit").(string)
	apiUrl := fmt.Sprintf("%s?environment=%s&rowkey=%s&status=%s&businessunit=%s", url, environment, vmName, status, businessUnit)
	req, err := http.NewRequest("PUT", apiUrl, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}
	return resourceVmNameRead(d, m)
}

func resourceVmNameDelete(d *schema.ResourceData, m interface{}) error {
	url := m.(string)
	vmName := d.Id()
	environment := d.Get("environment").(string)
	apiUrl := fmt.Sprintf("%s?environment=%s&rowkey=%s", url, environment, vmName)
	req, err := http.NewRequest("DELETE", apiUrl, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}
	d.SetId("")
	return nil
}

func resourceVmNameImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// The import ID should be in the format: environment/vm_name
	// Example: dev/lcpdevuks-0001
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID format. Expected: environment/vm_name (e.g., dev/lcpdevuks-0001)")
	}

	environment := parts[0]
	vmName := parts[1]

	// Set the VM name as the ID
	d.SetId(vmName)
	d.Set("vm_name", vmName)
	d.Set("environment", environment)

	// Extract location from VM name
	location := extractLocationFromVmName(vmName)
	if location != "" {
		d.Set("location", location)
	}

	// Call the Read function to populate the rest of the state from the API
	// This will fetch status and business_unit from the API
	err := resourceVmNameRead(d, m)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
