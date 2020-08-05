package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestResourceCorpRule_basic(t *testing.T) {
	resourceName := "sigsci_corp_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCCheckCorpRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "sigsci_corp_rule" "test"{
					site_short_names=["%s"]
					type= "signal"
					corp_scope="specificSites"
					enabled=true
					group_operator="any"
					signal="SQLI"
					reason="Example corp rule"
					expiration= ""
					
					conditions {
						type="single"
						field="ip"
						operator="equals"
						value="1.2.3.4"
					}
					conditions {
						type="single"
						field="ip"
						operator="equals"
						value="1.2.3.5"
					}
					actions {
						type="excludeSignal"
					}
			}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1859487095.type", "excludeSignal"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.conditions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.field", "ip"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.group_operator", ""),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.operator", "equals"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.type", "single"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2534374319.value", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2383694574.conditions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2383694574.field", "ip"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2383694574.group_operator", ""),
					resource.TestCheckResourceAttr(resourceName, "conditions.2383694574.operator", "equals"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2383694574.type", "single"),
					resource.TestCheckResourceAttr(resourceName, "conditions.2383694574.value", "1.2.3.5"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "expiration", ""),
					resource.TestCheckResourceAttr(resourceName, "group_operator", "any"),
					resource.TestCheckResourceAttr(resourceName, "reason", "Example corp rule"),
					resource.TestCheckResourceAttr(resourceName, "signal", "SQLI"),
					resource.TestCheckResourceAttr(resourceName, "site_short_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "site_short_names.1785148924", testSite),
					resource.TestCheckResourceAttr(resourceName, "type", "signal"),
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

func testACCCheckCorpRuleDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client
	resourceType := "sigsci_corp_rule"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		readResp, err := sc.GetCorpRuleByID(pm.Corp, resource.Primary.ID)
		if err == nil {
			return fmt.Errorf("%s %#v still exists", resourceType, readResp)
		}
	}
	return nil
}
