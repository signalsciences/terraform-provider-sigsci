package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TODO implement sweepers for everyone
func TestAccResourceSiteListCRUD(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site_list.test_list"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_list" "test_list"{
                        site_short_name = "%s"
                        name = "test list"
                        type            = "ip"
                        description     = "Some IPs we are putting in a list"
                        entries = [
                          "4.5.6.7",
                          "2.3.4.5",
                          "1.2.3.4",
                        ]
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "test list"),
					resource.TestCheckResourceAttr(resourceName, "type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "description", "Some IPs we are putting in a list"),
					resource.TestCheckResourceAttr(resourceName, "entries.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "entries.1592319998", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "entries.2683765312", "2.3.4.5"),
					resource.TestCheckResourceAttr(resourceName, "entries.402539219", "4.5.6.7"),
				),
			},
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_list" "test_list"{
                        site_short_name = "%s"
                        name = "test list"
                        type            = "ip"
                        description     = "Some IPs we are putting in a list"
                        entries = [
                          "4.5.6.7",
                          "7.8.9.0",
                          "1.2.3.4",
                        ]
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "test list"),
					resource.TestCheckResourceAttr(resourceName, "type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "description", "Some IPs we are putting in a list"),
					resource.TestCheckResourceAttr(resourceName, "entries.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "entries.1592319998", "1.2.3.4"),
					resource.TestCheckNoResourceAttr(resourceName, "entries.2683765312"),
					resource.TestCheckResourceAttr(resourceName, "entries.852349055", "7.8.9.0"),
					resource.TestCheckResourceAttr(resourceName, "entries.402539219", "4.5.6.7"),
				),
			},
		},
	})
}
