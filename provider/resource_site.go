package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: createSite,
		Update: updateSite,
		Read:   readSite,
		Delete: deleteSite,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				d.Set("short_name", d.Id())
				d.SetId(d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"short_name": {
				Type:        schema.TypeString,
				Description: "Identifying name of the site",
				Required:    true,
				ForceNew:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name of the site",
				Required:    true,
			},
			"agent_level": {
				Type:        schema.TypeString,
				Description: "Agent action level - 'block', 'log' or 'off'",
				Optional:    true,
				Default:     "log",
			},
			"agent_anon_mode": { // Has issues on create -- will always be default, will update just fine to the correct value
				Type:        schema.TypeString,
				Description: "Agent IP anonymization mode - \"\" (empty string) or 'EU'",
				Optional:    true,
				Default:     "",
			},
			"block_duration_seconds": { // Has issues on create -- will always be default, will update just fine to the correct value
				Type:        schema.TypeInt,
				Description: "Duration to block an IP in seconds",
				Optional:    true,
				Default:     86400,
			},
			"block_http_code": { // CANNOT UPDATE THIS FIELD,
				Type:        schema.TypeInt,
				Description: "HTTP response code to send when traffic is being blocked",
				Computed:    true,
			},
			"primary_agent_key": {
				Type:        schema.TypeMap,
				Description: "The sites primary Agent key",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"access_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func createSite(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site, err := sc.CreateSite(pm.Corp, sigsci.CreateSiteBody{
		Name:                 d.Get("short_name").(string),
		DisplayName:          d.Get("display_name").(string),
		AgentLevel:           d.Get("agent_level").(string),
		AgentAnonMode:        d.Get("agent_anon_mode").(string),
		BlockHTTPCode:        d.Get("block_http_code").(int),
		BlockDurationSeconds: d.Get("block_duration_seconds").(int),
	})

	if err != nil {
		return err
	}
	d.SetId(site.Name)

	// For whatever reason, you cannot create without default values, but you may update them later
	// If these are not the default values, update
	if d.Get("block_duration_seconds").(int) != 86400 ||
		d.Get("agent_anon_mode").(string) != "" {
		return updateSite(d, m)
	}

	return readSite(d, m)
}

func readSite(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	sitename := d.Get("short_name").(string)
	site, err := sc.GetSite(corp, sitename)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(site.Name)
	err = d.Set("agent_level", site.AgentLevel)
	if err != nil {
		return err
	}
	err = d.Set("block_duration_seconds", site.BlockDurationSeconds)
	if err != nil {
		return err
	}
	err = d.Set("block_http_code", site.BlockHTTPCode)
	if err != nil {
		return err
	}
	err = d.Set("agent_anon_mode", site.AgentAnonMode)
	if err != nil {
		return err
	}
	err = d.Set("display_name", site.DisplayName)
	if err != nil {
		return err
	}
	err = d.Set("short_name", site.Name)
	if err != nil {
		return err
	}

	primaryAgentKey, err := sc.GetSitePrimaryAgentKey(corp, sitename)
	if err != nil {
		return err
	}

	return d.Set("primary_agent_key", map[string]interface{}{
		"name":       primaryAgentKey.Name,
		"secret_key": primaryAgentKey.SecretKey,
		"access_key": primaryAgentKey.AccessKey,
	})
}

func updateSite(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("short_name").(string)
	_, err := sc.UpdateSite(corp, site, sigsci.UpdateSiteBody{
		DisplayName:          d.Get("display_name").(string),
		AgentLevel:           d.Get("agent_level").(string),
		BlockDurationSeconds: d.Get("block_duration_seconds").(int),
		BlockHTTPCode:        d.Get("block_http_code").(int),
		AgentAnonMode:        d.Get("agent_anon_mode").(string),
	})
	if err != nil {
		return err
	}

	return readSite(d, m)
}

func deleteSite(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	err := sc.DeleteSite(pm.Corp, d.Get("short_name").(string))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
