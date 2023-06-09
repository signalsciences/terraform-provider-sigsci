resource "sigsci_site" "my-site" {
  short_name             = "manual_test"
  display_name           = "manual terraform test"
  block_duration_seconds = 86400
  agent_anon_mode        = ""
  agent_level            = "block"
}

resource "sigsci_site_signal_tag" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new signal tag"
  description     = "description of tag"
}
