### Example Usage

```hcl-terraform
resource "sigsci_corp_cloudwaf_certificate" "test_corp_cloudwaf_certificate" {
    name = "Test Cloud WAF Certificate"
    certificate_body = <<CERT
-----BEGIN CERTIFICATE-----
[encoded certificate]
-----END CERTIFICATE-----
CERT
    certificate_chain = <<CHAIN
-----BEGIN CERTIFICATE-----
[encoded certificate chain]
-----END CERTIFICATE-----
CHAIN
    private_key = <<PRIVATEKEY
-----BEGIN PRIVATE KEY-----
[encoded privatekey]]
----END PRIVATE KEY-----
PRIVATEKEY
}
```

### Argument Reference
- `name` - (Required) Friendly name to identify a CloudWAF certificate.
- `certificate_body` - (Required) Body of the certificate in PEM format.
- `certificate_chain` - (Optional) Certificate chain in PEM format.
- `private_key` - (Required) Private key of the certificate in PEM format - must be unencrypted.

### Attributes Reference
In addition to all arguments, the following fields are also available
- `id` - CloudWAF certificate unique identifier.
- `common_name` - Common name of the uploaded certificate.
- `expires_at` - TimeStamp for when certificate expires in RFC3339 date time format.
- `fingerprint` - SHA1 fingerprint of the certififcate.
- `status` - Current status of the certificate - could be one of "unknown", "active", "pendingverification", "expired", "error".
- `subject_alternative_names` - Subject alternative names from the uploaded certificate.

### Import
You can import corp lists with the generic site import formula

Example:
```shell script
$ terraform import sigsci_corp_cloudwaf_certificate.test id
```
