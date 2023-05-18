package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceCorpRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCorpRuleCreate,
		Update: resourceCorpRuleUpdate,
		Read:   resourceCorpRuleRead,
		Delete: resourceCorpRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"site_short_names": {
				Type:        schema.TypeSet,
				Description: "Sites with the rule available. Rules with a global corpScope will return '[]'.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of rule (request, signal exclusion)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if existsInString(val.(string), "request", "signal") {
						return nil, nil
					}
					return nil, []error{}
				},
			},
			"corp_scope": {
				Type:        schema.TypeString,
				Description: "Whether the rule is applied to all sites or to specific sites. (global, specificSites)",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "enable the rule",
				Required:    true,
			},
			"group_operator": {
				Type:        schema.TypeString,
				Description: "Conditions that must be matched when evaluating the request (all, any)",
				Required:    true,
			},
			"actions": {
				Type:        schema.TypeSet,
				Description: "Actions",
				Required:    true,
				MaxItems:    2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "(block, allow, addSignal, excludeSignal)",
							Required:    true,
						},
						"signal": {
							Type:        schema.TypeString,
							Description: "signal id",
							Optional:    true,
						},
					},
				},
			},
			"conditions": {
				Type:        schema.TypeSet,
				Description: "Conditions",
				Required:    true,
				MaxItems:    10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "(group, multival, single)",
							Required:    true,
						},
						"field": {
							Type:         schema.TypeString,
							Description:  fmt.Sprintf("types:\n    - single - (%s)\n    - multival - (%s)", strings.Join(KnownSingleConditionFields, ", "), strings.Join(KnownMultivalConditionFields, ", ")),
							Optional:     true,
							ValidateFunc: validateConditionField,
						},
						"operator": {
							Type:        schema.TypeString,
							Description: "type: single - (equals, doesNotEqual, contains, doesNotContain, like, notLike, exists, doesNotExist, inList, notInList)",
							Optional:    true,
						},
						"group_operator": {
							Type:        schema.TypeString,
							Description: "type: group, multival - Conditions that must be matched when evaluating the request (all, any)",
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
							MaxItems:    10,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "(group, multival, single)",
										Required:    true,
									},
									"field": {
										Type:         schema.TypeString,
										Description:  fmt.Sprintf("types:\n    - single - (%s)\n    - multival - (%s)", strings.Join(KnownSingleConditionFields, ", "), strings.Join(KnownMultivalConditionFields, ", ")),
										Optional:     true,
										ValidateFunc: validateConditionField,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "type: single - (equals, doesNotEqual, contains, doesNotContain, like, notLike, exists, doesNotExist, inList, notInList)",
										Optional:    true,
									},
									"group_operator": {
										Type:        schema.TypeString,
										Description: "type: group, multival - Conditions that must be matched when evaluating the request (all, any)",
										Optional:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "type: single - See request fields (https://docs.fastly.com/signalsciences/using-signal-sciences/rules/defining-rule-conditions/#fields)",
										Optional:    true,
									},
									"conditions": {
										Type:        schema.TypeSet,
										Description: "Conditions",
										Optional:    true,
										MaxItems:    10,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Description: "(group, multival, single)",
													Required:    true,
												},
												"field": {
													Type:         schema.TypeString,
													Description:  fmt.Sprintf("types:\n    - single - (%s)\n    - multival - (%s)", strings.Join(KnownSingleConditionFields, ", "), strings.Join(KnownMultivalConditionFields, ", ")),
													Optional:     true,
													ValidateFunc: validateConditionField,
												},
												"operator": {
													Type:        schema.TypeString,
													Description: "type: single - (equals, doesNotEqual, contains, doesNotContain, like, notLike, exists, doesNotExist, inList, notInList)",
													Optional:    true,
												},
												"group_operator": {
													Type:        schema.TypeString,
													Description: "type: group, multival - Conditions that must be matched when evaluating the request (all, any)",
													Optional:    true,
												},
												"value": {
													Type:        schema.TypeString,
													Description: "type: single - See request fields (https://docs.fastly.com/signalsciences/using-signal-sciences/rules/defining-rule-conditions/#fields)",
													Optional:    true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"requestlogging": {
				Type:        schema.TypeString,
				Description: "Indicates whether to store the logs for requests that match the rule's conditions (sampled) or not store them (none). This field is only available for rules of type `request`. Not valid for `signal`.",
				Optional:    true,
			},
			"signal": {
				Type:        schema.TypeString,
				Description: "The signal id of the signal being excluded",
				Optional:    true,
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
		},
	}
}

func resourceCorpRuleCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	corpRuleBody := sigsci.CreateCorpRuleBody{
		Type:          d.Get("type").(string),
		CorpScope:     d.Get("corp_scope").(string),
		Enabled:       d.Get("enabled").(bool),
		GroupOperator: d.Get("group_operator").(string),
		Reason:        d.Get("reason").(string),
		Signal:        d.Get("signal").(string),
		Expiration:    d.Get("expiration").(string),
	}

	corpRuleBody.SiteNames = expandStringArray(d.Get("site_short_names").(*schema.Set))
	corpRuleBody.Conditions = expandRuleConditions(d.Get("conditions").(*schema.Set))
	corpRuleBody.Actions = expandRuleActions(d.Get("actions").(*schema.Set))

	rule, err := sc.CreateCorpRule(corp, corpRuleBody)
	if err != nil {
		return err
	}
	d.SetId(rule.ID)
	return resourceCorpRuleRead(d, m)
}

func resourceCorpRuleUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp

	updateCorpRuleBody := sigsci.CreateCorpRuleBody{
		Type:          d.Get("type").(string),
		CorpScope:     d.Get("corp_scope").(string),
		Enabled:       d.Get("enabled").(bool),
		GroupOperator: d.Get("group_operator").(string),
		Reason:        d.Get("reason").(string),
		Signal:        d.Get("signal").(string),
		Expiration:    d.Get("expiration").(string),
	}

	updateCorpRuleBody.SiteNames = expandStringArray(d.Get("site_short_names").(*schema.Set))
	updateCorpRuleBody.Conditions = expandRuleConditions(d.Get("conditions").(*schema.Set))
	updateCorpRuleBody.Actions = expandRuleActions(d.Get("actions").(*schema.Set))

	_, err := sc.UpdateCorpRuleByID(corp, d.Id(), updateCorpRuleBody)
	if err != nil {
		return fmt.Errorf("%s. Could not update rule with ID %s in corp %s ", err.Error(), corp, d.Id())
	}
	return resourceCorpRuleRead(d, m)
}

func resourceCorpRuleRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp

	rule, err := sc.GetCorpRuleByID(corp, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(rule.ID)
	err = d.Set("type", rule.Type)
	if err != nil {
		return err
	}
	err = d.Set("corp_scope", rule.CorpScope)
	if err != nil {
		return err
	}
	err = d.Set("enabled", rule.Enabled)
	if err != nil {
		return err
	}
	err = d.Set("group_operator", rule.GroupOperator)
	if err != nil {
		return err
	}
	err = d.Set("signal", rule.Signal)
	if err != nil {
		return err
	}
	err = d.Set("reason", rule.Reason)
	if err != nil {
		return err
	}
	err = d.Set("expiration", rule.Expiration)
	if err != nil {
		return err
	}
	err = d.Set("site_short_names", flattenStringArray(rule.SiteNames))
	if err != nil {
		return err
	}
	err = d.Set("actions", flattenRuleActions(rule.Actions, false))
	if err != nil {
		return err
	}
	err = d.Set("conditions", flattenRuleConditions(rule.Conditions))
	if err != nil {
		return err
	}
	return nil
}

func resourceCorpRuleDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp

	err := sc.DeleteCorpRuleByID(corp, d.Id())
	if err != nil {
		return err
	}
	_, err = sc.GetCorpRuleByID(corp, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	return fmt.Errorf("could not delete rule with ID %s in corp %s. Please re-run", corp, d.Id())
}
