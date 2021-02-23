package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

//TODO implement sweepers for everyone
func TestAccResourceSiteAlertCRUD(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site_alert.test_site_alert"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_alert" "test_site_alert" {
                      site_short_name    = "%s"
                      tag_name           = "CMDEXE"
                      long_name          = "test_alert"
                      interval           = 10
                      threshold          = 12
                      enabled            = true
                      action             = "info"
                      skip_notifications = true 
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "tag_name", "CMDEXE"),
					resource.TestCheckResourceAttr(resourceName, "long_name", "test_alert"),
					resource.TestCheckResourceAttr(resourceName, "interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "12"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "action", "info"),
					resource.TestCheckResourceAttr(resourceName, "skip_notifications", "true"),
				),
			},
			{
				Config: fmt.Sprintf(`
                     resource "sigsci_site_alert" "test_site_alert" {
                      site_short_name    = "%s"
                      tag_name           = "SQLI"
                      long_name          = "test_alert 2"
                      interval           = 60
                      threshold          = 13
                      enabled            = false
                      action             = "flagged"
                      skip_notifications = false 
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "tag_name", "SQLI"),
					resource.TestCheckResourceAttr(resourceName, "long_name", "test_alert 2"),
					resource.TestCheckResourceAttr(resourceName, "interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "13"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "action", "flagged"),
					resource.TestCheckResourceAttr(resourceName, "skip_notifications", "false"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateIdPrefix: fmt.Sprintf("%s:", testSite),
				ImportStateCheck:    testAccImportStateCheckFunction(1),
				ImportStateVerify:   true,
			},
		},
	})
}
