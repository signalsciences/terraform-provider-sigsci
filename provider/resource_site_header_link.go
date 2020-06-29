package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceSiteHeaderLink() *schema.Resource {
	return &schema.Resource{
		Create:   resourceSiteHeaderLinkCreate,
		Update:   resourceSiteHeaderLinkUpdate,
		Read:     resourceSiteHeaderLinkRead,
		Delete:   resourceSiteHeaderLinkDelete,
		Importer: &siteImporter,
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Descriptive list name",
				Required:    true,
				ForceNew:    true, // Hopefully this can be changed in the api later
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of header, either 'request' or 'response'",
				Required:    true,
				ForceNew:    true, // Hopefully this can be changed in the api later
			},
			"link_name": {
				Type:        schema.TypeString,
				Description: "Name of header link for display purposes",
				Optional:    true,
			},
			"link": {
				Type:        schema.TypeString,
				Description: "External link",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSiteHeaderLinkCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	links, err := sc.AddHeaderLink(pm.Corp, d.Get("site_short_name").(string), sigsci.HeaderLinkBody{
		Type:     d.Get("type").(string),
		Name:     d.Get("name").(string),
		LinkName: d.Get("link_name").(string),
		Link:     d.Get("link").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(links[len(links)-1].ID)
	return resourceSiteHeaderLinkRead(d, m)
}

func resourceSiteHeaderLinkRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	link, err := sc.GetHeaderLink(pm.Corp, site, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(d.Id())
	err = d.Set("site_short_name", site)
	if err != nil {
		return err
	}
	err = d.Set("name", link.Name)
	if err != nil {
		return err
	}
	err = d.Set("type", link.Type)
	if err != nil {
		return err
	}
	err = d.Set("link_name", link.LinkName)
	if err != nil {
		return err
	}
	err = d.Set("link", link.Link)
	if err != nil {
		return err
	}
	return nil
}

// No update, will generate a unique ID
func resourceSiteHeaderLinkUpdate(d *schema.ResourceData, m interface{}) error {
	err := resourceSiteHeaderLinkDelete(d, m)
	if err != nil {
		return err
	}

	err = resourceSiteHeaderLinkCreate(d, m)
	if err != nil {
		return err
	}

	return resourceSiteHeaderLinkRead(d, m)
}

func resourceSiteHeaderLinkDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	err := sc.DeleteHeaderLink(pm.Corp, site, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
