package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

//TODO implement sweepers for everyone
func TestAccResourceSiteSignalTagCRUD(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site_signal_tag.test_tag"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_signal_tag" "test_tag" {
                      site_short_name = "%s" 
                      name            = "My new tag"
                      description     = "test description"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "My new tag"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
				),
			},
			{
				Config: fmt.Sprintf(`
                     resource "sigsci_site_signal_tag" "test_tag" {
                      site_short_name = "%s" 
                      name            = "My new tag"
                      description     = "test description 2"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "My new tag"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description 2"),
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
