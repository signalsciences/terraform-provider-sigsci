terraform {
  required_providers {
    sigsci = {
      source = "signalsciences/sigsci"
    }
  }
}

// To build locally:
// make && cp terraform-provider-sigsci ~/.terraform.d/plugins/signalsciences/local/sigsci/1.2.1/darwin_amd64/terraform-provider-sigsci && rm .terraform.lock.hcl && tf init

provider "sigsci" {
  //  corp = ""       // Required. may also provide via env variable SIGSCI_CORP
  //  email = ""      // Required. may also provide via env variable SIGSCI_EMAIL
  //  auth_token = "" //may also provide via env variable SIGSCI_TOKEN
  //  password = ""   //may also provide via env variable SIGSCI_PASSWORD
  //  fastly_key = ""  //may also provide via env variable FASTLY_KEY. Required for Edge Deployments functionality.
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
  site_short_name        = sigsci_site.my-site.short_name
  tag_name               = sigsci_site_signal_tag.test_tag.id
  long_name              = "test_alert"
  interval               = 10
  threshold              = 12
  enabled                = true
  action                 = "info"
  block_duration_seconds = 86400
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
  }
  conditions {
    type           = "multival"
    field          = "queryParameter"
    operator       = "exists"
    group_operator = "all"

    conditions {
      type     = "single"
      field    = "name"
      operator = "equals"
      value    = "hello"
    }

    conditions {
      type     = "single"
      field    = "value"
      operator = "equals"
      value    = "world"
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
  reason          = ""
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
    field          = "postParameter"
    operator       = "exists"
    group_operator = "all"
    conditions {
      field    = "name"
      operator = "equals"
      type     = "single"
      value    = "submit"
    }
  }
}

resource "sigsci_site_integration" "test_integration" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "slack"
  url             = "https://wat.slack.com"
  events          = ["listCreated"]
}

resource "sigsci_corp_cloudwaf_certificate" "test_cloudwaf_certificate" {
  name             = "Certificate Name"
  certificate_body = <<CERT
-----BEGIN CERTIFICATE-----
MIIDzjCCArYCCQD6uBPuCbaDuDANBgkqhkiG9w0BAQsFADCBqDELMAkGA1UEBhMC
VVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28x
HTAbBgNVBAoMFEV4YW1wbGUgT3JnYW5pemF0aW9uMRMwEQYDVQQLDApFeGFtcGxl
IE9VMRQwEgYDVQQDDAtleGFtcGxlLmNvbTEiMCAGCSqGSIb3DQEJARYTZXhhbXBs
ZUBleGFtcGxlLmNvbTAeFw0yMjA5MjQwMDE4MjRaFw0zMjA5MjEwMDE4MjRaMIGo
MQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2Fu
IEZyYW5jaXNjbzEdMBsGA1UECgwURXhhbXBsZSBPcmdhbml6YXRpb24xEzARBgNV
BAsMCkV4YW1wbGUgT1UxFDASBgNVBAMMC2V4YW1wbGUuY29tMSIwIAYJKoZIhvcN
AQkBFhNleGFtcGxlQGV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAscvDb2j2s9bdiAIHbqRoM2qZBxdM4atSwAJQrXVe3pbne2KLZw53
kHpVtjaugfMKBnXueR1iilYu5eXtgNfrNHgq0X0+NToL/xtSgYthp89lxBYArUVy
kiM5gy8BqpApfAwQ5MMDgGflIV/mTCcCyNK3DwuOgO7oVp0V2zdtJhgvZ8e3qkuT
3dOxC27aUFYf/P88UILoc9YWRCkw2Gww/Zr908a/mgVBJ9v+/sKP3/yk8jzrRhL5
JsGWC5Gbv1gpkyzSjKyboYePvJJo5D6Fue9XZmzry3wepG1oUcLO6QpH+lTBfTjd
xHKA4sIza1J/RDBLgUBney1nMxLN8RzU5QIDAQABMA0GCSqGSIb3DQEBCwUAA4IB
AQAPRvwDkKTKCDQj5F4ZUTE9AIEs0w99KuXiWBGz3RmYl5zwZCrVWeOI+lPfCG0v
prMgh5ydUgUOqrs8S7MAkt8GaU5lb0MSKmz1jPgEEbLBp6VYv2UbrWlBz9JIxTLw
riPHNUzKb6SXk5wuoO8w7+GsBNI8fWPDSQqSWLlNsi0r4ReLxlM5WBNC10d3q2ia
jV6r8iMpiArwbJn4WSTlFuJ6crrjgbBVCFxxwoF1sHhwGg+5idxm2AHSzvENyFW4
UVhVTn9w3UPLMkEl7nAVzydpdMb/M/GLCV787BrQL35EtiCr9MSL9Gc8vR/9PzPP
QodC+xWXbig7xKLqZgQ/PbPt
-----END CERTIFICATE-----
CERT
  private_key      = <<PRIVATEKEY
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCxy8NvaPaz1t2I
AgdupGgzapkHF0zhq1LAAlCtdV7elud7YotnDneQelW2Nq6B8woGde55HWKKVi7l
5e2A1+s0eCrRfT41Ogv/G1KBi2Gnz2XEFgCtRXKSIzmDLwGqkCl8DBDkwwOAZ+Uh
X+ZMJwLI0rcPC46A7uhWnRXbN20mGC9nx7eqS5Pd07ELbtpQVh/8/zxQguhz1hZE
KTDYbDD9mv3Txr+aBUEn2/7+wo/f/KTyPOtGEvkmwZYLkZu/WCmTLNKMrJuhh4+8
kmjkPoW571dmbOvLfB6kbWhRws7pCkf6VMF9ON3EcoDiwjNrUn9EMEuBQGd7LWcz
Es3xHNTlAgMBAAECggEAV0K/f52Pf0JUZd1BEo+ESL/nrTBFXni8W1qHiCqTzkFY
CRmbe5ABJJq2GIEL8uF6qSMWUMEYTPbxe4n2oAbY/F6B/WEvt+XuX11kiAoFetvy
gWOfH2t3SLwbDQR0F+c7RROS8wO3Yz0amt+7YuK+nhu1FqBAZ41Z4LCmOnogitIA
cCErKpHqCJbYT99eaXTt2QpXJNI8fItXaO4p8zfKxzBKybhsyu3tEerKWnhqz+25
Xr7OieYkZM3ryrIsVWZ299wH2D+gA9O+PbY2RJ0Vf3YVf8VdyYRJ5oCAaioNUfZw
HeGGlfClZzrpX1MjiNNfli05cLqX0iE2bIO+jvZkIQKBgQDadIMHfGQs8CPxbPqi
fZRuHdPedogM80f4E5RyKhTzEqTG6x0pfjkDr57rdqJEWhM3TjsMIiJRnEIomWar
2tyckmvkxuSioiWb/+HXJN5u6AtsMgQ4WLMm9HOO0ir5uSKQd0iQrCMAaGDNdV52
6eipkWLdhYAhW31bycv9cX04+QKBgQDQWlkTovazzNElJU4h1YQCq99VS1gmYb9U
HAzg3Jmu7WF4Oln+HxZMWqwR38vCuMHiwtCmsqGAEKK2ev6W6iq7iVJLKzTXV++a
612Mr+JohbHNL0bKlgTMt/i2TnmBWOlhL7xuIwduru1pQ4mM7Vh7Hv+CrVTGm4VZ
Khzq+vbCTQKBgEnprgO0ZLiHr8GZ29tqnfP8B5l3hWTMU4duKIXQEzKDFllvZ3iI
ioXiv+RvSUvTJjlKMNRUIER4mDHgZUq0THx1Vigb23PjZNI5a5I9mTzxKhw7eA4Q
hN0jTI4AMiY4K6exlE3O0DDtIAOkOIgHcH8e/9JvvwCKUgniZzCjW3kRAoGAeEId
tgLSyFbIxNrybP7zciNIBdA2MfkrWN3T5RoPLnNfVejANrg0w592P97fmiXP6xWt
Hvpt0yBG+nKlbe/8+D+7mx12I3FjIBUH6wM9+Dxqsta90oKihJMPYBKNeUYbdnf6
F8vqJ02aRK6xvwDjmDT9H6zyCKyNXDi9djeio+UCgYEAm1g/2O42eVBzYG/WTJ2H
7jplYemfbIFVpl3Uo18UJKZl/AIzm9tw/+c884naSubwQ8TukI3bDjwWu99R26bo
HQRmLMnP+t7xp84Rn4jXReWlr9sexXHPg/Lj25MdR1t3Ow53qSh4nw/cUPr42N9o
cX4iWLb38v7KEornZfofXEw=
-----END PRIVATE KEY-----
PRIVATEKEY
}

resource "sigsci_edge_deployment" "edge" {
  site_short_name = sigsci_site.my-site.short_name
}

resource "sigsci_edge_deployment_service" "edge" {
  site_short_name  = sigsci_site.my-site.short_name
  fastly_sid       = "[Fastly service id]"
  activate_version = true
  percent_enabled  = 100
}
