package provider

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"net/url"
)

func resourceSiteMonitor() *schema.Resource {
	return &schema.Resource{
		Create:   resourceSiteMonitorCreate,
		Update:   resourceSiteMonitorUpdate,
		Read:     resourceSiteMonitorRead,
		Delete:   resourceSiteMonitorDelete,
		Importer: &siteImporter,
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
				ForceNew:    true,
			},
			"dashboard_id": {
				Type:        schema.TypeString,
				Description: "Dashboard ID",
				Required:    true,
				ForceNew:    true,
			},
			"share": {
				Type:        schema.TypeBool,
				Description: "Enables or Disables the site monitor",
				Optional:    true,
				Default:     true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "Share url",
				Computed:    true,
			},
		},
	}
}

func resourceSiteMonitorCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	monitor, err := sc.GenerateSiteMonitorDashboard(
		pm.Corp,
		d.Get("site_short_name").(string),
		d.Get("dashboard_id").(string),
	)

	if err != nil {
		return err
	}
	d.SetId(monitor.ID)

	// Share will always come back true at creation, disable if needed
	if d.Get("share") == false {
		err := resourceSiteMonitorUpdate(d, m)
		if err != nil {
			return err
		}
	}

	return resourceSiteMonitorRead(d, m)
}

func resourceSiteMonitorRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	monitors, err := sc.GetSiteMonitor(pm.Corp, site, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}
	var monitor sigsci.SiteMonitor
	for _, m := range monitors {
		if m.ID == d.Id() {
			monitor = m
		}
		if monitor.ID == "" {
			return errors.New(fmt.Sprintf("Could not find monitor with ID %s", d.Id()))
		}
	}
	if monitor.ID == "" {
		d.SetId("")
		return errors.New(fmt.Sprintf("Could not find site monitor with id %s", d.Id()))
	}

	d.SetId(monitor.ID)
	err = d.Set("site_short_name", site)
	if err != nil {
		return err
	}
	parse, err := url.Parse(monitor.URL)
	if err != nil {
		return err
	}
	dashboardID := parse.Query().Get("dashboardId")
	err = d.Set("dashboard_id", dashboardID)
	if err != nil {
		return err
	}
	err = d.Set("share", monitor.Share)
	if err != nil {
		return err
	}
	err = d.Set("url", monitor.URL)
	if err != nil {
		return err
	}
	return nil
}

func resourceSiteMonitorUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	err := sc.UpdateSiteMonitor(pm.Corp, site, d.Id(), sigsci.UpdateSiteMonitorBody{
		ID:    d.Id(),
		Share: d.Get("share").(bool),
	})
	if err != nil {
		d.SetId("")
		return err
	}
	return resourceSiteMonitorRead(d, m)
}

func resourceSiteMonitorDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	err := sc.DeleteSiteMonitor(pm.Corp, site, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
