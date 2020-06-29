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
			"name": {
				Type:        schema.TypeString,
				Description: "Identifying name of the site",
				Required:    true,
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
				Description: "Agent IP anonimization mode - 'EU' or 'off'",
				Optional:    true,
				Removed:     "Documented but causes an error when provided",
				Default:     "off", // TODO Default is off in the docs, but not enforced by api
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
				Default:     406,
			},
		},
	}
}

func createSite(d *schema.ResourceData, m interface{}) error {
	sc := m.(sigsci.Client)
	corp := d.Get("corp").(string)
	site, err := sc.CreateSite(corp, sigsci.CreateSiteBody{
		Name:                 d.Get("name").(string),
		DisplayName:          d.Get("display_name").(string),
		AgentLevel:           d.Get("agent_level").(string),
		AgentAnonMode:        d.Get("agent_anon_mode").(string),
		BlockHTTPCode:        d.Get("block_http_code").(int),
		BlockDurationSeconds: d.Get("block_duration_seconds").(int),
	})

	if err != nil {
		return err
	}
	d.SetId(corpSiteToId(corp, site.Name))

	return readSite(d, m)
}

func readSite(d *schema.ResourceData, m interface{}) error {
	sc := m.(sigsci.Client)
	corp, sitename := idToCorpSite(d.Id())
	site, err := sc.GetSite(corp, sitename)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("[ERROR] No site found with name %s in %s", sitename, corp)
	}

	d.SetId(corpSiteToId(corp, site.Name)) // No inherent id, combination of corp and site should be unique
	d.Set("agent_level", site.AgentLevel)
	d.Set("block_duration_seconds", site.BlockDurationSeconds)
	d.Set("block_http_code", site.BlockHTTPCode)
	d.Set("corp", corp)
	d.Set("display_name", site.DisplayName)
	d.Set("name", site.Name)

	return nil
}

func updateSite(d *schema.ResourceData, m interface{}) error {
	sc := m.(sigsci.Client)
	updateSiteBody := sigsci.UpdateSiteBody{
		DisplayName:          d.Get("display_name").(string),
		AgentLevel:           d.Get("agent_level").(string),
		BlockDurationSeconds: d.Get("block_duration_seconds").(int),
	}
	corp, site := idToCorpSite(d.Id())
	_, err := sc.UpdateSite(corp, site, updateSiteBody)
	if err != nil {
		return err
	}

	return readSite(d, m)
}

func deleteSite(d *schema.ResourceData, m interface{}) error {
	sc := m.(sigsci.Client)
	err := sc.DeleteSite(idToCorpSite(d.Id()))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
