# Sigsci Terraform Provider


## Requirements
* [Terraform](https://www.terraform.io/downloads.html) > 0.12.x
* [Go](https://golang.org/doc/install) 1.20

Check out the [Terraform Documentation](https://www.terraform.io/docs/configuration/index.html) and their [Introduction](https://www.terraform.io/intro/index.html) for more information on terraform

## Building the provider
If you are using terraform >0.13.x, our release can be automatically downloaded from their registry using the block described in "Using the provider"

If you are using terraform 0.12.x, you must either build or copy the appropriate executable to your plugin directory. ex `terraform.d/plugins/darwin_amd64`

You may find prebuilt binaries in our [Releases](https://github.com/signalsciences/terraform-provider-sigsci/releases).

If you wish to build from source, first make the correct directory, cd to it, and checkout the repo.  Running `make build` will then build the provider and output it to terraform-provider-sigsci
```shell script
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
  //  corp           = "" // Required. may also provide via env variable SIGSCI_CORP
  //  email          = "" // Required. may also provide via env variable SIGSCI_EMAIL
  //  auth_token     = "" // May also provide via env variable SIGSCI_TOKEN
  //  password       = "" // May also provide via env variable SIGSCI_PASSWORD
  //  fastly_api_key = "" // May also provide via env variable FASTLY_API_KEY. Required for Edge Deployments functionality.
}
```

## Resources

Resource documentation and examples can be found in [docs/resources](./docs/resources).

## FAQ

FAQ can be found in [docs/guides/FAQ.md](./docs/guides/FAQ.md).

## Example
[main.tf](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/main.tf) has an example of every resource.
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

## Errors

Errors occasionally occur when updating certain resources. If an error occurs please try re-running with `-parallelism=1`:

```
$ terraform apply -parallelism=1
```

If running with `-parallelism=1` does not resolve the error, please open an issue.
