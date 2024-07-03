package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TODO implement sweepers for everyone
func TestAccResourceSiteAgentAlertCRUD(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site_alert.test_site_alert"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_agent_alert" "test_site_agent_alert" {
						action                 = "siteMetricInfo"
						block_duration_seconds = 21600
						enabled                = true
						interval               = 10
						long_name              = "test_alert"
						site_short_name        = "%s"
						skip_notifications     = true
						tag_name               = "requests_total"
						threshold              = 8400
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "siteMetricInfo"),
					resource.TestCheckResourceAttr(resourceName, "block_duration_seconds", "21600"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "tag_name", "requests_total"),
					resource.TestCheckResourceAttr(resourceName, "long_name", "test_alert"),
					resource.TestCheckResourceAttr(resourceName, "interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "skip_notifications", "true"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "8400"),
				),
			},
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_agent_alert" "test_site_agent_alert" {
						action                 = "siteMetricInfo"
						block_duration_seconds = 21600
						enabled                = false
						interval               = 5
						long_name              = "test_alert"
						site_short_name        = "%s"
						skip_notifications     = true
						tag_name               = "requests_total"
						threshold              = 6300
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "siteMetricInfo"),
					resource.TestCheckResourceAttr(resourceName, "block_duration_seconds", "21600"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "tag_name", "requests_total"),
					resource.TestCheckResourceAttr(resourceName, "long_name", "test_alert"),
					resource.TestCheckResourceAttr(resourceName, "interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "skip_notifications", "true"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "6300"),
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
