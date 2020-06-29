package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"log"
)

func resourceSiteIntegration() *schema.Resource {
	return &schema.Resource{
		Create:   resourceSiteIntegrationCreate,
		Update:   resourceSiteIntegrationUpdate,
		Read:     resourceSiteIntegrationRead,
		Delete:   resourceSiteIntegrationDelete,
		Importer: &siteImporter,
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "One of (mailingList, slack, datadog, generic, pagerduty, microsoftTeams, jira, opsgenie, victorops, pivotaltracker)",
				Required:    true,
				ForceNew:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "Integration URL",
				Required:    true,
			},
			"events": {
				Type:        schema.TypeSet,
				Description: "Array of event types. Visit https://docs.signalsciences.net/integrations to find out which events the service you are connecting allows.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name",
				Computed:    true,
			},
		},
	}
}

func resourceSiteIntegrationCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	integrations, err := sc.AddIntegration(pm.Corp, d.Get("site_short_name").(string), sigsci.IntegrationBody{
		Type:   d.Get("type").(string),
		URL:    d.Get("url").(string),
		Events: expandStringArray(d.Get("events").(*schema.Set)),
	})
	if err != nil {
		return err
	}

	d.SetId(integrations[len(integrations)-1].ID)
	return resourceSiteIntegrationRead(d, m)
}

func resourceSiteIntegrationRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	integration, err := sc.GetIntegration(pm.Corp, site, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(d.Id())
	err = d.Set("site_short_name", site)
	if err != nil {
		return err
	}
	err = d.Set("name", integration.Name)
	if err != nil {
		return err
	}
	err = d.Set("type", integration.Type)
	if err != nil {
		return err
	}
	err = d.Set("url", integration.URL)
	if err != nil {
		return err
	}
	err = d.Set("events", flattenStringArray(integration.Events))
	if err != nil {
		return err
	}
	return nil
}

func resourceSiteIntegrationUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	err := sc.UpdateIntegration(pm.Corp, site, d.Id(), sigsci.UpdateIntegrationBody{
		URL:    d.Get("url").(string),
		Events: expandStringArray(d.Get("events").(*schema.Set)),
	})
	if err != nil {
		log.Printf("[ERROR] %s. Could not update integration with ID %s in corp %s site %s", err.Error(), d.Id(), pm.Corp, site)
		d.SetId("")
		return nil
	}
	return resourceSiteIntegrationRead(d, m)
}

func resourceSiteIntegrationDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	err := sc.DeleteIntegration(pm.Corp, site, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
