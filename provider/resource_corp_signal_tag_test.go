package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestACCResourceCorpSignalTag_basic(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_corp_signal_tag.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckCorpSignalTagDestroy,

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprint(`
					resource "sigsci_corp_signal_tag" "test"{
						short_name="testacc-signal-tagg"
						description= "An example of a custom signal tag"
				}`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "configurable"),
					resource.TestCheckResourceAttrSet(resourceName, "informational"),
					resource.TestCheckResourceAttrSet(resourceName, "needs_response"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "testacc-signal-tagg"),
					resource.TestCheckResourceAttr(resourceName, "description", "An example of a custom signal tag"),
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

func testACCCheckCorpSignalTagDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client

	resourceType := "sigsci_corp_signal_tag"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		readresp, err := sc.GetCorpSignalTagByID(pm.Corp, resource.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("%s %#v still exists", resourceType, readresp)
		}
	}
	return nil
}
