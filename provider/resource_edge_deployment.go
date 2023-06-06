package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		},
	}
}

func createOrUpdateEdgeDeployment(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)

	d.SetId(d.Get("site_short_name").(string))

	return pm.Client.CreateOrUpdateEdgeDeployment(pm.Corp, d.Get("site_short_name").(string))
}

func readEdgeDeployment(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteEdgeDeployment(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)

	return pm.Client.DeleteEdgeDeployment(pm.Corp, d.Get("site_short_name").(string))
}
