resource "sigsci_site_rule" "test-request-rule" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "request"
  group_operator  = "all"
  enabled         = true
  reason          = "Example request site rule"
  requestlogging  = "sampled"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "ip"
    operator = "equals"
    value    = "1.2.3.4"
  }

  conditions {
    type           = "multival"
    field          = "requestHeader"
    operator       = "exists"
    group_operator = "all"

    conditions {
      type     = "single"
      field    = "name"
      operator = "equals"
      value    = "Content-Type"
    }

    conditions {
      type     = "single"
      field    = "valueString"
      operator = "equals"
      value    = "application/json"
    }
  }

  actions {
    type = "block"
  }
}

resource "sigsci_site_rule" "test-ratelimit-rule" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "rateLimit"
  group_operator  = "all"
  enabled         = true
  reason          = "Example rate limit rule"
  signal          = "site.signal_id"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "ip"
    operator = "equals"
    value    = "1.2.3.5"
  }

  rate_limit = {
    threshold = 6
    interval  = 10
    duration  = 300
  }

  actions {
    type   = "logRequest"
    signal = "site.signal_id"
  }
}


resource "sigsci_site_rule" "test-signal-exclusion" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "signal"
  group_operator  = "all"
  enabled         = true
  reason          = "Example signal exclusion site rule"
  signal          = "SQLI"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "ip"
    operator = "equals"
    value    = "1.2.3.6"
  }

  actions {
    type = "excludeSignal"
  }
}
