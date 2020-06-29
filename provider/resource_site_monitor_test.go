package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccResourceSiteMonitorCRUD(t *testing.T) {
	resourceName := "sigsci_site_monitor.test_monitor"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_monitor" "test_monitor"{
                        site_short_name = "%s"
                        dashboard_id = "000000000000000000000001"
                        share = false
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "share", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_monitor" "test_monitor"{
                        site_short_name = "%s"
                        dashboard_id="000000000000000000000001"
                        share = true
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "share", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
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
