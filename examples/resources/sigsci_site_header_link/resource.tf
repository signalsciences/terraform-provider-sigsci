resource "sigsci_site_header_link" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "test_header_link"
  type            = "request"
  link_name       = "signal sciences"
  link            = "https://www.signalsciences.net"
}
