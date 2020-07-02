provider "sigsci" {
  email = "jhanrahan+staff+corp2@signalsciences.com"
//  auth_token = "" //provide via env variable
//  password = ""
  corp = "jhanrahan_test_corp"
}

resource "sigsci_site" "my-site" {
  short_name = "wattt"
  display_name = "testt"
  block_duration_seconds = 1000
  block_http_code = 303
  agent_anon_mode = ""
}
