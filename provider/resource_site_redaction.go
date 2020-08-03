package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceSiteRedaction() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteRedactionCreate,
		Update: resourceSiteRedactionUpdate,
		Read:   resourceSiteRedactionRead,
		Delete: resourceSiteRedactionDelete,
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
			"field": {
				Type:        schema.TypeString,
				Description: "Field Name",
				Required:    true,
			},
			"redaction_type": {
				Type:        schema.TypeInt,
				Description: "Type of redaction (0: Request Parameter, 1: Request Header, 2: Response Header)",
				Required:    true,
			},
		},
	}
}

func resourceSiteRedactionCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	redaction, err := sc.CreateSiteRedaction(pm.Corp, d.Get("site_short_name").(string), sigsci.CreateSiteRedactionBody{
		Field:         d.Get("field").(string),
		RedactionType: d.Get("redaction_type").(int),
	})
	if err != nil {
		return err
	}
	d.SetId(redaction.ID)
	return resourceSiteRedactionRead(d, m)
}

func resourceSiteRedactionRead(d *schema.ResourceData, m interface{}) error {

	pm := m.(providerMetadata)
	sc := pm.Client

	redaction, err := sc.GetSiteRedactionByID(pm.Corp, d.Get("site_short_name").(string), d.Id())
	if err != nil {
		return err
	}

	d.SetId(redaction.ID)
	err = d.Set("site_short_name", d.Get("site_short_name").(string))
	if err != nil {
		return err
	}
	err = d.Set("field", redaction.Field)
	if err != nil {
		return err
	}
	err = d.Set("redaction_type", redaction.RedactionType)
	if err != nil {
		return err
	}
	return nil
}

func resourceSiteRedactionUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	redaction, err := sc.UpdateSiteRedactionByID(pm.Corp, d.Get("site_short_name").(string), d.Id(), sigsci.CreateSiteRedactionBody{
		Field:         d.Get("field").(string),
		RedactionType: d.Get("redaction_type").(int),
	})
	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(redaction.ID)
	return resourceSiteRedactionRead(d, m)
}

func resourceSiteRedactionDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	err := sc.DeleteSiteRedactionByID(pm.Corp, d.Get("site_short_name").(string), d.Id())
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
