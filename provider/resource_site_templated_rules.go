package provider

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceSiteTemplatedRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplatedRuleCreate,
		Update: resourceTemplatedRuleUpdate,
		Read:   resourceTemplatedRuleRead,
		Delete: resourceTemplatedRuleDelete,
		//Importer: &schema.ResourceImporter{ //TODO try importing, make sure it works
		//	State: schema.ImportStatePassthrough, // this only sets the id. Probably a better way
		//},
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of templated rule.  This must match an existing templated rule e.g., LOGINATTEMPT, CMDEXE, XSS...",
				Required:    true,
				ForceNew:    true,
			},
			"detections": {
				Type:        schema.TypeSet,
				Description: "description",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"fields": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"alerts": {
				Type:        schema.TypeSet,
				Description: "Alerts",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"long_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								if intArrContains([]int{1, 10, 60}, val.(int)) {
									return nil, nil
								}
								return nil, []error{errors.New("alerts.interval must be 1, 10, or 60")}
							},
						},
						"threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"skip_notifications": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceTemplatedRuleCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)
	templateID := d.Get("name").(string)

	detectionAdds, _, _ := diffTemplateDetections(templateID, []sigsci.Detection{}, expandDetections(d.Get("detections").(*schema.Set)))
	alertAdds, _, _ := diffTemplateAlerts([]sigsci.Alert{}, expandAlerts(d.Get("alerts").(*schema.Set)))

	template, err := sc.UpdateSiteTemplateRuleByID(pm.Corp, site, templateID, sigsci.SiteTemplateRuleBody{
		DetectionAdds: detectionAdds,
		AlertAdds:     alertAdds,
	})
	if err != nil {
		return err
	}
	d.SetId(template.Name)
	return resourceTemplatedRuleRead(d, m)
}

func resourceTemplatedRuleRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	template, err := sc.GetSiteTemplateRuleByID(pm.Corp, d.Get("site_short_name").(string), d.Id())
	if err != nil {
		return err
	}

	d.SetId(template.Name)
	err = d.Set("site_short_name", d.Get("site_short_name").(string))
	if err != nil {
		return err
	}
	err = d.Set("name", template.Name)
	if err != nil {
		return err
	}
	err = d.Set("detections", flattenDetections(template.Detections))
	if err != nil {
		return err
	}
	err = d.Set("alerts", flattenAlerts(template.Alerts))
	if err != nil {
		return err
	}
	return nil
}

func resourceTemplatedRuleUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	existingTemplate, err := sc.GetSiteTemplateRuleByID(pm.Corp, site, d.Id())
	if err != nil {
		return err
	}

	detectionAdds, detectionUpdates, detectionDeletes := diffTemplateDetections(d.Id(), existingTemplate.Detections, expandDetections(d.Get("detections").(*schema.Set)))
	alertAdds, alertUpdates, alertDeletes := diffTemplateAlerts(existingTemplate.Alerts, expandAlerts(d.Get("alerts").(*schema.Set)))
	siteTemplate, err := sc.UpdateSiteTemplateRuleByID(pm.Corp, site, d.Id(), sigsci.SiteTemplateRuleBody{
		DetectionAdds:    detectionAdds,
		DetectionUpdates: detectionUpdates,
		DetectionDeletes: detectionDeletes,
		AlertAdds:        alertAdds,
		AlertUpdates:     alertUpdates,
		AlertDeletes:     alertDeletes,
	})
	if err != nil {
		return err
	}

	d.SetId(siteTemplate.Name)
	return resourceTemplatedRuleRead(d, m)
}

func resourceTemplatedRuleDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	site := d.Get("site_short_name").(string)

	existingTemplate, err := sc.GetSiteTemplateRuleByID(pm.Corp, site, d.Id())
	if err != nil {
		return err
	}

	_, err = sc.UpdateSiteTemplateRuleByID(pm.Corp, site, d.Id(), sigsci.SiteTemplateRuleBody{
		DetectionDeletes: existingTemplate.Detections,
		AlertDeletes:     existingTemplate.Alerts,
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
