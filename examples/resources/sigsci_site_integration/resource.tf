resource "sigsci_site_integration" "test_integration" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "slack"
  url             = "https://wat.slack.com"
  events          = ["listCreated"]
}
