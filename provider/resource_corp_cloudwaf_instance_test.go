package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceCorpCloudWAFInstanceCRUD(t *testing.T) {
	//t.Parallel() //TODO figure out why we can't run this in parallel
	resourceName := "sigsci_corp_cloudwaf_instance.test_cloudwaf_instance"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCorpCloudWAFInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "sigsci_corp_cloudwaf_instance" "test_cloudwaf_instance"{
					name = "Cloud WAF created by SigSci Terraform provider test"
					description = "Test CWAF Created by SigSci Terraform provider"
					region = "us-west-1"
					tls_min_version = "1.2"
					use_uploaded_certificates=false

					workspace_configs {
						site_name = "%s"
						instance_location = "advanced"
						client_ip_header = "Fastly-Client-IP"
						listener_protocols = ["https"]

						routes {
							domains = ["example.net"]
							origin = "https://example.com"
							connection_pooling = null
							pass_host_header = false
							trust_proxy_headers = true
						}
					}
				}`, testSite),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Cloud WAF created by SigSci Terraform provider test"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test CWAF Created by SigSci Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-west-1"),
					resource.TestCheckResourceAttr(resourceName, "tls_min_version", "1.2"),
					resource.TestCheckResourceAttr(resourceName, "use_uploaded_certificates", "false"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.client_ip_header", "Fastly-Client-IP"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.instance_location", "advanced"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.listener_protocols.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.listener_protocols.1552086545", "https"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.routes.3172309932.certificate_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.routes.3172309932.connection_pooling", "true"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.routes.3172309932.domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.routes.3172309932.domains.3053388764", "example.net"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.routes.3172309932.origin", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.routes.3172309932.pass_host_header", "false"),
					resource.TestCheckResourceAttr(resourceName, "workspace_configs.368698502.routes.3172309932.trust_proxy_headers", "true"),
					resource.TestCheckResourceAttr(resourceName, "deployment.0.status", "done"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment.0.dns_entry"),
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

func testAccCheckCorpCloudWAFInstanceDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client

	resourceType := "sigsci_corp_cloudwaf_instance"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		readResp, err := sc.GetCloudWAFInstance(pm.Corp, resource.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("%s %#v still exists", resourceType, readResp)
		}
	}
	return nil
}
