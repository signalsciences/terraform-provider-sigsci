# Sigsci Terraform Provider

This terraform provider is currently in beta

## Requirements
* [Terraform](https://www.terraform.io/downloads.html) 0.12.x
* [Go](https://golang.org/doc/install) 1.14

## Building the provider
Build with make and the resulting binary will be terraform-provider-sigsci.

First make the correct directory, cd to it, and checkout the repo.  make build will then build the provider and output it to terraform-provider-sigsci
```shell script
mkdir -p $GOPATH/src/github.com/signalsciences/terraform-provider-sigsci
cd $GOPATH/src/github.com/signalsciences/terraform-provider-sigsci
git clone git@github.com:signalsciences/terraform-provider-sigsci.git
make build
```

## Using the provider
You must provide corp, email, and either form of authentication.  This can be added in the provider block or with environment variables (recommended).

```hcl-terraform
provider "sigsci" {
  //  corp = ""       // Required. may also provide via env variable SIGSCI_CORP
  //  email = ""      // Required. may also provide via env variable SIGSCI_EMAIL
  //  auth_token = "" //may also provide via env variable SIGSCI_TOKEN
  //  password = ""   //may also provide via env variable SIGSCI_PASSWORD
}
```
## Corp level resources
##### Site
```hcl-terraform
resource "sigsci_site" "my-site" {
  short_name             = "manual_test"
  display_name           = "manual terraform test"
  block_duration_seconds = 1000
  block_http_code        = 303
  agent_anon_mode        = ""
}
```
##### List
```hcl-terraform
resource "sigsci_corp_list" "test_list" {
  name        = "My corp list"
  type        = "ip"
  description = "Some IPs"
  entries = [
    "4.5.6.7",
    "2.3.4.5",
    "1.2.3.4",
  ]
}
```

##### Rules
```hcl-terraform
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
    type = "excludeSignal"
  }
}
```

##### Tags
```hcl-terraform
resource "sigsci_corp_signal_tag" "test" {
  short_name  = "example-signal-tag"
  description = "An example of a custom signal tag"
}
```


## Site level resources
##### Lists
```hcl-terraform
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
```

##### Rules
```hcl-terraform
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
```

##### Tags
```hcl-terraform
resource "sigsci_site_signal_tag" "test_tag" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new signal tag"
  description     = "description"
}
```
##### Redactions
Warning: if using redactions, you **must** terraform apply with the option parallelism=1
```hcl-terraform
resource "sigsci_site_redaction" "test_redaction" {
  site_short_name = sigsci_site.my-site.short_name
  field           = "redacted_field"
  redaction_type  = 1
}
```

##### Alerts
```hcl-terraform
resource "sigsci_site_alert" "test_site_alert" {
  site_short_name = sigsci_site.my-site.short_name
  tag_name        = sigsci_site_signal_tag.test_tag.id
  long_name       = "test_alert"
  interval        = 10
  threshold       = 12
  enabled         = true
  action          = "info"
}
```

##### Templated Rules
```hcl-terraform
resource "sigsci_site_templated_rule" "test_template_rule" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "LOGINATTEMPT"
  detections {
    enabled = "true"
    fields {
      name  = "path"
      value = "/auth/*"
    }
  }

  alerts {
    long_name          = "alert 1"
    interval           = 60
    threshold          = 10
    skip_notifications = true
    enabled            = true
    action             = "info"
  }

  alerts {
    long_name          = "alert 2"
    interval           = 60
    threshold          = 1
    skip_notifications = false
    enabled            = false
    action             = "info"
  }
}
```

##### BlackList
```hcl-terraform
resource "sigsci_site_blacklist" "test" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.3.4"
  note            = "sample blacklist"
}
```

##### Whitelist
```hcl-terraform
resource "sigsci_site_whitelist" "test" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.3.4"
  note            = "sample whitelist"
}
```

##### Header Links
Warning: if using header links, you **must** terraform apply with the option parallelism=1 
```hcl-terraform
resource "sigsci_site_header_link" "test_header_link" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "test_header_link"
  type            = "request"
  link_name       = "signal sciences"
  link            = "https://www.signalsciences.net"
}
```

More information on each resource and field can be found on the [Signal Sciences Api Docs](https://docs.signalsciences.net/api/).

## Importing

Importing will vary depending on if you are importing a corp level resource or a site level resource
##### Corp Resources
```hcl-terraform
terraform import resource.name id // General form
terraform import sigsci_site.my-site test_site // Example
```

##### Site Resources
```hcl-terraform
terraform import resource.name site_short_name:id //General form
terraform import sigsci_site_list.manual-list test_site:site.manual-list //Example
```
