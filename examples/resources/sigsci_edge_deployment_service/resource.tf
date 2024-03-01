resource "sigsci_edge_deployment_service" "my-service" {
  site_short_name  = "manual_test"
  fastly_sid       = "test_sid"
  activate_version = true
  percent_enabled  = 100
}
