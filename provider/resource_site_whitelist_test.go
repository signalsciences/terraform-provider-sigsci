package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestACCResourceSiteWhitelist_basic(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site_whitelist.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckSiteWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "sigsci_site_whitelist" "test"{
						site_short_name = "%s"
						source          = "1.2.3.4"
						note            = "sample whitelist"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "source", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "note", "sample whitelist"),
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

func testACCCheckSiteWhitelistDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client

	resourceType := "sigsci_site_whitelist"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		whitelistIPs, err := sc.ListWhitelistIPs(pm.Corp, resource.Primary.Attributes["site_short_name"])
		if err != nil {
			return fmt.Errorf("%s couldn't check whitelist ips", resourceType)
		}

		for _, w := range whitelistIPs {
			if w.ID == resource.Primary.ID {
				return fmt.Errorf("%s %#v still exists", resourceType, whitelistIPs)
			}
		}

	}
	return nil
}
