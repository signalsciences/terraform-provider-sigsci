package provider

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSiteAgentAlert() *schema.Resource {
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
				Description: "Integer value for interval. Must be 5, 10 or 60.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if existsInInt(val.(int), 5, 10, 60) {
						return nil, nil
					}
					return nil, []error{errors.New("interval must be 5, 10, or 60")}
				},
			},
			"threshold": {
				Type:        schema.TypeInt,
				Description: "The number of occurrences of the tag in the interval needed to trigger the alert. Min 0, Max 10000",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if existsInRange(val.(int), 0, 10000) {
						return nil, nil
					}
					return nil, []error{errors.New("threshold must be between 0 and 10000")}
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "A flag to toggle this alert.",
				Optional:    true,
			},
			"action": {
				Type:        schema.TypeString,
				Description: "Action for agent alert.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if existsInString(val.(string), "siteMetricInfo") {
						return nil, nil
					}
					return nil, []error{errors.New("action must be 'siteMetricInfo'")}
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
