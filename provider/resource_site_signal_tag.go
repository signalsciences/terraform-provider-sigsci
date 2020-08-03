package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceSiteSignalTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteSignalTagCreate,
		Update: resourceSiteSignalTagUpdate,
		Read:   resourceSiteSignalTagRead,
		Delete: resourceSiteSignalTagDelete,
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
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the signal tag",
				Required:    true,
				ForceNew:    true, // Hopefully this can be changed in the api later
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description",
				Optional:    true,
			},
		},
	}
}

func resourceSiteSignalTagCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	tag, err := sc.CreateSiteSignalTag(pm.Corp, d.Get("site_short_name").(string), sigsci.CreateSignalTagBody{
		ShortName:   d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return err
	}
	d.SetId(tag.TagName)
	return resourceSiteSignalTagRead(d, m)
}

func resourceSiteSignalTagRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	tag, err := sc.GetSiteSignalTagByID(pm.Corp, d.Get("site_short_name").(string), d.Id())
	if err != nil {
		return err
	}

	d.SetId(tag.TagName)
	err = d.Set("site_short_name", d.Get("site_short_name").(string))
	if err != nil {
		return err
	}
	err = d.Set("name", tag.LongName)
	if err != nil {
		return err
	}
	err = d.Set("description", tag.Description)
	if err != nil {
		return err
	}
	return nil
}

func resourceSiteSignalTagUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	tag, err := sc.UpdateSiteSignalTagByID(pm.Corp, d.Get("site_short_name").(string), d.Id(), sigsci.UpdateSignalTagBody{
		Description: d.Get("description").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(tag.TagName)
	return resourceSiteSignalTagRead(d, m)
}

func resourceSiteSignalTagDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	err := sc.DeleteSiteSignalTagByID(pm.Corp, d.Get("site_short_name").(string), d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
