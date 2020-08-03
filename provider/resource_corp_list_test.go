package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestACCResourceCorpList_basic(t *testing.T) {
	resourceName := "sigsci_corp_list.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckCorpListDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprint(`
					resource "sigsci_corp_list" "test"{
						name="My new list"
						type= "ip"
						description= "Some IPs we are putting in a list"
						entries= [
							"4.5.6.7",
							"2.3.4.5",
							"1.2.3.4",
						]
				}`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "My new list"),
					resource.TestCheckResourceAttr(resourceName, "type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "description", "Some IPs we are putting in a list"),
					resource.TestCheckResourceAttr(resourceName, "entries.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "entries.1592319998", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "entries.2683765312", "2.3.4.5"),
					resource.TestCheckResourceAttr(resourceName, "entries.402539219", "4.5.6.7"),
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

func testACCCheckCorpListDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client

	resourceType := "sigsci_corp_list"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		readResp, err := sc.GetCorpListByID(pm.Corp, resource.Primary.ID)
		if err == nil {
			return fmt.Errorf("%s %#v still exists", resourceType, readResp)
		}
	}
	return nil
}
