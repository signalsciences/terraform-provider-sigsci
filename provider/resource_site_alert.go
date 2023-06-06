package provider

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceSiteAlert() *schema.Resource {
	return &schema.Resource{
		Create:   resourceSiteAlertCreate,
		Update:   resourceSiteAlertUpdate,
		Read:     resourceSiteAlertRead,
		Delete:   resourceSiteAlertDelete,
		Importer: &siteImporter,
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
				ForceNew:    true,
			},
			"tag_name": {
				Type:        schema.TypeString,
				Description: "The name of the tag whose occurrences the alert is watching. Must match an existing tag",
				Required:    true,
			},
			"long_name": {
				Type:        schema.TypeString,
				Description: "description",
				Optional:    true,
			},
			"interval": {
				Type:        schema.TypeInt,
				Description: "The number of minutes of past traffic to examine. Must be 1, 10 or 60.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if existsInInt(val.(int), 1, 10, 60) {
						return nil, nil
					}
					return nil, []error{errors.New("interval must be 1, 10, or 60")}
				},
			},
			"threshold": {
				Type:        schema.TypeInt,
				Description: "The number of occurrences of the tag in the interval needed to trigger the alert. Min 1, Max 10000",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if existsInRange(val.(int), 1, 10000) {
						return nil, nil
					}
					return nil, []error{errors.New("threshold must be between 1 and 10000")}
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "A flag to toggle this alert.",
				Optional:    true,
			},
			"action": {
				Type:        schema.TypeString,
				Description: "A flag that describes what happens when the alert is triggered. 'info' creates an incident in the dashboard. 'flagged' creates an incident and blocks traffic for 24 hours. Must be info or flagged.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if existsInString(val.(string), "info", "flagged") {
						return nil, nil
					}
					return nil, []error{errors.New("action must be 'info' or 'flagged'")}
				},
			},
			"skip_notifications": {
				Type:        schema.TypeBool,
				Description: "A flag to skip notifications",
				Optional:    true,
			},
			"block_duration_seconds": {
				Type:        schema.TypeInt,
				Description: "The number of seconds this alert is active.",
				Optional:    true,
			},
		},
	}
}

func resourceSiteAlertCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	alert, err := sc.CreateCustomAlert(pm.Corp, d.Get("site_short_name").(string), sigsci.CustomAlertBody{
		TagName:              d.Get("tag_name").(string),
		LongName:             d.Get("long_name").(string),
		Interval:             d.Get("interval").(int),
		Threshold:            d.Get("threshold").(int),
		Enabled:              d.Get("enabled").(bool),
		Action:               d.Get("action").(string),
		SkipNotifications:    d.Get("skip_notifications").(bool),
		BlockDurationSeconds: d.Get("block_duration_seconds").(int),
	})
	if err != nil {
		return err
	}
	d.SetId(alert.ID)
	return resourceSiteAlertRead(d, m)
}

func resourceSiteAlertRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	alert, err := sc.GetCustomAlert(pm.Corp, d.Get("site_short_name").(string), d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(alert.ID)
	err = d.Set("site_short_name", d.Get("site_short_name").(string))
	if err != nil {
		return err
	}
	err = d.Set("tag_name", alert.TagName)
	if err != nil {
		return err
	}
	err = d.Set("long_name", alert.LongName)
	if err != nil {
		return err
	}
	err = d.Set("interval", alert.Interval)
	if err != nil {
		return err
	}
	err = d.Set("threshold", alert.Threshold)
	if err != nil {
		return err
	}
	err = d.Set("enabled", alert.Enabled)
	if err != nil {
		return err
	}
	err = d.Set("action", alert.Action)
	if err != nil {
		return err
	}
	err = d.Set("skip_notifications", alert.SkipNotifications)
	if err != nil {
		return err
	}
	err = d.Set("block_duration_seconds", alert.BlockDurationSeconds)
	if err != nil {
		return err
	}

	return nil
}

func resourceSiteAlertUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	alert, err := sc.UpdateCustomAlert(pm.Corp, d.Get("site_short_name").(string), d.Id(), sigsci.CustomAlertBody{
		TagName:              d.Get("tag_name").(string),
		LongName:             d.Get("long_name").(string),
		Interval:             d.Get("interval").(int),
		Threshold:            d.Get("threshold").(int),
		Enabled:              d.Get("enabled").(bool),
		Action:               d.Get("action").(string),
		SkipNotifications:    d.Get("skip_notifications").(bool),
		BlockDurationSeconds: d.Get("block_duration_seconds").(int),
	})
	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(alert.ID)
	return resourceSiteAlertRead(d, m)
}

func resourceSiteAlertDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	err := sc.DeleteCustomAlert(pm.Corp, d.Get("site_short_name").(string), d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
