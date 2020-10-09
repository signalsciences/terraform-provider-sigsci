### Example Usage

```hcl-terraform
resource "sigsci_site_blocklist" "test" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.3.4"
  note            = "sample blocklist"
  expires         = "2012-11-01T22:08:41+00:00"
}
```

### Argument Reference
 - `site_short_name` - (Required) Identifying name of the site
 - `source` - (Required) IP address
 - `note` - (Required) Note associated with the tag
 - `expires` - Optional RFC3339-formatted datetime in the future. Omit this parameter if it does not expire.
 
 ### Import
You can import corp lists with the generic site import formula

### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource
 
Example:
```shell script
terraform import sigsci_site_blocklist.test site_short_name:id
```