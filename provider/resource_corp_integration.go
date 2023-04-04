package provider

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceCorpIntegration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCorpIntegrationCreate,
		Update: resourceCorpIntegrationUpdate,
		Read:   resourceCorpIntegrationRead,
		Delete: resourceCorpIntegrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "One of (mailingList, slack, microsoftTeams)",
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !existsInString(val.(string), "mailingList", "slack", "microsoftTeams") {
						return nil, []error{fmt.Errorf(`received type %q is invalid. should be "mailingList", "slack", or "microsoftTeams"`, val.(string))}
					}
					return nil, nil
				},
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

func resourceCorpIntegrationCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	integration, err := sc.AddCorpIntegration(pm.Corp, sigsci.IntegrationBody{
		Type:   d.Get("type").(string),
		URL:    d.Get("url").(string),
		Events: expandStringArray(d.Get("events").(*schema.Set)),
	})
	if err != nil {
		return err
	}

	d.SetId(integration.ID)
	return resourceCorpIntegrationRead(d, m)
}

func resourceCorpIntegrationRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	integration, err := sc.GetCorpIntegration(pm.Corp, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(d.Id())
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

func resourceCorpIntegrationUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	err := sc.UpdateCorpIntegration(pm.Corp, d.Id(), sigsci.UpdateIntegrationBody{
		URL:    d.Get("url").(string),
		Events: expandStringArray(d.Get("events").(*schema.Set)),
	})
	if err != nil {
		log.Printf("[ERROR] %s. Could not update corp integration with ID %s in corp %s", err.Error(), d.Id(), pm.Corp)
		d.SetId("")
		return nil
	}
	return resourceCorpIntegrationRead(d, m)
}

func resourceCorpIntegrationDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	err := sc.DeleteCorpIntegration(pm.Corp, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
