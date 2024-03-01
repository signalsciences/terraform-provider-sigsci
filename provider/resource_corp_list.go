package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceCorpList() *schema.Resource {
	return &schema.Resource{
		Create: resourceCorpListCreate,
		Update: resourceCorpListUpdate,
		Read:   resourceCorpListRead,
		Delete: resourceCorpListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Descriptive list name",
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !validStringLength(val.(string), 3, 32) {
						return nil, []error{fmt.Errorf(`received name %q is invalid. should be min len 3, max len 32`, val.(string))}
					}
					return nil, nil
				},
			},
			"type": {
				Type:        schema.TypeString,
				Description: "List types (string, ip, country, wildcard, signal)",
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !existsInString(val.(string), "string", "ip", "country", "wildcard", "signal") {
						return nil, []error{fmt.Errorf(`received type %q is invalid. should be "string", "ip", "country", "wildcard" or "signal"`, val.(string))}
					}
					return nil, nil
				},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Optional list description",
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !validStringLength(val.(string), 0, 140) {
						return nil, []error{fmt.Errorf(`received description %q is invalid. should be max len 140`, val.(string))}
					}
					return nil, nil
				},
			},
			"entries": {
				Type:        schema.TypeSet,
				Description: "List entries",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCorpListCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp

	corpListBody := sigsci.CreateListBody{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Entries:     expandStringArray(d.Get("entries").(*schema.Set)),
	}

	list, err := sc.CreateCorpList(corp, corpListBody)
	if err != nil {
		return err
	}
	d.SetId(list.ID)

	return resourceCorpListRead(d, m)
}

func resourceCorpListRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp

	list, err := sc.GetCorpListByID(corp, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(list.ID)
	err = d.Set("name", list.Name)
	if err != nil {
		return err
	}
	err = d.Set("type", list.Type)
	if err != nil {
		return err
	}
	err = d.Set("description", list.Description)
	if err != nil {
		return err
	}
	err = d.Set("entries", flattenStringArray(list.Entries))
	if err != nil {
		return err
	}
	return nil
}

func resourceCorpListDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp

	err := sc.DeleteCorpListByID(corp, d.Id())
	if err != nil {
		return err
	}
	_, err = sc.GetCorpListByID(corp, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	return fmt.Errorf("could not delete list with ID %s in corp %s. Please re-run", d.Id(), corp)

}
func resourceCorpListUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp

	before, err := sc.GetCorpListByID(corp, d.Id())
	if err != nil {
		return fmt.Errorf("%s. Could not find list with ID %s in corp %s", err.Error(), d.Id(), corp)
	}
	existingEntries := before.Entries
	newEntries := expandStringArray(d.Get("entries").(*schema.Set))
	additions, deletions := getListAdditionsDeletions(existingEntries, newEntries)
	updateCorpListBody := sigsci.UpdateListBody{
		Description: d.Get("description").(string),
		Entries: sigsci.Entries{
			Additions: additions,
			Deletions: deletions,
		},
	}
	_, err = sc.UpdateCorpListByID(corp, d.Id(), updateCorpListBody)
	if err != nil {
		return fmt.Errorf("%s. Could not update list with ID %s in corp %s. Please re-run", err.Error(), d.Id(), corp)
	}
	return resourceCorpListRead(d, m)
}
