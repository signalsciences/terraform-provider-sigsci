### Example Usage

```hcl-terraform
resource "sigsci_corp_cloudwaf_instance" "test_corp_cloudwaf" {
    name                      = "Test CloudWAF"
    description               = "for test"
    region                    = "ap-northeast-1"
    tls_min_version           = "1.2"
    use_uploaded_certificates = true

    workspace_configs {
        site_name          = sigsci_site.this.short_name
        instance_location  = "direct"
        listener_protocols = [
            "https",
        ]

        routes {
            certificate_ids     = [
                "a01bc234-5678-9de0-a12b-3456c789d12d",
            ]
            connection_pooling  = true
            domains             = [
                "example.com",
            ]
            origin              = "https://origin.example.com"
            pass_host_header    = true
            trust_proxy_headers = false
        }
    }
}
```

### Argument Reference
- `name` - (Required) Friendly name to identify a CloudWAF instance.
- `description` - (Required) Friendly description to identify a CloudWAF instance.
- `region` - (Required) Region the CloudWAF Instance is being deployed to. See the [documentation](https://docs.fastly.com/signalsciences/api/#_corps__corpName__cloudwafInstances_post) for a list of available regions.
- `tls_min_version` - (Required) TLS minimum version. Versions Available: "1.0", "1.2".
- `use_uploaded_certificates` - (Required) If "true", use the uploaded certificate.
- `workspace_configs` - (Required) Workspace Configs. Detailed below.
  - `site_name` - (Required) Site name.
  - `instance_location` - (Required) Set instance location to "direct" or "advanced".
  - `client_ip_header` - (Optional) Specify the request header containing the client IP address, available when InstanceLocation is set to "advanced". Default: "X-Forwarded-For".
  - `listener_protocols` - (Required) Specify the protocol or protocols required. ex. ["http", "https"], ["https"].
  - `routes` - (Required) Routes. Detailed below.
    - `certificate_ids` - (Optional) List of certificate IDs in string associated with request URI or domains. IDs will be available in certificate GET request.
    - `connection_pooling` - (Optional) If enabled, this will allow open TCP connections to be reused (default: true).
    - `domains` - (Required) List of domain or request URIs, up to 100 entries.
    - `origin` - (Required) Origin server URI.
    - `pass_host_header` - (Optional) Pass the client supplied host header through to the upstream (including the upstream TLS handshake for use with SNI and certificate validation). If using Heroku or Server Name Indications (SNI), this must be disabled (default: false).
    - `trust_proxy_headers` - (Optional) If true, will trust proxy headers coming into the agent. If false, will ignore and drop those headers (default: false).

### Attributes Reference
In addition to all arguments, the following fields are also available
- `id` - the identifier of the resource

### Import
```
$ terraform import sigsci_corp_cloudwaf_instance.test_corp_cloudwaf id
```
