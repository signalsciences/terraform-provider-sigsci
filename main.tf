provider "sigsci" {
  email = "jhanrahan+staff+corp2@signalsciences.com"
  //  auth_token = "" //provide via env variable
  //  password = ""
  corp = "jhanrahan_test_corp"
}

resource "sigsci_site" "my-site" {
  short_name             = "test"
  display_name           = "testt"
  block_duration_seconds = 1000
  block_http_code        = 303
  agent_anon_mode        = ""
}

resource "sigsci_site_list" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new list"
  type            = "ip"
  description     = "Some IPs we are putting in a list"
  entries = [
    "4.5.6.7",
    "2.3.4.5",
    "1.2.3.4",
  ]
}
