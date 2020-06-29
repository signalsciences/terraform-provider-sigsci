package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"log"
)

func resourceSiteList() *schema.Resource {
	return &schema.Resource{
		Create:   resourceSiteListCreate,
		Update:   resourceSiteListUpdate,
		Read:     resourceSiteListRead,
		Delete:   resourceSiteListDelete,
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
				Description: "List types (string, ip, country, wildcard)",
				Required:    true,
				ForceNew:    true, // Hopefully this can be changed in the api later
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

func resourceSiteListCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	list, err := sc.CreateSiteList(pm.Corp, d.Get("site_short_name").(string), sigsci.CreateListBody{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Entries:     expandStringArray(d.Get("entries").(*schema.Set)),
	})
	if err != nil {
		return err
	}
	d.SetId(list.ID)
	return resourceSiteListRead(d, m)
}

func resourceSiteListRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	list, err := sc.GetSiteListByID(pm.Corp, site, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(d.Id())
	err = d.Set("site_short_name", site)
	if err != nil {
		return err
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

func resourceSiteListUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	readList, err := sc.GetSiteListByID(pm.Corp, site, d.Id())
	if err != nil {
		log.Printf("[ERROR] %s. Could not find list with ID %s in corp %s site %s", err.Error(), pm.Corp, site, d.Id())
		d.SetId("")
		return nil
	}
	existingEntries := readList.Entries
	newEntries := expandStringArray(d.Get("entries").(*schema.Set))
	additions, deletions := getListAdditionsDeletions(existingEntries, newEntries)

	updateSiteListBody := sigsci.UpdateListBody{
		Description: d.Get("description").(string),
		Entries: sigsci.Entries{
			Additions: additions,
			Deletions: deletions,
		},
	}

	_, err = sc.UpdateSiteListByID(pm.Corp, site, d.Id(), updateSiteListBody)
	if err != nil {
		d.SetId("")
		return err
	}

	return resourceSiteListRead(d, m)
}

func resourceSiteListDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	err := sc.DeleteSiteListByID(pm.Corp, site, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
