resource "sigsci_site" "my-site" {
  short_name             = "manual_test"
  display_name           = "manual terraform test"
  block_duration_seconds = 86400
  agent_anon_mode        = ""
  agent_level            = "block"
  client_ip_rules        = ["X-Client-IP"]
}
