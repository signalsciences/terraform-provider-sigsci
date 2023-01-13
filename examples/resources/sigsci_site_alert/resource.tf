resource "sigsci_site_alert" "test" {
  site_short_name    = sigsci_site.my-site.short_name
  tag_name           = sigsci_site_signal_tag.test_tag.id
  long_name          = "test_alert"
  interval           = 10
  threshold          = 12
  enabled            = true
  action             = "info"
  skip_notifications = true
}
