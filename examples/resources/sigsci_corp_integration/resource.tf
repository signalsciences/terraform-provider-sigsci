resource "sigsci_corp_integration" "test_corp_integration" {
  type   = "slack"
  url    = "https://signalsciences.slack.com"
  events = ["newSite", "enableSSO"]
}
