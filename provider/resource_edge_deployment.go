package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
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

	authorizedServices := []string{}
	if v, ok := d.GetOk("authorized_services"); ok {
		for _, serviceID := range v.([]interface{}) {
			authorizedServices = append(authorizedServices, serviceID.(string))
		}
	}

	err := pm.Client.CreateOrUpdateEdgeDeployment(pm.Corp, d.Get("site_short_name").(string), sigsci.CreateOrUpdateEdgeDeploymentBody{
		AuthorizedServices: &authorizedServices,
	})
	if err != nil {
		return err
	}

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
