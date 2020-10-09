package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"time"
)

func resourceSiteAllowlist() *schema.Resource {
	return &schema.Resource{
		Create:   resourceSiteAllowlistCreate,
		Update:   resourceSiteAllowlistUpdate,
		Read:     resourceSiteAllowlistRead,
		Delete:   resourceSiteAllowlistDelete,
		Importer: &siteImporter,
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
				ForceNew:    true,
			},
			"source": {
				Type:        schema.TypeString,
				Description: "Source IP Address to allowlist",
				Required:    true,
			},
			"note": {
				Type:        schema.TypeString,
				Description: "Note/Description associated with the tag.",
				Required:    true,
			},
			"expires": {
				Type:        schema.TypeString,
				Description: "Optional RFC3339-formatted datetime in the future. Omit this paramater if it does not expire.",
				Optional:    true,
			},
		},
	}
}

func resourceSiteAllowlistCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	var expires time.Time
	if expiresStr := d.Get("expires").(string); expiresStr != "" {
		parse, err := time.Parse(time.RFC3339, expiresStr)
		if err != nil {
			expires = parse
		}
	}

	createResp, err := sc.AddWhitelistIP(corp, site, sigsci.ListIPBody{
		Source:  d.Get("source").(string),
		Note:    d.Get("note").(string),
		Expires: expires,
	})

	if err != nil {
		return fmt.Errorf("%s. Could not create allowlist", err.Error())
	}
	d.SetId(createResp.ID)
	return resourceSiteAllowlistRead(d, m)
}

func resourceSiteAllowlistRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	allowlists, err := sc.ListWhitelistIPs(corp, site)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("%s. Could not find allowlists for site %s in corp %s", err.Error(), site, corp)
	}
	var allowlist *sigsci.ListIP
	for _, w := range allowlists {
		if w.ID == d.Id() {
			allowlist = &w
		}
	}

	if allowlist == nil {
		d.SetId("")
		return fmt.Errorf("could not find allowlist with id %s", d.Id())
	}

	err = d.Set("source", allowlist.Source)
	if err != nil {
		return err
	}
	err = d.Set("note", allowlist.Note)
	if err != nil {
		return err
	}
	if !allowlist.Expires.IsZero() {
		err = d.Set("expires", allowlist.Expires.Format(time.RFC3339))
		if err != nil {
			return err
		}
	}
	d.SetId(allowlist.ID)
	return nil
}

// There is no update api, we must delete and recreate every update :(
// This function should never be called but should work anyways
func resourceSiteAllowlistUpdate(d *schema.ResourceData, m interface{}) error {
	err := resourceSiteAllowlistDelete(d, m)
	if err != nil {
		return err
	}
	return resourceSiteAllowlistCreate(d, m)
}

func resourceSiteAllowlistDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	err := sc.DeleteWhitelistIP(corp, site, d.Id())

	if err == nil {
		d.SetId("")
		return nil
	}
	return fmt.Errorf("could not delete allowlist with ID %s for site %s in corp %s", d.Id(), site, corp)
}
