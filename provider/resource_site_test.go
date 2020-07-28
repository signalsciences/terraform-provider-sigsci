package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/signalsciences/go-sigsci"
	"testing"
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
			if err != nil {
				return err
			}
			return nil
		},
	})
}

func TestAccResourceSiteBasic(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site.test"
	testSite := randStringRunes(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site" "test"{
                        short_name = "%s"
                        display_name = "test 2"
                        block_duration_seconds = 86400
                        block_http_code = 406
                        agent_anon_mode = ""
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test 2"),
					resource.TestCheckResourceAttr(resourceName, "block_duration_seconds", "86400"), // TODO change these values once api is fixed
					resource.TestCheckResourceAttr(resourceName, "block_http_code", "406"),          // TODO change these values once api is fixed
					resource.TestCheckResourceAttr(resourceName, "agent_anon_mode", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site" "test"{
                        short_name = "%s"
                        display_name = "test"
                        block_duration_seconds = 86401
                        block_http_code = 406
                        agent_anon_mode = "EU"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "block_duration_seconds", "86401"),
					resource.TestCheckResourceAttr(resourceName, "block_http_code", "406"), // TODO change these values once api is fixed
					resource.TestCheckResourceAttr(resourceName, "agent_anon_mode", "EU"),
				),
			},
		},
	})
}
