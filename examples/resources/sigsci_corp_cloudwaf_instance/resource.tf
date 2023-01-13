resource "sigsci_corp_cloudwaf_instance" "test_corp_cloudwaf" {
  name                      = "Test CloudWAF"
  description               = "for test"
  region                    = "ap-northeast-1"
  tls_min_version           = "1.2"
  use_uploaded_certificates = true

  workspace_configs {
    site_name         = sigsci_site.this.short_name
    instance_location = "direct"
    listener_protocols = [
      "https",
    ]

    routes {
      certificate_ids = [
        "a01bc234-5678-9de0-a12b-3456c789d12d",
      ]
      connection_pooling = true
      domains = [
        "example.com",
      ]
      origin              = "https://origin.example.com"
      pass_host_header    = true
      trust_proxy_headers = false
    }
  }
}
