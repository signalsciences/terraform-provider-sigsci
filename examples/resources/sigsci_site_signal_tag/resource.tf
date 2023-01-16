resource "sigsci_site_signal_tag" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new signal tag"
  description     = "description of tag"
}
