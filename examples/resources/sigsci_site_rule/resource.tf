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

resource "sigsci_site_rule" "test-ratelimit-rule-conditions" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "rateLimit"
  group_operator  = "all"
  enabled         = true
  reason          = "Example rate limit rule that rate limits clients who match the rule conditions after exceeding threshold"
  signal          = "site.count-ratelimit-rule1"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "path"
    operator = "equals"
    value    = "/login"
  }

  rate_limit {
    threshold = 6 
    interval  = 10
    duration  = 300

    client_identifiers {
      type = "ip"
    }
  }

  actions {
    type   = "logRequest"
    signal = "site.count-ratelimit-rule1"
  }
}

resource "sigsci_site_rule" "test-ratelimit-other-signal" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "rateLimit"
  group_operator  = "all"
  enabled         = true
  reason          = "Example rate limit rule that rate limits clients who match a different signal after exceeding threshold"
  signal          = "site.count-ratelimit-rule2"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "path"
    operator = "equals"
    value    = "/reset_password"
  }

  rate_limit = {
    threshold = 6
    interval  = 10
    duration  = 300
  }

  actions {
    type   = "logRequest"
    signal = "site.action-on-other-signal"
  }
}

resource "sigsci_site_rule" "test-ratelimit-all-requests" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "rateLimit"
  group_operator  = "all"
  enabled         = true
  reason          = "Example rule that rate limits all requests from clients after exceeding threshold"
  signal          = "site.count-ratelimit-rule3"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "path"
    operator = "equals"
    value    = "/signup"
  }

  rate_limit = {
    threshold = 6
    interval  = 10
    duration  = 300
  }

  actions {
    type   = "logRequest"
    signal = "ALL-REQUESTS"
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
