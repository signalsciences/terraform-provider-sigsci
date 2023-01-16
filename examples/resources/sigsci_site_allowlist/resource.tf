resource "sigsci_site_allowlist" "test" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.2.1"
  note            = "sample allowlist"
}
