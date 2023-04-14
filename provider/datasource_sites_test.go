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
				Config: `data "sigsci_sites" "example" {}`,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						r := s.RootModule().Resources["data.sigsci_sites.example"]
						a := r.Primary.Attributes

						sites, err := strconv.Atoi(a["sites.#"])
						if err != nil {
							return err
						}

						if sites < 1 {
							return fmt.Errorf("expected at least one site to be returned")
						}

						return nil
					},
				),
			},
		},
	})
}
