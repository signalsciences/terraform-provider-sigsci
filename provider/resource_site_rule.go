package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceSiteRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of rule (request, signal, rateLimit)",
				Required:    true,
				ForceNew:    true,
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
			"requestlogging": {
				Type:             schema.TypeString,
				Description:      "Indicates whether to store the logs for requests that match the rule's conditions (sampled) or not store them (none). This field is only available for rules of type `request`. Not valid for `signal` or `rateLimit`.",
				Optional:         true,
				DiffSuppressFunc: suppressRequestLoggingDefaultDiffs,
			},
			"actions": {
				Type:        schema.TypeSet,
				Description: "Actions",
				Optional:    true,
				MaxItems:    2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "(addSignal, allow, block, browserChallenge, excludeSignal) (rateLimit rule valid values: logRequest, blockSignal)",
							Required:    true,
						},
						"signal": {
							Type:        schema.TypeString,
							Description: "signal id to tag",
							Optional:    true,
						},
						"response_code": {
							Type:         schema.TypeInt,
							Description:  "HTTP code agent for agent to respond with. range: 301, 302, or 400-599, defaults to '406' if not provided. Only valid with the 'block' action type.",
							Optional:     true,
							ValidateFunc: validateActionResponseCode,
						},
						"redirect_url": {
							Type:         schema.TypeString,
							Description:  "URL to redirect to when blocking response code is set to 301 or 302",
							Optional:     true,
							ValidateFunc: validateActionRedirectURL,
						},
						"allow_interactive": {
							Type:        schema.TypeBool,
							Description: "Allows toggling between a non-interactive and interactive browser challenge. Only valid with the 'browserChallenge' action type.",
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
							Description: "type: single - (equals, doesNotEqual, contains, doesNotContain, greaterEqual, lesserEqual, like, notLike, exists, doesNotExist, matches, doesNotMatch, inList, notInList)",
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
										Description: "type: single - (equals, doesNotEqual, contains, doesNotContain, greaterEqual, lesserEqual, like, notLike, exists, doesNotExist, matches, doesNotMatch, inList, notInList)",
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
													Description: "type: single - (equals, doesNotEqual, contains, doesNotContain, greaterEqual, lesserEqual, like, notLike, exists, doesNotExist, matches, doesNotMatch, inList, notInList)",
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
			"rate_limit": {
				Type:        schema.TypeSet,
				Description: "Rate Limit",
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"threshold": {
							Type:        schema.TypeInt,
							Description: "threshold",
							Required:    true,
						},
						"interval": {
							Type:        schema.TypeInt,
							Description: "interval in minutes (1, 5, 10)",
							Required:    true,
						},
						"duration": {
							Type:        schema.TypeInt,
							Description: "duration in seconds (300 < x < 3600)",
							Required:    true,
						},
						"client_identifiers": {
							Type:        schema.TypeSet,
							Description: "Client Identifiers",
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "(ip, requestHeader, requestCookie, postParameter, signalPayload)",
										Required:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "",
										Optional:    true,
									},
									"key": {
										Type:        schema.TypeString,
										Description: "",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
		Create:   resourceSiteRuleCreate,
		Read:     resourceSiteRuleRead,
		Update:   resourceSiteRuleUpdate,
		Delete:   resourceSiteRuleDelete,
		Importer: &siteImporter,
	}
}

func resourceSiteRuleCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	siteRulesBody := sigsci.CreateSiteRuleBody{
		Type:           d.Get("type").(string),
		GroupOperator:  d.Get("group_operator").(string),
		Enabled:        d.Get("enabled").(bool),
		Reason:         d.Get("reason").(string),
		Signal:         d.Get("signal").(string),
		Expiration:     d.Get("expiration").(string),
		RequestLogging: d.Get("requestlogging").(string),
	}

	siteRulesBody.Conditions = expandRuleConditions(d.Get("conditions").(*schema.Set))
	siteRulesBody.Actions = expandRuleActions(d.Get("actions").(*schema.Set))
	siteRulesBody.RateLimit = expandRuleRateLimit(d.Get("rate_limit").(*schema.Set))

	rule, err := sc.CreateSiteRule(corp, site, siteRulesBody)
	if err != nil {
		return err
	}
	_, err = sc.GetSiteRuleByID(corp, site, rule.ID)
	if err != nil {
		return err
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
		return nil
	}

	d.SetId(rule.ID)
	err = d.Set("site_short_name", site)
	if err != nil {
		return err
	}
	err = d.Set("type", rule.Type)
	if err != nil {
		return err
	}
	err = d.Set("group_operator", rule.GroupOperator)
	if err != nil {
		return err
	}
	err = d.Set("enabled", rule.Enabled)
	if err != nil {
		return err
	}
	err = d.Set("reason", rule.Reason)
	if err != nil {
		return err
	}
	err = d.Set("signal", rule.Signal)
	if err != nil {
		return err
	}
	err = d.Set("expiration", rule.Expiration)
	if err != nil {
		return err
	}
	err = d.Set("requestlogging", rule.RequestLogging)
	if err != nil {
		return err
	}
	err = d.Set("actions", flattenRuleActions(rule.Actions, true))
	if err != nil {
		return err
	}
	err = d.Set("conditions", flattenRuleConditions(rule.Conditions))
	if err != nil {
		return err
	}
	err = d.Set("rate_limit", flattenRuleRateLimit(rule.RateLimit))
	if err != nil {
		return err
	}
	return nil
}

func resourceSiteRuleUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	corp := pm.Corp
	site := d.Get("site_short_name").(string)

	updateSiteRuleBody := sigsci.CreateSiteRuleBody{
		Type:           d.Get("type").(string),
		GroupOperator:  d.Get("group_operator").(string),
		Enabled:        d.Get("enabled").(bool),
		Reason:         d.Get("reason").(string),
		Signal:         d.Get("signal").(string),
		Expiration:     d.Get("expiration").(string),
		RequestLogging: d.Get("requestlogging").(string),
	}

	updateSiteRuleBody.Conditions = expandRuleConditions(d.Get("conditions").(*schema.Set))
	updateSiteRuleBody.Actions = expandRuleActions(d.Get("actions").(*schema.Set))
	updateSiteRuleBody.RateLimit = expandRuleRateLimit(d.Get("rate_limit").(*schema.Set))

	_, err := sc.UpdateSiteRuleByID(corp, site, d.Id(), updateSiteRuleBody)
	if err != nil {
		return err
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
		return err
	}
	_, err = sc.GetSiteRuleByID(corp, site, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("could not delete rule with ID %s in corp %s site %s. Please re-run", d.Id(), corp, site)
}
