### Example Usage

```hcl-terraform
resource "sigsci_site_header_link" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "test_header_link"
  type            = "request"
  link_name       = "signal sciences"
  link            = "https://www.signalsciences.net"
}
```
|Warning: You must terraform apply with the option parallelism=1 when using this resource or risk data inconsistencies! [See the FAQ.](https://github.com/signalsciences/terraform-provider-sigsci/blob/master/docs/guides/FAQ.md)|
|---|

### Argument Reference
 - `site_short_name` - (Required) Identifying name of the site
 - `name` - (Required) Name of header
 - `type` - (Required) The type of header, either 'request' or 'response'
 - `link_name` - Name of header link for display purposes
 - `link` - (Required) External link
 
 ### Import
You can import corp lists with the generic site import formula

### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource
 
Example:
```shell script
terraform import sigsci_site_header_link.test site_short_name:id
```