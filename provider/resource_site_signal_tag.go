package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

// TODO maybe rename to corp tag
func resourceSiteSignalTag() *schema.Resource {
	return &schema.Resource{
		Create:   resourceSiteSignalTagCreate,
		Update:   resourceSiteSignalTagUpdate,
		Read:     resourceSiteSignalTagRead,
		Delete:   resourceSiteSignalTagDelete,
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
				Description: "The display name of the signal tag. Must be 3-25 char.",
				Required:    true,
				ForceNew:    true, // TODO Hopefully this can be changed in the api later
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !validStringLength(val.(string), 3, 25) {
						return nil, []error{fmt.Errorf(`received name %q is invalid. should be min len 3, max len 25`, val.(string))}
					}
					return nil, nil
				},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description",
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !validStringLength(val.(string), 0, 140) {
						return nil, []error{fmt.Errorf(`received description is invalid. should be max len 140`)}
					}
					return nil, nil
				},
			},
			"configurable": {
				Type:        schema.TypeBool,
				Description: "configurable",
				Computed:    true,
			},
			"informational": {
				Type:        schema.TypeBool,
				Description: "informational",
				Computed:    true,
			},
			"needs_response": {
				Type:        schema.TypeBool,
				Description: "need response",
				Computed:    true,
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
	site := d.Get("site_short_name").(string)
	tag, err := sc.GetSiteSignalTagByID(pm.Corp, site, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(tag.TagName)
	err = d.Set("site_short_name", site)
	if err != nil {
		return err
	}
	err = d.Set("name", tag.ShortName)
	if err != nil {
		return err
	}
	err = d.Set("description", tag.Description)
	if err != nil {
		return err
	}
	err = d.Set("configurable", tag.Configurable)
	if err != nil {
		return err
	}
	err = d.Set("informational", tag.Informational)
	if err != nil {
		return err
	}
	err = d.Set("needs_response", tag.NeedsResponse)
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
