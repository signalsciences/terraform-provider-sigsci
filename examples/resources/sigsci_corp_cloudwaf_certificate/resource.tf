resource "sigsci_corp_cloudwaf_certificate" "test_corp_cloudwaf_certificate" {
  name              = "Test Cloud WAF Certificate"
  certificate_body  = <<CERT
-----BEGIN CERTIFICATE-----
[encoded certificate]
-----END CERTIFICATE-----
CERT
  certificate_chain = <<CHAIN
-----BEGIN CERTIFICATE-----
[encoded certificate chain]
-----END CERTIFICATE-----
CHAIN
  private_key       = <<PRIVATEKEY
-----BEGIN PRIVATE KEY-----
[encoded privatekey]]
----END PRIVATE KEY-----
PRIVATEKEY
}
