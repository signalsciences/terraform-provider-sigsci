package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TODO implement sweepers for everyone
func TestAccResourceTemplatedRulesCRUD(t *testing.T) {
	resourceName := "sigsci_site_templated_rule.test_template_rule"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_templated_rule" "test_template_rule" {
                      site_short_name = "%s"
                      name            = "LOGINATTEMPT"

                      detections {
                        enabled = "true"
                        fields {
                          name = "path"
                          value = "/auth/*"
                        }
                      }
                      detections {
                        enabled = "false"
                        fields {
                          name = "path"
                          value = "/login/*"
                        }
                      }

                      alerts {
                        long_name          = "alert1"
                        interval           = 10
                        threshold          = 10
                        skip_notifications = true
                        enabled            = true
                        action             = "info"
                        block_duration_seconds = 54321
                      }

                      alerts {
                        long_name          = "alert2"
                        interval           = 60
                        threshold          = 11
                        skip_notifications = false
                        enabled            = false 
                        action             = "info"
                        block_duration_seconds = 0
                      }
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "LOGINATTEMPT"),

					resource.TestCheckResourceAttr(resourceName, "detections.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "detections.1895549889.id"),
					resource.TestCheckResourceAttr(resourceName, "detections.1895549889.name", "LOGINATTEMPT"),
					resource.TestCheckResourceAttr(resourceName, "detections.1895549889.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "detections.1895549889.fields.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "detections.1895549889.fields.4019068922.name", "path"),
					resource.TestCheckResourceAttr(resourceName, "detections.1895549889.fields.4019068922.value", "/auth/*"),

					resource.TestCheckResourceAttrSet(resourceName, "detections.2076537939.id"),
					resource.TestCheckResourceAttr(resourceName, "detections.2076537939.name", "LOGINATTEMPT"),
					resource.TestCheckResourceAttr(resourceName, "detections.2076537939.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "detections.2076537939.fields.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "detections.2076537939.fields.2831008451.name", "path"),
					resource.TestCheckResourceAttr(resourceName, "detections.2076537939.fields.2831008451.value", "/login/*"),

					resource.TestCheckResourceAttr(resourceName, "alerts.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "alerts.3578226761.id"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3578226761.long_name", "alert1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3578226761.interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3578226761.threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3578226761.skip_notifications", "true"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3578226761.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3578226761.action", "info"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3578226761.block_duration_seconds", "54321"),

					resource.TestCheckResourceAttrSet(resourceName, "alerts.3034807924.id"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3034807924.long_name", "alert2"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3034807924.interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3034807924.threshold", "11"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3034807924.skip_notifications", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3034807924.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3034807924.action", "info"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3034807924.block_duration_seconds", "0"),
				),
			},
			{
				Config: fmt.Sprintf(`
			       resource "sigsci_site_templated_rule" "test_template_rule" {
			        site_short_name = "%s"
			        name            = "LOGINATTEMPT"
			        detections {
			          enabled = "false"
			          fields {
                        name = "path"
			            value = "/admin/*"
			          }
                      fields {
                        name = "second"
			            value = "/backdoor/*"
			          }
			        }
			        alerts {
			          long_name          = "alert1"
			          interval           = 1
			          threshold          = 14
			          skip_notifications = false 
			          enabled            = false 
			          action             = "template"
			          block_duration_seconds = 54321
			        }
                    alerts {
			          long_name          = "alert3"
			          interval           = 60
			          threshold          = 12
			          skip_notifications = false
			          enabled            = true
			          action             = "template"
			          block_duration_seconds = 54321
			        }
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "LOGINATTEMPT"),

					resource.TestCheckResourceAttr(resourceName, "detections.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "detections.1573782214.id"),
					resource.TestCheckResourceAttr(resourceName, "detections.1573782214.name", "LOGINATTEMPT"),
					resource.TestCheckResourceAttr(resourceName, "detections.1573782214.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "detections.1573782214.fields.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "detections.1573782214.fields.2418319502.name", "path"),
					resource.TestCheckResourceAttr(resourceName, "detections.1573782214.fields.2418319502.value", "/admin/*"),
					resource.TestCheckResourceAttr(resourceName, "detections.1573782214.fields.4057994516.name", "second"),
					resource.TestCheckResourceAttr(resourceName, "detections.1573782214.fields.4057994516.value", "/backdoor/*"),

					resource.TestCheckResourceAttrSet(resourceName, "alerts.4126686736.id"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4126686736.long_name", "alert1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4126686736.interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4126686736.threshold", "14"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4126686736.skip_notifications", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4126686736.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4126686736.action", "template"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4126686736.action", "template"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4126686736.block_duration_seconds", "54321"),

					resource.TestCheckResourceAttr(resourceName, "alerts.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "alerts.3498615432.id"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3498615432.long_name", "alert3"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3498615432.interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3498615432.threshold", "12"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3498615432.skip_notifications", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3498615432.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3498615432.action", "template"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3498615432.block_duration_seconds", "54321"),
				),
			},
			{
				Config: fmt.Sprintf(`
			       resource "sigsci_site_templated_rule" "test_template_rule" {
			        site_short_name = "%s"
			        name            = "LOGINFAILURE"

                    detections {
			          enabled = "true"
			          fields {
                        name = "path"
			            value = "/admin/*"
			          }
                      fields {
                        name = "responseCode"
			            value = "303"
			          }
                      fields {
                        name = "responseHeaderName"
			            value = "Content-Type"
			          }
                      fields {
                        name = "responseHeaderValue"
			            value = "application/json"
			          }
			        }
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "LOGINFAILURE"),

					resource.TestCheckResourceAttr(resourceName, "detections.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "detections.900190133.id"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.name", "LOGINFAILURE"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.fields.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.fields.2418319502.name", "path"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.fields.2418319502.value", "/admin/*"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.fields.3194812060.name", "responseCode"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.fields.3194812060.value", "303"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.fields.37716642.name", "responseHeaderName"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.fields.37716642.value", "Content-Type"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.fields.2135544841.name", "responseHeaderValue"),
					resource.TestCheckResourceAttr(resourceName, "detections.900190133.fields.2135544841.value", "application/json"),

					resource.TestCheckResourceAttr(resourceName, "alerts.#", "0"),
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
		CheckDestroy: func(state *terraform.State) error {
			return nil
		},
	})
}

func TestAccResourceTemplatedRulesSSM(t *testing.T) {
	resourceName := "sigsci_site_templated_rule.aws_ssm"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_templated_rule" "aws_ssm" {
                      site_short_name = "%s"
                      name            = "AWS-SSRF"
                      detections {
					    enabled = "true"
					  }

					  alerts {
					    long_name          = ""
					    interval           = 0
						threshold          = 0
						skip_notifications = false
						enabled            = true
						action             = "blockImmediate"
						block_duration_seconds = 54321
                      }

					}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "AWS-SSRF"),

					resource.TestCheckResourceAttr(resourceName, "detections.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "detections.2705029964.id"),
					resource.TestCheckResourceAttr(resourceName, "detections.2705029964.name", "AWS-SSRF"),
					resource.TestCheckResourceAttr(resourceName, "detections.2705029964.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "detections.2705029964.fields.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "alerts.3025084325.id"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3025084325.long_name", ""),
					resource.TestCheckResourceAttr(resourceName, "alerts.3025084325.interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3025084325.threshold", "0"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3025084325.skip_notifications", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3025084325.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3025084325.action", "blockImmediate"),
					resource.TestCheckResourceAttr(resourceName, "alerts.3025084325.block_duration_seconds", "54321"),
				),
			},
		},
		CheckDestroy: func(state *terraform.State) error {
			return nil
		},
	})
}
