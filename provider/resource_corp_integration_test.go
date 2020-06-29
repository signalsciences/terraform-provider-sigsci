package provider

//func TestAccResourceCorpIntegrationCRUD(t *testing.T) {
//	//t.Parallel() //TODO figure out why we can't run this in parallel
//	resourceName := "sigsci_corp_integration.test_integration"
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:  func() { testAccPreCheck(t) },
//		Providers: testAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: fmt.Sprintf(`
//                    resource "sigsci_corp_integration" "test_integration"{
//                        type = "slack"
//                        url = "https://hooks.slack.com/services/blah/blah"
//                        events = [
//                          "webhookEvents"
//                        ]
//				}`),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "type", "slack"),
//					resource.TestCheckResourceAttr(resourceName, "url", "https://hooks.slack.com/services/blah/blah"),
//					resource.TestCheckResourceAttr(resourceName, "events.#", "1"),
//					resource.TestCheckResourceAttr(resourceName, "events.3935235602", "webhookEvents"),
//				),
//			},
//			{
//				Config: fmt.Sprintf(`
//			       resource "sigsci_corp_integration" "test_integration"{
//			           type = "slack"
//			           url = "https://hooks.slack.com/services/blah/blah"
//			           events = [
//			             "corpUpdated",
//			             "listDeleted"
//			           ]
//				}`),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					testInspect(),
//					resource.TestCheckResourceAttr(resourceName, "type", "slack"),
//					resource.TestCheckResourceAttr(resourceName, "url", "https://hooks.slack.com/services/blah/blah"),
//					resource.TestCheckResourceAttr(resourceName, "events.#", "2"),
//					resource.TestCheckResourceAttr(resourceName, "events.963169190", "corpUpdated"),
//					resource.TestCheckResourceAttr(resourceName, "events.2369979095", "listDeleted"),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateCheck:  testAccImportStateCheckFunction(1),
//			},
//		},
//	})
//}
