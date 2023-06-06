package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestACCResourceSiteAllowlist_basic(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site_allowlist.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckSiteAllowlistDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "sigsci_site_allowlist" "test"{
						site_short_name = "%s"
						source          = "1.2.3.4"
						note            = "sample allowlist"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "source", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "note", "sample allowlist"),
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

func testACCCheckSiteAllowlistDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client

	resourceType := "sigsci_site_allowlist"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		allowlistIPs, err := sc.ListWhitelistIPs(pm.Corp, resource.Primary.Attributes["site_short_name"])
		if err != nil {
			return fmt.Errorf("%s couldn't check allowlist ips", resourceType)
		}

		for _, w := range allowlistIPs {
			if w.ID == resource.Primary.ID {
				return fmt.Errorf("%s %#v still exists", resourceType, allowlistIPs)
			}
		}

	}
	return nil
}
