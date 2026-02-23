package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceEdgeDeploymentService() *schema.Resource {
	return &schema.Resource{
		Create:   createEdgeDeploymentService,
		Read:     readEdgeDeploymentService,
		Update:   updateEdgeDeploymentService,
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
				Description: "activate Fastly service version after clone. Possible values are true or false. Defaults to true.",
				Optional:    true,
				Default:     true,
			},
			"custom_client_ip": {
				Type:        schema.TypeBool,
				Description: "enable to prevent Fastly-Client-IP from being overwritten by the NGWAF. Intended for advanced use cases. Defaults to false.",
				Optional:    true,
				Default:     false,
			},
			"percent_enabled": {
				Type:        schema.TypeInt,
				Description: "percentage of traffic to send to NGWAF@Edge. Possible values are integers values 0 to 100. Defaults to 0.",
				Optional:    true,
				Default:     0,
			},
			"sync_id": {
				Type:        schema.TypeInt,
				Description: "A numeric identifier used to trigger a synchronization of the NGWAF VCL module. Incrementing this value forces a non-destructive update of the Edge WAF VCL version on the Fastly service without requiring a resource replacement. Set to 0 to disable this trigger; any value greater than 0 will initiate a sync when the value is changed. Defaults to 0.",
				Optional:    true,
				Default:     0,
			},
		},
	}
}

func createEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)

	activateVersion := d.Get("activate_version").(bool)
	custom_client_ip := d.Get("custom_client_ip").(bool)
	percent_enabled := d.Get("percent_enabled").(int)
	sync_id := d.Get("sync_id").(int)

	err := pm.Client.CreateOrUpdateEdgeDeploymentService(pm.Corp, d.Get("site_short_name").(string), d.Get("fastly_sid").(string), sigsci.CreateOrUpdateEdgeDeploymentServiceBody{
		ActivateVersion: &activateVersion,
		CustomClientIP:  &custom_client_ip,
		PercentEnabled:  &percent_enabled,
	})
	if err != nil {
		return err
	}

	d.SetId(d.Get("fastly_sid").(string))
	d.Set("sync_id", sync_id)

	return nil
}

func updateEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	fastlySID := d.Get("fastly_sid").(string)
	siteName := d.Get("site_short_name").(string)
	activateVersion := d.Get("activate_version").(bool)
	custom_client_ip := d.Get("custom_client_ip").(bool)
	percent_enabled := d.Get("percent_enabled").(int)
	sync_id := d.Get("sync_id").(int)

	// Handle site remapping (existing logic)
	if d.HasChange("site_short_name") {
		oldSite, newSite := d.GetChange("site_short_name")
		siteName = newSite.(string)

		// First detach site from service
		if err := pm.Client.DetachEdgeDeploymentService(pm.Corp, oldSite.(string), fastlySID); err != nil {
			return err
		}
	}

	// Only call CreateOrUpdateEdgeDeploymentService when:
	// 1. sync_id > 0 (user wants to trigger sync), OR
	// 2. Functional fields have changed (activate_version, custom_client_ip, percent_enabled)
	syncTriggerActive := sync_id > 0
	functionalFieldsChanged := d.HasChange("activate_version") || d.HasChange("custom_client_ip") || d.HasChange("percent_enabled")

	if syncTriggerActive || functionalFieldsChanged || d.HasChange("site_short_name") {
		err := pm.Client.CreateOrUpdateEdgeDeploymentService(pm.Corp, siteName, fastlySID, sigsci.CreateOrUpdateEdgeDeploymentServiceBody{
			ActivateVersion: &activateVersion,
			CustomClientIP:  &custom_client_ip,
			PercentEnabled:  &percent_enabled,
		})
		if err != nil {
			return err
		}
	}

	d.SetId(fastlySID)
	d.Set("sync_id", sync_id)

	return nil
}

func readEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	return nil
}

func detachEdgeDeploymentService(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)

	return pm.Client.DetachEdgeDeploymentService(pm.Corp, d.Get("site_short_name").(string), d.Get("fastly_sid").(string))
}
