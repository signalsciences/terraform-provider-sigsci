---
page_title: "sigsci Provider"
subcategory: ""
description: |-
  
---

## Requirements
* [Terraform](https://www.terraform.io/downloads.html) 0.12.x, 0.13.x
* [Go](https://golang.org/doc/install) 1.18

## Building the provider
Build with make and the resulting binary will be terraform-provider-sigsci.

First make the correct directory, cd to it, and checkout the repo.  make build will then build the provider and output it to terraform-provider-sigsci
```shell script
git clone git@github.com:signalsciences/terraform-provider-sigsci.git
cd terraform-provider-sigsci && make build
```

# sigsci Provider



## Example Usage

```terraform
provider "sigsci" {
  corp  = "" // Required. may also provide via env variable SIGSCI_CORP
  email = "" // Required. may also provide via env variable SIGSCI_EMAIL
  //  auth_token = "" // May also provide via env variable SIGSCI_TOKEN
  //  password = ""   // May also provide via env variable SIGSCI_PASSWORD
  //  fastly_key = "" // May also provide via env variable FASTLY_KEY. Required for Edge Deployments functionality.
}
```

## Schema

### Required
- `corp` (String) Corp short name (id)
- `email` (String) The email to be used for authentication

### Optional

- `api_url` (String) URL override for testing
- `auth_token` (String, Sensitive) The token used for authentication specify either the password or the token
- `fastly_key` (String, Sensitive) The Fastly API key used for deploying Signal Sciences as a Fastly edge security service. For edge deployment service       calls, the Fastly key must have write access to the given service.
- `password` (String, Sensitive) The password used to for authentication specify either the password or the token
- `validate` (Bool) Enable validation of API credentials during provider initialization. Default is true

# Resources

For examples of how to use each resource, see [docs/resources](./resources).
