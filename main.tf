provider "sigsci" {
  //  email = ""  //may also provide via env variable
  //  auth_token = "" //may also provide via env variable
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

resource "sigsci_site_list" "test_list" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new list 2"
  type            = "ip"
  description     = "Some IPs we are putting in a list"
  entries = [
    "4.5.6.7",
    "2.3.4.5",
    "1.2.3.4",
  ]
}

resource "sigsci_site_signal_tag" "test_tag" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new list 2"
  description     = "descriptionnn"
}