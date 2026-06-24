package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSite() *schema.Resource {
	return &schema.Resource{
		Read: readSiteDataSource,
		Schema: map[string]*schema.Schema{
			"short_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifying name of the site",
			},
			"agent_anon_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent IP anonimization mode - 'EU' or 'off'",
			},
			"agent_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent action level - 'block', 'log' or 'off'",
			},
			"agents_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's agents",
			},
			"alerts_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's alerts",
			},
			"analytics_events_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's analytics events",
			},
			"blacklist_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's blacklist",
			},
			"block_duration_secs": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Duration to block an IP in seconds",
			},
			"block_http_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "HTTP response code to send when traffic is being blocked",
			},
			"block_redirect_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to redirect to when blockHTTPCode is 301 or 302",
			},
			"client_ip_rules": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Headers used for assigning client IPs to requests",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created RFC3339 date time",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name of the site",
			},
			"events_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's events",
			},
			"header_links_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's header links",
			},
			"integrations_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's integrations",
			},
			"members_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's members",
			},
			"monitors_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's monitors",
			},
			"redactions_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's redactions",
			},
			"requests_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's requests",
			},
			"suspicious_ips_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's suspicious IPs",
			},
			"top_attacks_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's top attacks",
			},
			"whitelist_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the site's whitelist",
			},
		},
	}
}

func readSiteDataSource(d *schema.ResourceData, m any) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp

	shortName := d.Get("short_name").(string)

	// API documentation:
	// https://docs.fastly.com/signalsciences/api/#_corps__corpName__sites__siteName__get
	site, err := sc.GetSite(corp, shortName)
	if err != nil {
		return fmt.Errorf("could not get site %q in corp %q: %w", shortName, corp, err)
	}

	d.SetId(site.Name)

	fields := map[string]any{
		"agent_anon_mode":      site.AgentAnonMode,
		"agent_level":          site.AgentLevel,
		"agents_uri":           site.Agents["uri"],
		"alerts_uri":           site.Alerts["uri"],
		"analytics_events_uri": site.AnalyticsEvents["uri"],
		"blacklist_uri":        site.Blacklist["uri"],
		"block_duration_secs":  site.BlockDurationSeconds,
		"block_http_code":      site.BlockHTTPCode,
		"block_redirect_url":   site.BlockRedirectURL,
		"client_ip_rules":      flattenClientIPRules(site.ClientIPRules),
		"created":              site.Created.String(),
		"display_name":         site.DisplayName,
		"events_uri":           site.Events["uri"],
		"header_links_uri":     site.HeaderLinks["uri"],
		"integrations_uri":     site.Integrations["uri"],
		"members_uri":          site.Members["uri"],
		"monitors_uri":         site.Monitors["uri"],
		"redactions_uri":       site.Redactions["uri"],
		"requests_uri":         site.Requests["uri"],
		"suspicious_ips_uri":   site.SuspiciousIPs["uri"],
		"top_attacks_uri":      site.TopAttacks["uri"],
		"whitelist_uri":        site.Whitelist["uri"],
	}

	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("could not set %q: %w", k, err)
		}
	}

	return nil
}
