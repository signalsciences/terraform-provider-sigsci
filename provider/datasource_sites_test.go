package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSigSciDataSourceSites(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
        resource "sigsci_site" "example1" {
          short_name   = "terraform_site1"
          display_name = "terraform test example site1"
        }
        resource "sigsci_site" "example2" {
          short_name   = "terraform_site2"
          display_name = "terraform test example site2"
        }
        data "sigsci_sites" "example" {
          depends_on = [sigsci_site.example1, sigsci_site.example2]
          filter = "terraform_site2"
        }
        `,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						r := s.RootModule().Resources["data.sigsci_sites.example"]
						a := r.Primary.Attributes

						sites, err := strconv.Atoi(a["sites.#"])
						if err != nil {
							return err
						}

						if sites != 1 {
							return fmt.Errorf("expected one site to be returned as per the filter")
						}

						got := a["sites.0.name"]
						want := "terraform_site2"
						if got != want {
							return fmt.Errorf("got: %s, want: %s", got, want)
						}

						return nil
					},
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
