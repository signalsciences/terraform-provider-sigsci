package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: createSite,
		Update: updateSite,
		Read:   readSite,
		Delete: deleteSite,
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
				Default:     "log", // TODO not in docs, but enforced by api
			},
			"agent_anon_mode": {
				Type:        schema.TypeString,
				Description: "Agent IP anonimization mode - \"\" (empty string) or 'EU'",
				Optional:    true,
				Default:     "",
			},
			"block_duration_seconds": {
				Type:        schema.TypeInt,
				Description: "Duration to block an IP in seconds",
				Optional:    true,
				Default:     86400,
			},
			"block_http_code": {
				Type:        schema.TypeInt,
				Description: "HTTP response code to send when when traffic is being blocked",
				Optional:    true,
				Default:     407,
			},
		},
	}
}

func createSite(d *schema.ResourceData, m interface{}) error {
	pm := m.(ProviderMetadata)
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
	d.SetId(corpSiteToId(pm.Corp, site.Name))

	return readSite(d, m)
}

func readSite(d *schema.ResourceData, m interface{}) error {
	pm := m.(ProviderMetadata)
	sc := pm.Client
	corp := pm.Corp
	sitename := d.Get("short_name").(string)
	site, err := sc.GetSite(corp, sitename)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("[ERROR] No site found with name %s in %s", sitename, corp)
	}

	d.SetId(corpSiteToId(corp, site.Name)) // No inherent id, combination of corp and site should be unique
	d.Set("agent_level", site.AgentLevel)
	d.Set("block_duration_seconds", site.BlockDurationSeconds)
	d.Set("block_http_code", site.BlockHTTPCode)
	d.Set("agent_anon_mode", site.AgentAnonMode)
	d.Set("corp", corp)
	d.Set("display_name", site.DisplayName)
	d.Set("short_name", site.Name)

	return nil
}

func updateSite(d *schema.ResourceData, m interface{}) error {
	pm := m.(ProviderMetadata)
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
	pm := m.(ProviderMetadata)
	sc := pm.Client
	err := sc.DeleteSite(pm.Corp, d.Get("short_name").(string))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
