package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSigSciDataSourceSite(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
        resource "sigsci_site" "example" {
          short_name   = "terraform_site_ds"
          display_name = "terraform test example site (data source)"
        }
        data "sigsci_site" "example" {
          depends_on = [sigsci_site.example]
          short_name = "terraform_site_ds"
        }
        `,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						r := s.RootModule().Resources["data.sigsci_site.example"]
						a := r.Primary.Attributes

						if got, want := r.Primary.ID, "terraform_site_ds"; got != want {
							return fmt.Errorf("id: got %q, want %q", got, want)
						}

						if got, want := a["short_name"], "terraform_site_ds"; got != want {
							return fmt.Errorf("short_name: got %q, want %q", got, want)
						}

						if got, want := a["display_name"], "terraform test example site (data source)"; got != want {
							return fmt.Errorf("display_name: got %q, want %q", got, want)
						}

						return nil
					},
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
