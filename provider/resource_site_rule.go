package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"reflect"
)

func resourceSiteRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteRuleCreate,
		Update: resourceSiteRuleUpdate,
		Read:   resourceSiteRuleRead,
		Delete: resourceSiteRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, // this only sets the id. Probably a better way
		},
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of rule (request, signal exclusion)",
				Required:    true,
			},
			"group_operator": {
				Type:        schema.TypeString,
				Description: "Conditions that must be matched when evaluating the request (all, any)",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "enable the rule",
				Required:    true,
			},
			"signal": {
				Type:        schema.TypeString,
				Description: "The signal id of the signal being excluded",
				Required:    true,
			},
			"reason": {
				Type:        schema.TypeString,
				Description: "Description of the rule",
				Required:    true,
			},
			"expiration": {
				Type:        schema.TypeString,
				Description: "Date the rule will automatically be disabled. If rule is always enabled, will return empty string",
				Required:    true,
			},
			"actions": {
				Type:        schema.TypeSet,
				Description: "Actions",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "(block, allow, exclude)",
							Required:    true,
						},
					},
				},
			},
			"conditions": {
				Type:        schema.TypeSet,
				Description: "Conditions",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "(group, single)",
							Required:    true,
						},
						"field": {
							Type:        schema.TypeString,
							Description: "type: single - (scheme, method, path, useragent, domain, ip, responseCode, agentname, paramname, paramvalue, country, name, valueString, valueIp, signalType)",
							Optional:    true,
						},
						"operator": {
							Type:        schema.TypeString,
							Description: "type: single - (equals, doesNotEqual, contains, doesNotContain, like, notLike, exists, doesNotExist, inList, notInList)",
							Optional:    true,
						},
						"group_operator": {
							Type:        schema.TypeString,
							Description: "type: group - Conditions that must be matched when evaluating the request (all, any)",
							Optional:    true,
							// ConflictsWith: []string{"conditions.0.operator", "conditions.0.value", "conditions.0.field", "conditions.1.operator", "conditions.1.value", "conditions.1.field"}, does # work here
						},
						"value": {
							Type:        schema.TypeString,
							Description: "type: single - See request fields (https://docs.signalsciences.net/using-signal-sciences/features/rules/#request-fields)",
							Optional:    true,
						},
						"conditions": {
							Type:        schema.TypeSet,
							Description: "Conditions",
							Optional:    true,
							// ConflictsWith: []string{"conditions.0.operator"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "(group, single)",
										Required:    true,
									},
									"field": {
										Type:        schema.TypeString,
										Description: "type: single - (scheme, method, path, useragent, domain, ip, responseCode, agentname, paramname, paramvalue, country, name, valueString, valueIp, signalType)",
										Optional:    true,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "type: single - (equals, doesNotEqual, contains, doesNotContain, like, notLike, exists, doesNotExist, inList, notInList)",
										Optional:    true,
									},
									"group_operator": {
										Type:        schema.TypeString,
										Description: "type: group - Conditions that must be matched when evaluating the request (all, any)",
										Optional:    true,
										// ConflictsWith: []string{"conditions.0.operator", "conditions.0.value", "conditions.0.field", "conditions.1.operator", "conditions.1.value", "conditions.1.field"}, does # work here
									},
									"value": {
										Type:        schema.TypeString,
										Description: "type: single - See request fields (https://docs.signalsciences.net/using-signal-sciences/features/rules/#request-fields)",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSiteRuleCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	siteRulesBody := sigsci.CreateSiteRuleBody{
		Type:          d.Get("type").(string),
		GroupOperator: d.Get("group_operator").(string),
		Enabled:       d.Get("enabled").(bool),
		Reason:        d.Get("reason").(string),
		Signal:        d.Get("signal").(string),
		Expiration:    d.Get("expiration").(string),
	}

	siteRulesBody.Conditions = expandRuleConditions(d.Get("conditions").(*schema.Set))
	siteRulesBody.Actions = expandRuleActions(d.Get("actions").(*schema.Set))

	rule, err := sc.CreateSiteRule(corp, site, siteRulesBody)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	_, err = sc.GetSiteRuleByID(corp, site, rule.ID)
	if err != nil {
		return fmt.Errorf("%s. Could not create rule with ID %s in corp %s in site %s. Please re-run", err.Error(), rule.ID, corp, site)
	}
	d.SetId(rule.ID)
	return resourceSiteRuleRead(d, m)
}

func resourceSiteRuleRead(d *schema.ResourceData, m interface{}) error {

	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	rule, err := sc.GetSiteRuleByID(corp, site, d.Id())
	if err != nil {
		d.SetId("")
		return fmt.Errorf("%s. Could not find rule with ID %s in corp %s site %s", err.Error(), corp, site, d.Id())
	}

	err = d.Set("site_short_name", site)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("type", rule.Type)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("group_operator", rule.GroupOperator)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("enabled", rule.Enabled)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("reason", rule.Reason)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("signal", rule.Signal)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("expiration", rule.Expiration)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("actions", flattenRuleActions(rule.Actions))
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	err = d.Set("conditions", flattenRuleConditions(rule.Conditions))
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	return nil
}

func resourceSiteRuleUpdate(d *schema.ResourceData, m interface{}) error {

	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	updateSiteRuleBody := sigsci.CreateSiteRuleBody{
		Type:          d.Get("type").(string),
		GroupOperator: d.Get("group_operator").(string),
		Enabled:       d.Get("enabled").(bool),
		Reason:        d.Get("reason").(string),
		Signal:        d.Get("signal").(string),
		Expiration:    d.Get("expiration").(string),
	}

	updateSiteRuleBody.Conditions = expandRuleConditions(d.Get("conditions").(*schema.Set))
	updateSiteRuleBody.Actions = expandRuleActions(d.Get("actions").(*schema.Set))

	_, err := sc.UpdateSiteRuleByID(corp, site, d.Id(), updateSiteRuleBody)
	if err != nil {
		return fmt.Errorf("%s. Could not update redaction with Id %s in corp %s site %s", err.Error(), d.Id(), corp, site)
	}
	rule, err := sc.GetSiteRuleByID(corp, site, d.Id())
	if err == nil && !reflect.DeepEqual(updateSiteRuleBody, rule.CreateSiteRuleBody) {
		return fmt.Errorf("Update failed for rule ID %s in corp %s in site %s\ngot:\n%#v\nexpected:\n%#v\nPlease re-run",
			d.Id(), corp, site, rule.CreateSiteRuleBody, updateSiteRuleBody)
	}
	return resourceSiteRuleRead(d, m)
}

func resourceSiteRuleDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	err := sc.DeleteSiteRuleByID(corp, site, d.Id())
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	_, err = sc.GetSiteRuleByID(corp, site, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("Could not delete rule with ID %s in corp %s site %s. Please re-run", d.Id(), corp, site)
}
