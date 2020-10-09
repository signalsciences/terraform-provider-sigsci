### Example Usage

```hcl-terraform
resource "sigsci_site_allowlist" "test" {
  site_short_name = sigsci_site.my-site.short_name
  source          = "1.2.2.1"
  note            = "sample allowlist"
}
```

### Argument Reference
- `site_short_name` - (Required) Identifying name of the site
- `source` - (Required) IP address
- `note` -  (Required) Note associated with the tag
- `expires` -  Optional RFC3339-formatted datetime in the future. Omit this parameter if it does not expire.

### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource

### Import
You can import corp lists with the generic site import formula

Example:
```shell script
terraform import sigsci_site_allowlist.test site_short_name:id
```