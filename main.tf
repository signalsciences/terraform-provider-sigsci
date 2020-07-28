provider "sigsci" {
  //  email = ""  //may also provide via env variable
  //  auth_token = "" //may also provide via env variable
  //  password = ""
  corp = "jhanrahan_test_corp"
}

resource "sigsci_site" "my-site" {
  short_name             = "manual_test"
  display_name           = "manual terraform test"
  block_duration_seconds = 1000
  block_http_code        = 303
  agent_anon_mode        = ""
}

resource "sigsci_site_list" "test_list" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new list 2"
  type            = "ip"
  description     = "Some IPs we are putting in a list"
  entries = [
    "4.5.6.7",
    "2.3.4.5",
    "1.2.3.4",
  ]
}

resource "sigsci_site_signal_tag" "test_tag" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new list"
  description     = "description"
}

resource "sigsci_site_redaction" "test_redaction" {
  site_short_name = sigsci_site.my-site.short_name
  field           = "redacted_field"
  redaction_type  = 1
}

resource "sigsci_site_alert" "test_site_alert" {
  site_short_name = sigsci_site.my-site.short_name
  tag_name        = sigsci_site_signal_tag.test_tag.id
  long_name       = "test_alert"
  interval        = 10
  threshold       = 12
  enabled         = true
  action          = "info"
}

resource "sigsci_site_templated_rule" "test_template_rule" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "LOGINATTEMPT"
  detections {
    enabled = "true"
    fields {
      value = "awefwefa"
    }
  }

  alerts {
    long_name          = "awefawef"
    interval           = 60
    threshold          = 10
    skip_notifications = true
    enabled            = true
    action             = "info"
  }

  alerts {
    long_name          = "fwaasd"
    interval           = 60
    threshold          = 1
    skip_notifications = false
    enabled            = false
    action             = "info"
  }
}