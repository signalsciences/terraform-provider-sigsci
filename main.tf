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
MIICzDCCAbQCCQDV2NzCr6aPbDANBgkqhkiG9w0BAQsFADAoMQwwCgYDVQQLDAN3
ZWIxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0yMjA5MDMxNzE5MTVaFw0z
MjA4MzExNzE5MTVaMCgxDDAKBgNVBAsMA3dlYjEYMBYGA1UEAwwPd3d3LmV4YW1w
bGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3xLllGg3Kl0S
W16pOwH4VXyGTyByWk3gShoSnEXqWYwoGDF18YhGFFMYLUBNfbDG/jy8MLiKY20R
BhV5hpObd4ZQq4PIlTl4ZNKy07CUPX/AufdbQzrFCQy96lXBVjo6gR10TD+F/CjC
tOkM83dxtZoSPzH86eHteos41+apjgpfvVai3vkBNuZeeoxuERkxuGsfpcK2qWTg
ZFuncrWt6Plvlu70qGEIPtiFiPfQ8Rs2mdzKJEBC8nb4nqSWxIbY9Z87yS3X3C/A
0xr0W4YxOLyCN94qr+Cc3Zl6DvjOv3LWAfv4qFXApWD9f8ynAzjojqfnXtavV6+D
1SOdvMcTVQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQDG9lhmTU9EzE/B7hAhNPj2
LHlRPj8ASnSpxBzWn6Acyu9hHJbGhkR8BnRPPjihH8kv+zRaRVYxhG2sb99qg168
rPKWMbZI6ZvCQGKNjLpwUARwPOKeZ8zF+qyzxdpM9mMyzx9SI1QXDirA0BsUbAjm
RfioqCdT54F8gFrH1+AnUX4Kf2euTS65bHRgegDiIsrAmwcRrzC8ev1SDiPMUkyC
goD1A0LHXLN1LMTs6qBXIkbCjYNRkPZBRagEu68CkwjT5H4vBIl39+Lcvo2WCBSG
j4LzHGcDief95tMLhz0f5g0geV3ytrld5NSw0g1sEYJlDe8NA/aDi4gVviOt3Z5A
-----END CERTIFICATE-----
CERT
  private_key = <<PRIVATEKEY
-----BEGIN PRIVATE KEY-----
MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQDfEuWUaDcqXRJb
Xqk7AfhVfIZPIHJaTeBKGhKcRepZjCgYMXXxiEYUUxgtQE19sMb+PLwwuIpjbREG
FXmGk5t3hlCrg8iVOXhk0rLTsJQ9f8C591tDOsUJDL3qVcFWOjqBHXRMP4X8KMK0
6Qzzd3G1mhI/Mfzp4e16izjX5qmOCl+9VqLe+QE25l56jG4RGTG4ax+lwrapZOBk
W6dyta3o+W+W7vSoYQg+2IWI99DxGzaZ3MokQELydviepJbEhtj1nzvJLdfcL8DT
GvRbhjE4vII33iqv4JzdmXoO+M6/ctYB+/ioVcClYP1/zKcDOOiOp+de1q9Xr4PV
I528xxNVAgMBAAECggEBAKPsP/aJipg/4oBwFE2/SdyP8CZvQnjnpzzs4eYiXm7F
VqVIm1INAOpokWiXSxpk8CXdPbFTuqYLfKoK182z5Fe1xMv0wE4f+D+msTBsHtL+
cQJ3KYJCyo225kwwDi2uBlXg7hglyfCdh07nvtOeX1nCyUvVEPRRSHB3pCLLZqdv
ysCCL4Sowuebcpec3w3nCuMTg+L1nxdk25C53EjsYGqMQgq0YX/CTo+M2X2Infac
3Ig5bkQohaOz724L7mc63UaT9m36vgEUZfwndxxVUBHzxQ/tqr/O4XEKmSdFrExh
yw+YFyP43WcE0m1lc2hYHKEwTM1QF1nVY60fyPJ0zCECgYEA9ljzLd8OjXKbhK7S
HZWZyR1+Lo4HnrBsLgCJVIuGuWi6gDi8hsarYpTe8RVhgqBaudMpK+Eup76WuMBt
F9RVpzkNBE7T+p1jkmCWOIkCpvf2/IoqlZ1ao5qu1fhK5HpnsYCTmUHhJmid3da9
eT34oR2WnvJ9s6fvqC1C6URtsE0CgYEA59B+pvJxvV8klPDfIgN31bcR90srHJTJ
ST3lnVMJNND6PJ3D96YJSHV0EIotqyXlv/droqerhNJ+gNOHvrXMKTJ7udCDQKfW
6IFiWaRW1SkqUEXlv59JI6Ip3B9Rfi3w3rPhNLAi/7GjSm33C+OgsMVq9ZBTMjR6
qnS872CwsykCgYEAq22+3D8K+3ezraOSaDAA8qlpc7A2sUGIJoMNDh6CRGgS0MOq
vgdmoJWEhzQfxS0dtY6yaeyr8ON6M1sFD74dVN8opcTNUutPrT81imYdyF9qKtdj
RvZXat5rqE6+nzxnCGi3TcFAkt/ea8/RzptHd6cFd9q7itfkuJ22oGmUA0kCgYEA
zmvUO+kr6wtr0czjhLA952rLbr/ateqviq65ZmxoiEWGbq+1rzKElacxIQFKRVrL
yTMS/5X6n52o1CKIgAP2tsCjeAT6u3o5XnTIFTbHs6yiZzS2rvmx8S8Xw1GICanz
EPxwj7BAmhueYkqlcErT7lT9N4m667vbdynYi/g3oHECgYEA5dCjPFl94ZGStHF2
4+cI1NCbVTnQApFI1+Vmd4lED619KSnpk/77TaYTh2I8gsK3OZP4crvee5aL41TJ
Y20XmAq/Dvv3g7QR97ND/AghEU8nnpZo1fgHzCcyZSkaBzMFeYdsfaGDncM56I0B
65yblwYq9Vyzy3hBFY6XGaFdnZ0=
-----END PRIVATE KEY-----
PRIVATEKEY
}
