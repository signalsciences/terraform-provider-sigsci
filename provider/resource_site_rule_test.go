package provider

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/signalsciences/go-sigsci"
	"reflect"
	"testing"
)

func TestACCResourceSiteRule_basic(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckSiteRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_rule" "test"{
                        site_short_name="%s"
                        type= "signal"
                        group_operator="any"
                        enabled= true
                        reason= "Example site rule update"
                        signal= "SQLI"
                        expiration= ""
                        conditions {
                            type="single"
                            field="ip"
                            operator="equals"
                            value="1.2.3.4"
                        }
                        conditions {
                            type="single"
                            field="ip"
                            operator="equals"
                            value="1.2.3.5"
                        }
                        actions {
                            type="excludeSignal"
                        }
                }`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckSiteRuleExists(resourceName),
					testCheckSiteRulesAreEqual(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "signal"),
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "reason", "Example site rule update"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1859487095.type", "excludeSignal"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.conditions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.field", "ip"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.group_operator", ""),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.operator", "equals"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.type", "single"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.value", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.conditions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.field", "ip"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.group_operator", ""),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.operator", "equals"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.type", "single"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2383694574.value", "1.2.3.5"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: fmt.Sprintf("%s:", testSite),
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateCheck:    testAccImportStateCheckFunction(1),
			},
		},
	})
}

func TestACCResourceSiteRuleRateLimit_basic(t *testing.T) {
	resourceName := "sigsci_site_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckSiteRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_signal_tag" "test_tag" {
                      site_short_name = "%s" 
                      name            = "My new tag"
                      description     = "test description"
                    }
                    resource "sigsci_site_rule" "test" {
                        site_short_name="%s"
                        type= "rateLimit"
                        group_operator="any"
                        enabled= true
                        reason= "Example site rule update"
                        signal= sigsci_site_signal_tag.test_tag.id
                        expiration= ""
                        conditions {
                            type="single"
                            field="ip"
                            operator="equals"
                            value="1.2.3.4"
                        }
                        actions {
                            type="logRequest"
                        }
                        rate_limit = {
                            threshold=10
                            interval=10
                            duration=600
                        }
                }`, testSite, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					//testCheckSiteRuleExists(resourceName),
					//testCheckSiteRulesAreEqual(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "rateLimit"),
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "reason", "Example site rule update"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1301899462.type", "logRequest"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rate_limit.threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "rate_limit.interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "rate_limit.duration", "600"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: fmt.Sprintf("%s:", testSite),
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateCheck:    testAccImportStateCheckFunction(1),
			},
		},
	})
}

func TestACCResourceSiteRuleConditionSignal(t *testing.T) {

	resourceName := "sigsci_site_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckSiteRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_rule" "test" {
                        site_short_name = "%s" 
                        type            = "request"
                        group_operator  = "all"
                        enabled         = true
                        reason          = "Example site rule update"
                        expiration      = ""

                        conditions {
                            type     = "multival"
                            field    = "signal"
                            group_operator = "all"
                            operator = "exists"
                            conditions {
                                field    = "signalType"
                                operator = "equals"
                                type     = "single"
                                value    = "RESPONSESPLIT"
                            }
                        }

                        conditions {
                            type     = "group"
                            group_operator = "any"
                            conditions {
                                field    = "useragent"
                                operator = "like"
                                type     = "single"
                                value    = "python-requests*"
                            }

                            conditions {
                                type     = "multival"
                                field    = "requestHeader"
                                operator = "doesNotExist"
                                group_operator = "all"
                                conditions {
                                    field    = "name"
                                    operator = "equals"
                                    type     = "single"
                                    value    = "cookie"
                                }
                            }

                            conditions {
                                type     = "multival"
                                field    = "signal"
                                operator = "exists"
                                group_operator = "any"
                                conditions {
                                    field    = "signalType"
                                    operator = "equals"
                                    type     = "single"
                                    value    = "TORNODE"
                                }
                                conditions {
                                    field    = "signalType"
                                    operator = "equals"
                                    type     = "single"
                                    value    = "SIGSCI-IP"
                                }
                                conditions {
                                    field    = "signalType"
                                    operator = "equals"
                                    type     = "single"
                                    value    = "SCANNER"
                                }
                            }
                        }

                        actions {
                            type = "block"
                        }
                }`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					testInspect(),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2455721190.conditions.3887678098.conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conditions.1840769124.conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2455721190.conditions.2522856064.conditions.#", "3"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: fmt.Sprintf("%s:", testSite),
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateCheck:    testAccImportStateCheckFunction(1),
			},
		},
	})
}

func testCheckSiteRuleExists(name string) resource.TestCheckFunc {
	var testFunc resource.TestCheckFunc = func(s *terraform.State) error {
		rsrc, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[ERROR] the module %s does not have a resource called %s", s.RootModule().Path, name)
		}

		is := rsrc.Primary
		if is == nil {
			return fmt.Errorf("[ERROR] No primary instance: %s in %s", name, s.RootModule().Path)
		}
		pm := testAccProvider.Meta().(providerMetadata)
		sc := pm.Client
		rule, err := sc.GetSiteRuleByID(pm.Corp, is.Attributes["site_short_name"], is.Attributes["id"])
		if err != nil {
			return err
		}
		if rule.ID != is.Attributes["id"] {
			return fmt.Errorf("[ERROR] the rule ids did not match expected :%s\t actual: %s", is.Attributes["id"], rule.ID)
		}
		return nil
	}
	return testFunc
}
func testCheckSiteRulesAreEqual(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pm := testAccProvider.Meta().(providerMetadata)
		sc := pm.Client
		rsrc, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[ERROR] the module %s does not have a resource called %s", s.RootModule().Path, name)
		}

		is := rsrc.Primary
		if is == nil {
			return fmt.Errorf("[ERROR] No primary instance: %s in %s", name, s.RootModule().Path)
		}

		actual, err := sc.GetSiteRuleByID(pm.Corp, is.Attributes["site_short_name"], is.Attributes["id"])

		if err != nil {
			return err
		}
		expected := sigsci.CreateSiteRuleBody{
			Type:          "signal",
			GroupOperator: "any",
			Enabled:       true,
			Reason:        "Example site rule update",
			Signal:        "SQLI",
			Expiration:    "",
			Conditions: []sigsci.Condition{
				sigsci.Condition{
					Type:     "single",
					Field:    "ip",
					Operator: "equals",
					Value:    "1.2.3.5",
				},
				sigsci.Condition{
					Type:     "single",
					Field:    "ip",
					Operator: "equals",
					Value:    "1.2.3.4",
				},
			},
			Actions: []sigsci.Action{
				sigsci.Action{
					Type: "excludeSignal",
				},
			},
		}
		if !reflect.DeepEqual(expected, actual.CreateSiteRuleBody) {
			spewConf := spew.NewDefaultConfig()
			spewConf.SortKeys = true
			return fmt.Errorf("not equal: \nexpected\n%s\nactual\n%s", spewConf.Sdump(expected), spewConf.Sdump(actual.CreateSiteRuleBody))
		}
		return nil
	}
}
func testCheckSiteGroupRulesAreEqual(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pm := testAccProvider.Meta().(providerMetadata)
		sc := pm.Client
		rsrc, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[ERROR] the module %s does not have a resource called %s", s.RootModule().Path, name)
		}

		is := rsrc.Primary
		if is == nil {
			return fmt.Errorf("[ERROR] No primary instance: %s in %s", name, s.RootModule().Path)
		}

		actual, err := sc.GetSiteRuleByID(pm.Corp, is.Attributes["site_short_name"], is.Attributes["id"])

		if err != nil {
			return err
		}
		expected := sigsci.CreateSiteRuleBody{
			Type:          "signal",
			GroupOperator: "all",
			Enabled:       true,
			Reason:        "Example site rule group",
			Signal:        "SQLI",
			Expiration:    "",
			Conditions: []sigsci.Condition{
				{
					Type:          "group",
					GroupOperator: "any",
					Conditions: []sigsci.Condition{
						{
							Type:     "single",
							Field:    "ip",
							Operator: "equals",
							Value:    "9.10.11.12",
						},
					},
				},
			},
			Actions: []sigsci.Action{
				{
					Type: "excludeSignal",
				},
			},
		}
		if !reflect.DeepEqual(expected, actual.CreateSiteRuleBody) {
			spewConf := spew.NewDefaultConfig()
			spewConf.SortKeys = true
			return fmt.Errorf("not equal: \nexpected\n%s\nactual\n%s", spewConf.Sdump(expected), spewConf.Sdump(actual.CreateSiteRuleBody))
		}
		return nil
	}
}
func testAccImportStateCheckFunction(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %#v", expectedStates, len(s), s)
		}
		return nil
	}
}

func testACCCheckSiteRuleDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client
	resourceType := "sigsci_site_rule"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		readresp, err := sc.GetSiteRuleByID(pm.Corp, resource.Primary.Attributes["site_short_name"], resource.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("%s %#v still exists", resourceType, readresp)
		}
	}
	return nil
}
