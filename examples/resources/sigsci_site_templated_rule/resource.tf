resource "sigsci_site_templated_rule" "test_template_rule" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "LOGINATTEMPT"
  detections {
    enabled = "true"
    fields {
      name  = "path"
      value = "/login/*"
    }
  }

  alerts {
    long_name              = "alert 1"
    interval               = 60
    threshold              = 10
    skip_notifications     = true
    enabled                = true
    action                 = "info"
    block_duration_seconds = sigsci_site.my-site.block_duration_seconds
  }

  alerts {
    long_name              = "alert 2"
    interval               = 60
    threshold              = 1
    skip_notifications     = false
    enabled                = false
    action                 = "info"
    block_duration_seconds = 64000
  }
}
