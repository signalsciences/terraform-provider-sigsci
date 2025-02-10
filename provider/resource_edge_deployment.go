package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEdgeDeployment() *schema.Resource {
	return &schema.Resource{
		Create:   createOrUpdateEdgeDeployment,
		Read:     readEdgeDeployment,
		Update:   createOrUpdateEdgeDeployment,
		Delete:   deleteEdgeDeployment,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
			},
			"authorized_services": {
				Type:        schema.TypeList,
				Description: "List of Compute services. This field is only required if you are linking Compute services to the Next-Gen WAF.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func createOrUpdateEdgeDeployment(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)

	// Initialize an empty authorizedServices slice
	var authorizedServices []string
	siteName := d.Get("site_short_name").(string)
	// Check if "authorized_services" exists before accessing it
	if v, ok := d.GetOk("authorized_services"); ok {
		authorizedServicesInterface := v.([]interface{})

		// Convert []interface{} to []string
		for _, v := range authorizedServicesInterface {
			authorizedServices = append(authorizedServices, v.(string))
		}
	}
	log.Printf("[DEBUG] Creating/Updating Edge Deployment for site: %s, authorized_services: %+v", siteName, authorizedServices)

	// Call the updated client function, passing an empty slice if "authorized_services" wasn't set
	err := pm.Client.CreateOrUpdateEdgeDeployment(pm.Corp, d.Get("site_short_name").(string), authorizedServices)

	if err != nil {
		log.Printf("[ERROR] Failed to create/update Edge Deployment: %v", err)
		return err
	}
	log.Printf("[DEBUG] Successfully updated Edge Deployment for site: %s", siteName)
	d.SetId(d.Get("site_short_name").(string))

	return nil
}

func readEdgeDeployment(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteEdgeDeployment(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)

	return pm.Client.DeleteEdgeDeployment(pm.Corp, d.Get("site_short_name").(string))
}
