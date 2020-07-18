package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

//TODO implement sweepers for everyone
func TestAccResourceTemplatedRulesCRUD(t *testing.T) {
	resourceName := "sigsci_site_templated_rule.test_site_templated_rule"
	testSite := "test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_templated_rule" "test_template_rule" {
                      site_short_name = "%s"
                      name            = "LOGINATTEMPT"
                      detections {
                        name    = "detection1"
                        enabled = "true"
                        fields {
                          name  = "path"
                          value = "/auth/*"
                        }
                      }
                      alerts {
                        long_name          = "alert1"
                        interval           = 60
                        threshold          = 10
                        skip_notifications = true
                        enabled            = true
                        action             = "info"
                      }
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "LOGINATTEMPT"),
				),
			},
			{
				Config: fmt.Sprintf(`
                     resource "sigsci_site_templated_rule" "test_template_rule" {
                      site_short_name = "%s"
                      name            = "LOGINATTEMPT"
                      detections {
                        name    = "detection1"
                        enabled = "true"
                        fields {
                          name  = "path"
                          value = "/auth/*"
                        }
                      }
                      alerts {
                        long_name          = "alert1"
                        interval           = 60
                        threshold          = 10
                        skip_notifications = true
                        enabled            = true
                        action             = "info"
                      }
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "LOGINATTEMPT"),
				),
			},
		},
	})
}
