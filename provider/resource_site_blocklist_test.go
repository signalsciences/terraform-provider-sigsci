package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestACCResourceSiteBlocklist_basic(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site_blocklist.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckSiteBlocklistDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "sigsci_site_blocklist" "test"{
						site_short_name = "%s"
						source          = "1.2.3.4"
						note            = "sample blocklist"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "source", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "note", "sample blocklist"),
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

func testACCCheckSiteBlocklistDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client

	resourceType := "sigsci_site_blocklist"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		blocklistIPs, err := sc.ListBlacklistIPs(pm.Corp, resource.Primary.Attributes["site_short_name"])
		if err != nil {
			return fmt.Errorf("%s couldn't check Blocklist ips", resourceType)
		}

		for _, w := range blocklistIPs {
			if w.ID == resource.Primary.ID {
				return fmt.Errorf("%s %#v still exists", resourceType, blocklistIPs)
			}
		}

	}
	return nil
}
