package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceSiteBlocklist() *schema.Resource {
	return &schema.Resource{
		Create:   resourceSiteBlocklistCreate,
		Update:   resourceSiteBlocklistUpdate,
		Read:     resourceSiteBlocklistRead,
		Delete:   resourceSiteBlocklistDelete,
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
				Description: "Source IP Address to Blocklist",
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

func resourceSiteBlocklistCreate(d *schema.ResourceData, m interface{}) error {
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

	createResp, err := sc.AddBlacklistIP(corp, site, sigsci.ListIPBody{
		Source:  d.Get("source").(string),
		Note:    d.Get("note").(string),
		Expires: expires,
	})

	if err != nil {
		return fmt.Errorf("%s. Could not create Blocklist", err.Error())
	}
	d.SetId(createResp.ID)
	return resourceSiteBlocklistRead(d, m)
}

func resourceSiteBlocklistRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	Blocklists, err := sc.ListBlacklistIPs(corp, site)
	if err != nil {
		d.SetId("")
		return nil
	}

	var Blocklist *sigsci.ListIP
	for _, w := range Blocklists {
		if w.ID == d.Id() {
			Blocklist = &w
		}
	}

	if Blocklist == nil {
		id := d.Id()
		d.SetId("")
		return fmt.Errorf("could not find Blocklist with id %s", id)
	}

	err = d.Set("source", Blocklist.Source)
	if err != nil {
		return err
	}
	err = d.Set("note", Blocklist.Note)
	if err != nil {
		return err
	}
	if !Blocklist.Expires.IsZero() {
		err = d.Set("expires", Blocklist.Expires.Format(time.RFC3339))
		if err != nil {
			return err
		}
	}
	d.SetId(Blocklist.ID)
	return nil
}

// There is no update api, we must delete and recreate every update :(
// This function should never be called but should work anyways
func resourceSiteBlocklistUpdate(d *schema.ResourceData, m interface{}) error {
	err := resourceSiteBlocklistDelete(d, m)
	if err != nil {
		return err
	}
	return resourceSiteBlocklistCreate(d, m)
}

func resourceSiteBlocklistDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	err := sc.DeleteBlacklistIP(corp, site, d.Id())

	if err == nil {
		d.SetId("")
		return nil
	}
	return fmt.Errorf("could not delete Blocklist with ID %s for site %s in corp %s", d.Id(), site, corp)
}
