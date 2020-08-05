package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"reflect"
)

func resourceCorpSignalTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceCorpSignalTagCreate,
		Update: resourceCorpSignalTagUpdate,
		Read:   resourceCorpSignalTagRead,
		Delete: resourceCorpSignalTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"short_name": {
				Type:        schema.TypeString,
				Description: "The display name of the signal tag",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Optional signal tag description",
				Required:    true,
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

func resourceCorpSignalTagCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp

	createResp, err := sc.CreateCorpSignalTag(corp, sigsci.CreateSignalTagBody{
		ShortName:   d.Get("short_name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	_, err = sc.GetCorpSignalTagByID(corp, createResp.TagName)
	if err != nil {
		return fmt.Errorf("%s. Could note create signaltag with ID %s in corp %s.Please re-run", err.Error(), createResp.TagName, corp)
	}
	d.SetId(createResp.TagName)
	return resourceCorpSignalTagRead(d, m)
}

func resourceCorpSignalTagUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	before, err := sc.GetCorpSignalTagByID(corp, d.Id())
	if err != nil {
		d.SetId("")
		return fmt.Errorf("%s. Could not find signaltag to update with ID %s in corp %s", err.Error(), d.Id(), corp)
	}
	updateSignalTagBody := sigsci.UpdateSignalTagBody{
		Description: d.Get("description").(string),
	}

	_, err = sc.UpdateCorpSignalTagByID(corp, d.Id(), updateSignalTagBody)
	if err != nil {
		return fmt.Errorf("%s. Could not update signaltag with ID %s in corp %s", err.Error(), d.Id(), corp)
	}
	after, err := sc.GetCorpSignalTagByID(corp, d.Id())
	if err == nil && reflect.DeepEqual(before, after) {
		return fmt.Errorf("Update failed for signaltag ID %s in corp %s\ngot:\n%#v\nexpected:\n%#v\nPlease re-run",
			d.Id(), corp, after, updateSignalTagBody)
	}
	return resourceCorpSignalTagRead(d, m)
}

func resourceCorpSignalTagRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	signaltag, err := sc.GetCorpSignalTagByID(corp, d.Id())
	if err != nil {
		d.SetId("")
		return fmt.Errorf("%s. Could not find signaltag with ID %s in corp %s", err.Error(), d.Id(), corp)
	}
	err = d.Set("short_name", signaltag.ShortName)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("description", signaltag.Description)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("configurable", signaltag.Configurable)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("informational", signaltag.Informational)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("needs_response", signaltag.NeedsResponse)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

func resourceCorpSignalTagDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	err := sc.DeleteCorpSignalTagByID(corp, d.Id())
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	_, err = sc.GetCorpSignalTagByID(corp, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	return fmt.Errorf("Could not delete signaltag with ID %s in corp %s, Please re-run", d.Id(), corp)
}
