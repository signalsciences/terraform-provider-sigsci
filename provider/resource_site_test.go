package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

var testSite = "test" //acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

func TestAccResourceSiteBasic(t *testing.T) {
	resourceName := "sigsci_site.test"
	testSite = acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
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
