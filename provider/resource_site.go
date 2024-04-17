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
			"attack_threshold": {
				Type:        schema.TypeSet,
				Description: "List entries",
				Required:    false,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"block_duration_seconds": { // Has issues on create -- will always be default, will update just fine to the correct value
				Type:        schema.TypeInt,
				Description: "Duration to block an IP in seconds",
				Optional:    true,
				Default:     86400,
			},
			"block_http_code": {
				Type:        schema.TypeInt,
				Description: "HTTP response code to send when traffic is being blocked",
				Optional:    true,
				Default:     406,
			},
			"block_redirect_url": {
				Type:        schema.TypeString,
				Description: "URL to redirect to when blocking with a '301' or '302' HTTP status code",
				Optional:    true,
			},
			"client_ip_rules": {
				Type:        schema.TypeSet,
				Description: "Headers used for assigning client IPs to requests",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			"immediate_block": {
				Type:        schema.TypeBool,
				Description: "Immediately block requests that contain attack signals",
				Optional:    true,
			},
			"primary_agent_key": {
				Type:        schema.TypeMap,
				Description: "The sites primary Agent key",
				Computed:    true,
				Sensitive:   true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_key": {
							Type:     schema.TypeString,
							Computed: true,
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
		AttackThresholds:     expandAttackThresholds(d.Get("attack_threshold").(*schema.Set)),
		BlockHTTPCode:        d.Get("block_http_code").(int),
		BlockDurationSeconds: d.Get("block_duration_seconds").(int),
		BlockRedirectURL:     d.Get("block_redirect_url").(string),
		ClientIPRules:        expandClientIPRules(d.Get("client_ip_rules").(*schema.Set)),
		ImmediateBlock:       d.Get("immediate_block").(bool),
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
	err = d.Set("block_redirect_url", site.BlockRedirectURL)
	if err != nil {
		return err
	}
	err = d.Set("client_ip_rules", flattenClientIPRules(site.ClientIPRules))
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

	err = d.Set("immediate_block", site.ImmediateBlock)
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
		AttackThresholds:     expandAttackThresholds(d.Get("attack_threshold").(*schema.Set)),
		BlockDurationSeconds: d.Get("block_duration_seconds").(int),
		BlockHTTPCode:        d.Get("block_http_code").(int),
		BlockRedirectURL:     d.Get("block_redirect_url").(string),
		AgentAnonMode:        d.Get("agent_anon_mode").(string),
		ClientIPRules:        expandClientIPRules(d.Get("client_ip_rules").(*schema.Set)),
		ImmediateBlock:       d.Get("immediate_block").(bool),
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
