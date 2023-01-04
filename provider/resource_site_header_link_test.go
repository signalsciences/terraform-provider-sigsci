package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

//TODO implement sweepers for everyone
func TestAccResourceSiteHeaderLinkCRUD(t *testing.T) {
	resourceName := "sigsci_site_header_link.test_header_link"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_header_link" "test_header_link" {
                      site_short_name = "%s"
                      name = "test_header_link"
                      type =  "request"
                      link_name = "signal sciences"
                      link = "https://www.signalsciences.net"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "test_header_link"),
					resource.TestCheckResourceAttr(resourceName, "type", "request"),
					resource.TestCheckResourceAttr(resourceName, "link_name", "signal sciences"),
					resource.TestCheckResourceAttr(resourceName, "link", "https://www.signalsciences.net"),
				),
			},
			{
				Config: fmt.Sprintf(`
                     resource "sigsci_site_header_link" "test_header_link" {
                      site_short_name = "%s"
                      name = "test_header_link"
                      type =  "response"
                      link_name = "signal sciences 2"
                      link = "https://www.signalsciences.com"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "test_header_link"),
					resource.TestCheckResourceAttr(resourceName, "type", "response"),
					resource.TestCheckResourceAttr(resourceName, "link_name", "signal sciences 2"),
					resource.TestCheckResourceAttr(resourceName, "link", "https://www.signalsciences.com"),
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
