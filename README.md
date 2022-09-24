# Sigsci Terraform Provider


## Requirements
* [Terraform](https://www.terraform.io/downloads.html) > 0.12.x
* [Go](https://golang.org/doc/install) 1.17

Check out the [Terraform Documentation](https://www.terraform.io/docs/configuration/index.html) and their [Introduction](https://www.terraform.io/intro/index.html) for more information on terraform

## Building the provider
If you are using terraform >0.13.x, our release can be automatically downloaded from their registry using the block described in "Using the provider"
 

If you are using terraform 0.12.x, you must either build or copy the appropriate executable to your plugin directory. ex `terraform.d/plugins/darwin_amd64`

You may find prebuilt binaries in our [Releases](https://github.com/signalsciences/terraform-provider-sigsci/releases).

If you wish to build from source, first make the correct directory, cd to it, and checkout the repo.  Running `make build` will then build the provider and output it to terraform-provider-sigsci
```shell script
mkdir -p $GOPATH/src/github.com/signalsciences/
cd $GOPATH/src/github.com/signalsciences/
git clone git@github.com:signalsciences/terraform-provider-sigsci.git
cd terraform-provider-sigsci
make build
cp terraform-provider-sigsci ~/.terraform.d/plugins
```

## Using the provider
You must provide corp, email, and either form of authentication.  This can be added in the provider block or with environment variables (recommended).

```hcl-terraform
# Terraform 0.13.x
terraform {
  required_providers {
    sigsci = {
      source = "signalsciences/sigsci"
    }
  }
}

# Required configuration block (for all versions of terraform)
provider "sigsci" {
  //  corp = ""       // Required. may also provide via env variable SIGSCI_CORP
  //  email = ""      // Required. may also provide via env variable SIGSCI_EMAIL
  //  auth_token = "" //may also provide via env variable SIGSCI_TOKEN
  //  password = ""   //may also provide via env variable SIGSCI_PASSWORD
}
```
## Corp resources
[Site](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site.md)

[Lists](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/corp_list.md)

[Tags](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/corp_signal_tag.md)

[Rules](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/corp_rule.md)

[Integrations](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/corp_integration.md)

[Cloud WAF Instance](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/corp_cloudwaf_instance.md)

[Cloud WAF Certificate](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/corp_cloudwaf_certificate.md)

## Site resources

[Lists](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_list.md)

[Rules](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_rule.md)

[Tags](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_signal_tag.md)

[Redactions](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_redaction.md)

[Alerts](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_alert.md)

[Templated Rules](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_templated_rule.md)

[Allowlist](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_allowlist.md)

[Blocklist](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_blocklist.md)

[Header Links](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_header_link.md)

[Integrations](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_integration.md)

More information on each resource and field can be found on the [Signal Sciences Api Docs](https://docs.signalsciences.net/api/).


## Importing

All resources are importable. Importing will vary depending on if you are importing a corp level resource, or a site level resource
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


## Example
[main.tf](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/main.tf) has an example of every resource 
```hcl-terraform
resource "sigsci_site" "my-site" {
  short_name             = "manual_test"
  display_name           = "manual terraform test"
  block_duration_seconds = 86400
  block_http_code        = 406
  agent_anon_mode        = ""
  agent_level            = "block"
}

resource "sigsci_site_signal_tag" "test_tag" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new signal tag"
  description     = "description"
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
