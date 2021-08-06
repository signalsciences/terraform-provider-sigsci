terraform {
  required_providers {
    sigsci = {
      source  = "signalsciences/local/sigsci"
      version = "0.4.2"
    }
  }
}

provider "sigsci" {
  //  corp = ""       // Required. may also provide via env variable SIGSCI_CORP
  //  email = ""      // Required. may also provide via env variable SIGSCI_EMAIL
  //  auth_token = "" //may also provide via env variable SIGSCI_TOKEN
  //  password = ""   //may also provide via env variable SIGSCI_PASSWORD
}

############# Corp Level Resources #############

resource "sigsci_site" "my-site" {
  short_name             = "manual_test"
  display_name           = "manual terraform test"
  block_duration_seconds = 86400
  agent_anon_mode        = ""
  agent_level            = "block"
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
    type   = "excludeSignal"
    signal = sigsci_corp_signal_tag.test.id
  }
}

resource "sigsci_corp_signal_tag" "test" {
  short_name  = "example-signal-tag"
  description = "An example of a custom signal tag"
}

resource "sigsci_corp_integration" "test_corp_integration" {
  type   = "slack"
  url    = "https://wat.slack.com"
  events = ["newSite", "enableSSO"]
}

############# Site Level Resources #############

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


resource "sigsci_site_signal_tag" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "test"
  description     = "test 2"
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
  name            = "AWS-SSRF"
  detections {
    enabled = "true"
  }

  alerts {
    long_name              = ""
    interval               = 0
    threshold              = 0
    skip_notifications     = false
    enabled                = true
    action                 = "blockImmediate"
    block_duration_seconds = 54321
  }

  alerts {
    long_name              = "AWS-SSRF-10-in-1"
    interval               = 10
    threshold              = 1
    skip_notifications     = false
    enabled                = true
    action                 = "info"
    block_duration_seconds = 54321
  }

  alerts {
    long_name              = "AWS-SSRF-11-in-60"
    interval               = 60
    threshold              = 11
    skip_notifications     = false
    enabled                = true
    action                 = "template"
    block_duration_seconds = 54321
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

resource "sigsci_site_blocklist" "test" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.3.4"
  note            = "sample blocklist"
}

resource "sigsci_site_header_link" "test_header_link" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "test_header_link"
  type            = "request"
  link_name       = "signal sciences 89"
  link            = "https://www.signalsciences.net"
}

resource "sigsci_site_allowlist" "test" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.2.1"
  note            = "sample allowlistt"
}

resource "sigsci_site_redaction" "test_redaction" {
  site_short_name = sigsci_site.my-site.short_name
  field           = "redacted"
  redaction_type  = 0
}


resource "sigsci_site_rule" "testt" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "request"
  group_operator  = "all"
  enabled         = true
  reason          = "Example site rule update"
  expiration      = ""

  conditions {
    type           = "multival"
    field          = "signal"
    group_operator = "all"
    operator       = "exists"
    conditions {
      field    = "signalType"
      operator = "equals"
      type     = "single"
      value    = "RESPONSESPLIT2"
    }
  }

  conditions {
    type           = "group"
    group_operator = "any"
    conditions {
      field    = "useragent"
      operator = "like"
      type     = "single"
      value    = "python-requests*"
    }

    conditions {
      type           = "multival"
      field          = "requestHeader"
      operator       = "doesNotExist"
      group_operator = "all"
      conditions {
        field    = "name"
        operator = "equals"
        type     = "single"
        value    = "cookie"
      }
    }

    conditions {
      type           = "multival"
      field          = "signal"
      operator       = "exists"
      group_operator = "any"
      conditions {
        field    = "signalType"
        operator = "equals"
        type     = "single"
        value    = "TORNODE"
      }
      conditions {
        field    = "signalType"
        operator = "equals"
        type     = "single"
        value    = "SIGSCI-IP"
      }
      conditions {
        field    = "signalType"
        operator = "equals"
        type     = "single"
        value    = "SCANNER"
      }
    }
  }

  actions {
    type = "block"
  }
}

resource "sigsci_site_integration" "test_integration" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "slack"
  url             = "https://wat.slack.com"
  events          = ["listCreated"]
}
