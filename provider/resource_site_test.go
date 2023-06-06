package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/signalsciences/go-sigsci"
)

func init() {
	resource.AddTestSweepers("site_templated_rules_sweeper", &resource.Sweeper{
		Name: "site_templated_rules_sweeper",
		F: func(region string) error {
			metadata := testAccProvider.Meta().(providerMetadata)
			sc := metadata.Client
			_ = sc.DeleteSite(metadata.Corp, testSite)
			_, err := sc.CreateSite(metadata.Corp, sigsci.CreateSiteBody{
				Name:                 testSite,
				DisplayName:          testSite,
				AgentLevel:           "log",
				AgentAnonMode:        "",
				BlockHTTPCode:        400,
				BlockDurationSeconds: 60,
			})

			return err
		},
	})
}

func TestAccResourceSiteBasic(t *testing.T) {
	resourceName := "sigsci_site.test"
	testSite := randStringRunes(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site" "test"{
                        short_name = "%s"
                        display_name = "test 2"
                        agent_anon_mode = "EU"
                        block_duration_seconds = 86401
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test 2"),
					resource.TestCheckResourceAttr(resourceName, "agent_anon_mode", "EU"),
					resource.TestCheckResourceAttr(resourceName, "block_duration_seconds", "86401"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_agent_key.name"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_agent_key.access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_agent_key.secret_key"),
				),
			},
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site" "test"{
                        short_name = "%s"
                        display_name = "test"
                        agent_anon_mode = ""
                        block_duration_seconds = 86400
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "agent_anon_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "block_duration_seconds", "86400"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccImportStateCheckFunction(1),
			},
		},
	})
}

func testACCCheckSiteDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client

	resourceType := "sigsci_site"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		readResp, err := sc.GetSite(pm.Corp, resource.Primary.ID)
		if err == nil {
			return fmt.Errorf("%s %#v still exists", resourceType, readResp)
		}
	}
	return nil
}
