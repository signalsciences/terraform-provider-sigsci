package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

//TODO implement sweepers for everyone
func TestAccResourceSiteRedactionCRUD(t *testing.T) {
	resourceName := "sigsci_site_redaction.test_redaction"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_redaction" "test_redaction" {
                      site_short_name    = "%s" 
                      field              = "field"
                      redaction_type     = 0
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "field", "field"),
					resource.TestCheckResourceAttr(resourceName, "redaction_type", "0"),
				),
			},
			{
				Config: fmt.Sprintf(`
                     resource "sigsci_site_redaction" "test_redaction" {
                      site_short_name    = "%s" 
                      field              = "field 2"
                      redaction_type     = 1
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "field", "field 2"),
					resource.TestCheckResourceAttr(resourceName, "redaction_type", "1"),
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
