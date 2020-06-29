package provider

//func TestAccResourceSiteIntegrationCRUD(t *testing.T) {
//	resourceName := "sigsci_site_integration.test_integration"
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:  func() { testAccPreCheck(t) },
//		Providers: testAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: fmt.Sprintf(`
//                    resource "sigsci_site_integration" "test_integration"{
//                        site_short_name = "%s"
//                        type = "generic"
//                        url = "https://hooks.slack.com/services/blah/blah"
//                        events = [
//                          "webhookEvents"
//                        ]
//				}`, testSite),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
//					resource.TestCheckResourceAttr(resourceName, "type", "generic"),
//					resource.TestCheckResourceAttr(resourceName, "url", "https://hooks.slack.com/services/blah/blah"),
//					resource.TestCheckResourceAttr(resourceName, "events.#", "1"),
//					resource.TestCheckResourceAttr(resourceName, "events.3935235602", "webhookEvents"),
//				),
//			},
//			{
//				Config: fmt.Sprintf(`
//                    resource "sigsci_site_integration" "test_integration"{
//                        site_short_name = "%s"
//                        type = "generic"
//                        url = "https://hooks.slack.com/services/blah/blah2"
//                        events = [
//                          "flag",
//                          "loggingModeChanged"
//                        ]
//				}`, testSite),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					testInspect(),
//					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
//					resource.TestCheckResourceAttr(resourceName, "type", "generic"),
//					resource.TestCheckResourceAttr(resourceName, "url", "https://hooks.slack.com/services/blah/blah2"),
//					resource.TestCheckResourceAttr(resourceName, "events.#", "2"),
//					resource.TestCheckResourceAttr(resourceName, "events.1929545752", "flag"),
//					resource.TestCheckResourceAttr(resourceName, "events.3593650736", "loggingModeChanged"),
//				),
//			},
//			{
//				Config: fmt.Sprintf(`
//			       resource "sigsci_site_integration" "test_integration"{
//			           site_short_name = "%s"
//			           type = "slack"
//			           url = "https://hooks.slack.com/services/blah/blah3"
//			           events = [
//			           ]
//				}`, testSite),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					testInspect(),
//					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
//					resource.TestCheckResourceAttr(resourceName, "type", "slack"),
//					resource.TestCheckResourceAttr(resourceName, "url", "https://hooks.slack.com/services/blah/blah3"),
//					resource.TestCheckResourceAttr(resourceName, "events.#", "0"),
//				),
//			},
//			{
//				ResourceName:        resourceName,
//				ImportStateIdPrefix: fmt.Sprintf("%s:", testSite),
//				ImportState:         true,
//				ImportStateVerify:   true,
//				ImportStateCheck:    testAccImportStateCheckFunction(1),
//			},
//		},
//	})
//}
