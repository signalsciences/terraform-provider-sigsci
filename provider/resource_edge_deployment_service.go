package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceEdgeDeploymentService() *schema.Resource {
	return &schema.Resource{
		Create:   createOrUpdateEdgeDeploymentService,
		Read:     readEdgeDeploymentService,
		Update:   createOrUpdateEdgeDeploymentService,
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

			"activate_version": {
				Type:        schema.TypeBool,
				Description: "activate Fastly service version after clone. Possible values are true or false",
				Optional:    true,
			},

			"percent_enabled": {
				Type:        schema.TypeInt,
				Description: "percentage of traffic to send to NGWAF@Edge. Possible values are integers values 0 to 100",
				Optional:    true,
			},
		},
	}
}

func createOrUpdateEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)

	d.SetId(d.Get("fastly_sid").(string))

	activateVersion := d.Get("activate_version").(bool)
	return pm.Client.CreateOrUpdateEdgeDeploymentService(pm.Corp, d.Get("site_short_name").(string), d.Get("fastly_sid").(string), sigsci.CreateOrUpdateEdgeDeploymentServiceBody{
		ActivateVersion:         &activateVersion,
		PercentEnabled:          d.Get("percent_enabled").(int),
	})
}

func readEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	return nil
}

func detachEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)

	return pm.Client.DetachEdgeDeploymentService(pm.Corp, d.Get("site_short_name").(string), d.Get("fastly_sid").(string))
}
