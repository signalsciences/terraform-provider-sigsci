package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEdgeDeploymentServiceBackend() *schema.Resource {
	return &schema.Resource{
		Create:   updateEdgeDeploymentServiceBackend,
		Read:     readEdgeDeploymentServiceBackend,
		Update:   updateEdgeDeploymentServiceBackend,
		Delete:   deleteEdgeDeploymentServiceBackend,
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

			"fastly_service_vcl_active_version": {
				Type:        schema.TypeInt,
				Description: "Currently active version of the VCL service. Required to run to sync the VCL services origins to the NGWAF service",
				Required:    true,
			},
		},
	}
}

func updateEdgeDeploymentServiceBackend(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)

	d.SetId(d.Get("fastly_sid").(string))

	return pm.Client.UpdateEdgeDeploymentBackends(pm.Corp, d.Get("site_short_name").(string), d.Get("fastly_sid").(string))

}

func readEdgeDeploymentServiceBackend(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteEdgeDeploymentServiceBackend(d *schema.ResourceData, m interface{}) error {
	return nil
}
