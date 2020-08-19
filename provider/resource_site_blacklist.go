package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"time"
)

func resourceSiteBlacklist() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteBlacklistCreate,
		Update: resourceSiteBlacklistUpdate,
		Read:   resourceSiteBlacklistRead,
		Delete: resourceSiteBlacklistDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				site, id, err := resourceSiteImport(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("site_short_name", site)
				d.SetId(id)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
				ForceNew:    true,
			},
			"source": {
				Type:        schema.TypeString,
				Description: "Source IP Address to Blacklist",
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

func resourceSiteBlacklistCreate(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf("%s. Could not create Blacklist", err.Error())
	}
	d.SetId(createResp.ID)
	return resourceSiteBlacklistRead(d, m)
}

func resourceSiteBlacklistRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	Blacklists, err := sc.ListBlacklistIPs(corp, site)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("%s. Could not find Blacklists for site %s in corp %s", err.Error(), site, corp)
	}
	var Blacklist *sigsci.ListIP
	for _, w := range Blacklists {
		if w.ID == d.Id() {
			Blacklist = &w
		}
	}

	if Blacklist == nil {
		d.SetId("")
		return fmt.Errorf("could not find Blacklist with id %s", d.Id())
	}

	err = d.Set("source", Blacklist.Source)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("note", Blacklist.Note)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	if !Blacklist.Expires.IsZero() {
		err = d.Set("expires", Blacklist.Expires.Format(time.RFC3339))
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	}
	d.SetId(Blacklist.ID)
	return nil
}

// There is no update api, we must delete and recreate every update :(
// This function should never be called but should work anyways
func resourceSiteBlacklistUpdate(d *schema.ResourceData, m interface{}) error {
	err := resourceSiteBlacklistDelete(d, m)
	if err != nil {
		return err
	}
	return resourceSiteBlacklistCreate(d, m)
}

func resourceSiteBlacklistDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	err := sc.DeleteBlacklistIP(corp, site, d.Id())

	if err == nil {
		d.SetId("")
		return nil
	}
	return fmt.Errorf("could not delete Blacklist with ID %s for site %s in corp %s", d.Id(), site, corp)
}
