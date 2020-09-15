### Example Usage

```hcl-terraform
resource "sigsci_site_header_link" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "test_header_link"
  type            = "request"
  link_name       = "signal sciences 89"
  link            = "https://www.signalsciences.net"
}
```
Warning: You must terraform apply with the option parallelism=1 when using this resource or risk data inconsistencies!

### Argument Reference
 - `site_short_name` - (Required) Identifying name of the site
 - `name` - (Required) Name of header
 - `type` - (Required) The type of header, either 'request' or 'response'
 - `link_name` - Name of header link for display purposes
 - `link` - (Required) External link
 
 ### Import
You can import corp lists with the generic site import formula
 
Example:
```shell script
terraform import sigsci_site_header_link.test site_short_name:id
```