package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

//TODO implement sweepers for everyone
func TestAccResourceTemplatedRulesCRUD(t *testing.T) {
	t.Parallel()
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
                      }
                      alerts {
                        long_name          = "alert2"
                        interval           = 60
                        threshold          = 11
                        skip_notifications = false
                        enabled            = false 
                        action             = "info"
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
					resource.TestCheckResourceAttrSet(resourceName, "alerts.824284280.id"),
					resource.TestCheckResourceAttr(resourceName, "alerts.824284280.long_name", "alert1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.824284280.interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "alerts.824284280.threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "alerts.824284280.skip_notifications", "true"),
					resource.TestCheckResourceAttr(resourceName, "alerts.824284280.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alerts.824284280.action", "info"),

					resource.TestCheckResourceAttrSet(resourceName, "alerts.1907559371.id"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1907559371.long_name", "alert2"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1907559371.interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1907559371.threshold", "11"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1907559371.skip_notifications", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1907559371.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1907559371.action", "info"),
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
			        }
                    alerts {
			          long_name          = "alert3"
			          interval           = 60
			          threshold          = 12
			          skip_notifications = false
			          enabled            = true
			          action             = "template"
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

					resource.TestCheckResourceAttrSet(resourceName, "alerts.4052116893.id"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4052116893.long_name", "alert1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4052116893.interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4052116893.threshold", "14"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4052116893.skip_notifications", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4052116893.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.4052116893.action", "template"),

					resource.TestCheckResourceAttr(resourceName, "alerts.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "alerts.1132798798.id"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1132798798.long_name", "alert3"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1132798798.interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1132798798.threshold", "12"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1132798798.skip_notifications", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1132798798.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alerts.1132798798.action", "template"),
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
