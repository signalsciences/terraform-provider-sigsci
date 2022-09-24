terraform {
  required_providers {
    sigsci = {
      source  = "signalsciences/sigsci"
      version = "1.0.1"
    }
  }
}

// To build locally:
// make && cp terraform-provider-sigsci ~/.terraform.d/plugins/signalsciences/local/sigsci/0.4.2/darwin_amd64/terraform-provider-sigsci && rm .terraform.lock.hcl && tf init
//terraform {
//  required_providers {
//    sigsci = {
//      source  = "signalsciences/local/sigsci"
//      version = "1.0.1"
//    }
//  }
//}

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
    group_operator = "any"
    operator       = "exists"
    conditions {
      field    = "signalType"
      operator = "equals"
      type     = "single"
      value    = "RESPONSESPLIT"
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
        field    = "valueString"
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

resource "sigsci_site_rule" "testsignal" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "templatedSignal"
  group_operator  = "all"
  enabled         = true
  reason          = "Example site rule update"
  signal          = "PW-RESET-ATTEMPT"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "method"
    operator = "equals"
    value    = "POST"
  }

  conditions {
    type     = "single"
    field    = "path"
    operator = "equals"
    value    = "/change-password"
  }

  conditions {
    type           = "multival"
    group_operator = "all"
    conditions {
      field    = "name"
      operator = "equals"
      type     = "single"
      value    = "foo"
    }
  }
}

resource "sigsci_site_integration" "test_integration" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "slack"
  url             = "https://wat.slack.com"
  events          = ["listCreated"]
}

resource "sigsci_corp_cloudwaf_certificate" "test_cloudwaf_certificate"{
  name = "Certificate Name"
  certificate_body = <<CERT
-----BEGIN CERTIFICATE-----
MIIDvDCCAqQCCQDj4MMBbF4gWTANBgkqhkiG9w0BAQsFADCBnzELMAkGA1UEBhMC
VVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28x
FDASBgNVBAoMC0V4YW1wbGUgT3JnMRMwEQYDVQQLDApFeGFtcGxlIE9VMRQwEgYD
VQQDDAtleGFtcGxlLmNvbTEiMCAGCSqGSIb3DQEJARYTZXhhbXBsZUBleGFtcGxl
LmNvbTAeFw0yMjA5MjQwMDAxMzlaFw0yMjEwMjQwMDAxMzlaMIGfMQswCQYDVQQG
EwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNj
bzEUMBIGA1UECgwLRXhhbXBsZSBPcmcxEzARBgNVBAsMCkV4YW1wbGUgT1UxFDAS
BgNVBAMMC2V4YW1wbGUuY29tMSIwIAYJKoZIhvcNAQkBFhNleGFtcGxlQGV4YW1w
bGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnuLmCj29AAyW
fwkErscoHZ4V20DUF3DfsOPUKd5MtXSEQqFnI1fIDkDJC4KL2poTFQoZ4TBuGjcw
lmgAbQggUF4V/UobvEQaqiWeTMUh4YW0szNVvZEVzpwcE2M71KNPx72AuoHs+sgv
FszwAGIpZw1teAhwDqMPscHm/KsK4dxnOkAD+FdMVM5oYCQmf9sPS45FdYEZHueA
554QYObrh43G5tJcte9S9fESgjWfg951ESVcFCHWEG6XQwT9hux9KplgsZJfmgaf
LUaFlnuM8dldi6H9TPL+o4PRRdz8dO/NGD3IkmxncxPt6ATpPRUfgxUi8zr1wEvc
8/oo4C1VVwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQArwwMv9SYSQne12zNEEm7k
77toN9Ya+36mQOFFvNA6Vajd2b4EvlKzbnJox5OkZE6xcE1an4yhKyYYOpqApGr5
mLbdzrUHTqY9IeGrpOuBd2LXrpKtgBR27+lxXzHXd/CWIFn9YVr5IcNaYwCsYgMr
sskUi+lDJWXmkiYoRYKxvR5Ug0NYzEyxj8ZmrGYHk502BxjeW2bFHdXqAZGqoC0O
XjJ53nNvxZHhaIGK9WDDuzem+6b5r4mQVK76BLfwJ/JB2oO3BWYL5cxb6MX/1DXK
ctwO5KlIK0rx/s8nSQBB+QosaXMDP0DQGqEHWT7CQTuT1gNW8ktvzGPQrpjK7JXJ
-----END CERTIFICATE-----
CERT
  private_key = <<PRIVATEKEY
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCe4uYKPb0ADJZ/
CQSuxygdnhXbQNQXcN+w49Qp3ky1dIRCoWcjV8gOQMkLgovamhMVChnhMG4aNzCW
aABtCCBQXhX9Shu8RBqqJZ5MxSHhhbSzM1W9kRXOnBwTYzvUo0/HvYC6gez6yC8W
zPAAYilnDW14CHAOow+xweb8qwrh3Gc6QAP4V0xUzmhgJCZ/2w9LjkV1gRke54Dn
nhBg5uuHjcbm0ly171L18RKCNZ+D3nURJVwUIdYQbpdDBP2G7H0qmWCxkl+aBp8t
RoWWe4zx2V2Lof1M8v6jg9FF3Px0780YPciSbGdzE+3oBOk9FR+DFSLzOvXAS9zz
+ijgLVVXAgMBAAECggEBAImggSLd15jzTmk7ppK+cEE3bjc9MHodi6XtsxmRNWD4
TJhqtqwmnWO7Omp96iawz1aqKUCmcrjClZOzAqtvHo5+8Q015FBvrak0bKqTF4YC
C0Quc1aBFiKhlrA0hN7rl2+s9pSXdm7EeAWH/1xVqwdY2jnfFTGYjT+sdijm/8Yj
hpvlcyqUC3jGrc9hQHDhqwzlVrP/dhYpQIdGwTRHpUXDqwNuXJQIzLhcE0ex+5MM
gWrgAi0M/Qwbn7CyEwdcapDjX6Bt8dVloroeODEYrsDClAXB+45BnRj5HsgYSvQJ
Sn3Xqa5sGxpwOPESOuzhX/Y4f/v2iqbVGw+D0AOcvgECgYEA0dNk+MqpylUrG4dY
uKCV+QADMIllQesw7e8hqKBAubmkfynbml8FMTr9e6GFXIS2Ujwsdg5QQM3b0JqS
qDYSk5EXnWy5zy61Zz9LlGoqmLwEI9NEPtUpSctN1VglhRmW2L9XKnEJ7dgbHl+e
AcMPBI2ownETAClHg8qKm5hMsZsCgYEAwdnQTAtM3gInT6ItxmUIWjKmDQP2tVY0
mIYpiPcoPnFe0IJCzwcWuBezdvgK3x0i2UY/HfRu+d93vMoiXt6ij+zab9ewXVHQ
PmytXelIJuQ4tA38L9HZGJW6yIzIaJT6laZiQStjTEN4lDeKNIpY4SEuqGzXaAl0
CHk+DovJ1PUCgYB82Q6kZloe1QxgRelJeeuijBpZv/brARlNCdN6NVgt6kLxkyNi
uBUr1NDMxi/G/ARL7Bf8asnftV2MwtxukDX/bf6iIfZxS3aOp3++IGmWFZFVC7j4
tfbqPLjkL52rk61I7Jjd3QKubb69FOG8ZKbD69I1V/iZSPaPeW195WIE7wKBgFN0
6deDWfGOrcv7/4cVgjYK7jBWT4WceoJb6E/eUIYpmu9b1VV6MM7K7Wm/ujZ6PcGb
G5tS2+BZ1BwETi3X3dbm2tgh3P0gNu5ZLX5r67NKuBrUlokj6DpMZCDpc3KLCSMa
gdyayGJR/fyZuLeMBF3QQl0ils5km372a8ApcJhtAoGAfDup6P3pDp2Rjq3IRJod
NFJAJERVDnSYml2D9q/uitU51nqL1JcvAoWQJByMFOXuy8F97Olflw4MW/cjWnXJ
1IBEgPEGv5zlronLACOIynBo6/Iu8ULIZhx0BYDc4DXD++PxG2Dd0mbYD8vrlrUK
wOmYr1etSwJs1p3ESLrOscM=
-----END PRIVATE KEY-----
PRIVATEKEY
}
