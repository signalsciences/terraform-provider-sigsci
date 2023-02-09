package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEdgeDeploymentService() *schema.Resource {
	return &schema.Resource{
		Create:   createOrUpdateEdgeDeploymentService,
		Read:     readEdgeDeploymentService,
		Update:   updateEdgeDeploymentBackends,
		Delete:   detachEdgeDeploymentService,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
			},

			"fastly_sid": {
				Type:        schema.TypeString,
				Description: "Fastly service ID",
				Required:    true,
			},
		},
	}
}

func createOrUpdateEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	d.SetId(d.Get("site_short_name").(string))

	return sc.CreateOrUpdateEdgeDeploymentService(pm.Corp, d.Get("site_short_name").(string), d.Get("fastly_sid").(string))
}

func readEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	return nil
}

func updateEdgeDeploymentBackends(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	return sc.CreateOrUpdateEdgeDeploymentService(pm.Corp, d.Get("site_short_name").(string), d.Get("fastly_sid").(string))
}

func detachEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	return sc.DetachEdgeDeploymentService(pm.Corp, d.Get("site_short_name").(string), d.Get("fastly_sid").(string))
}
