resource "sigsci_site_list" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new list"
  type            = "ip"
  description     = "Some IPs"
  entries = [
    "4.5.6.7",
    "2.3.4.5",
    "1.2.3.4",
  ]
}
