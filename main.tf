provider "sigsci" {
  auth_email = "jhanrahan+staff+corp2@signalsciences.com"
//  auth_token = "bc7feacb-4d9b-4a4f-93e2-bd098ea0f163"
  auth_password = "Hickory290"
  corp = "jhanrahan_test_corp"
}

resource "sigsci_site" "my-site" {
  name = "wat"
  display_name = "test"
  block_duration_seconds = 1
}
