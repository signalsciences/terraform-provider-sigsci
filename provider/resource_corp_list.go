package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"reflect"
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
			},
			"type": {
				Type:        schema.TypeString,
				Description: "List types (string, ip, country, wildcard)",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Optional list description",
				Optional:    true,
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
		return fmt.Errorf("%s. Could not find list with ID %s in corp %s", err.Error(), d.Id(), corp)
	}
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
	after, err := sc.GetCorpListByID(corp, d.Id())
	if err == nil && reflect.DeepEqual(after.CreateListBody, before.CreateListBody) {
		return fmt.Errorf("Update failed for list ID %s in corp %s\ngot:\n%#v\nexpected:\n%#v.Please re-run",
			d.Id(), corp, after.CreateListBody, updateCorpListBody)
	}
	return resourceCorpListRead(d, m)
}
