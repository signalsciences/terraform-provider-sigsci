package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestACCResourceSiteBlacklist_basic(t *testing.T) {
	resourceName := "sigsci_site_blacklist.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckSiteBlacklistDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "sigsci_site_blacklist" "test"{
						site_short_name = "%s"
						source          = "1.2.3.4"
						note            = "sample blacklist"
						//expires         = ""
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "source", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "note", "sample blacklist"),
					//resource.TestCheckResourceAttr(resourceName, "expires", ""),
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

func testACCCheckSiteBlacklistDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client

	resourceType := "sigsci_site_blacklist"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		blacklistIPs, err := sc.ListBlacklistIPs(pm.Corp, resource.Primary.Attributes["site_short_name"])
		if err != nil {
			return fmt.Errorf("%s couldn't check Blacklist ips", resourceType)
		}

		for _, w := range blacklistIPs {
			if w.ID == resource.Primary.ID {
				return fmt.Errorf("%s %#v still exists", resourceType, blacklistIPs)
			}
		}

	}
	return nil
}
