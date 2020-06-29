provider "sigsci" {
  email = "jhanrahan+staff+corp2@signalsciences.com"
  auth_token = "bc7feacb-4d9b-4a4f-93e2-bd098ea0f163"
//  password = "Hickory290"
  corp = "jhanrahan_test_corp"
}

resource "sigsci_site" "my-site" {
  short_name = "wattt"
  display_name = "testt"
  block_duration_seconds = 1
  block_http_code = 303
//  agent_anon_mode = "off"
}
