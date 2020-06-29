provider "sigsci" {
  //  corp = "" // Required. may also provide via env variable SIGSCI_CORP
  //  email = ""  // Required. may also provide via env variable SIGSCI_EMAIL
  //  auth_token = "" //may also provide via env variable SIGSCI_TOKEN
  //  password = "" //may also provide via env variable SIGSCI_PASSWORD
}


resource "sigsci_corp_list" "test" {
  name        = "My corp list"
  type        = "ip"
  description = "Some IPs"
  entries = [
    "4.5.6.7",
    "2.3.4.5",
    "1.2.3.4",
  ]
}


resource "sigsci_corp_rule" "test" {
  site_short_names = [sigsci_site.my-site.short_name]
  type             = "signal"
  corp_scope       = "specificSites"
  enabled          = true
  group_operator   = "any"
  signal           = "SQLI"
  reason           = "Example corp rule"
  expiration       = ""

  conditions {
    type     = "single"
    field    = "ip"
    operator = "equals"
    value    = "1.2.3.4"
  }
  conditions {
    type     = "single"
    field    = "ip"
    operator = "equals"
    value    = "1.2.3.5"
  }
  actions {
    type = "excludeSignal"
  }
}

resource "sigsci_corp_signal_tag" "test" {
  short_name  = "example-signal-tag"
  description = "An example of a custom signal tag"
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
  name            = "My new signal tag"
  description     = "description"
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
      name  = "path"
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

resource "sigsci_site_rule" "test" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "signal"
  group_operator  = "any"
  enabled         = true
  reason          = "Example site rule update"
  signal          = "SQLI"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "ip"
    operator = "equals"
    value    = "1.2.3.4"
  }
  conditions {
    type     = "single"
    field    = "ip"
    operator = "equals"
    value    = "1.2.3.5"
    conditions {
      type           = "multival"
      field          = "ip"
      operator       = "equals"
      group_operator = "all"
      value          = "1.2.3.8"
    }
  }

  actions {
    type = "excludeSignal"
  }
}

resource "sigsci_site_blacklist" "test" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.3.4"
  note            = "sample blacklist"
}

resource "sigsci_site_header_link" "test_header_link" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "test_header_link"
  type            = "request"
  link_name       = "signal sciences"
  link            = "https://www.signalsciences.net"
}

resource "sigsci_site_whitelist" "test" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.3.1"
  note            = "sample whitelistt"
}

resource "sigsci_site_whitelist" "test2" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.3.5"
  note            = "sample whitelist"
}

resource "sigsci_site_whitelist" "test3" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.3.3"
  note            = "sample whitelist"
}

resource "sigsci_site_redaction" "test_redaction" {
  site_short_name = sigsci_site.my-site.short_name
  field           = "redacted_f"
  redaction_type  = 2
}

resource "sigsci_site_redaction" "test_redaction2" {
  site_short_name = sigsci_site.my-site.short_name
  field           = "redacted_field"
  redaction_type  = 1
}

resource "sigsci_site_redaction" "test_redaction3" {
  site_short_name = sigsci_site.my-site.short_name
  field           = "redacted_field"
  redaction_type  = 0
}

resource "sigsci_site" "my-site" {
  short_name             = "manual_test"
  display_name           = "manual terraform test"
  block_duration_seconds = 86400 // TODO check this...
  block_http_code        = 406   // default 406, field not really respected
  agent_anon_mode        = ""
  agent_level            = "block"
}
