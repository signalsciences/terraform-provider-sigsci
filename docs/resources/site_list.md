### Example Usage

```hcl-terraform
resource "sigsci_site_list" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new list"
  type            = "ip"
  description     = "Some IPs"
  entries = [
    "4.5.6.7",
    "2.3.4.5",
    "1.2.3.4",
  ]
}
```

### Argument Reference
- `site_short_name` - (Required) Identifying name of the site
- `name` - (Required) Descriptive list name
- `type` - (Required) List types (string, ip, country, wildcard)
- `description` - Optional list description
- `entries` - (Required) List entries

### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource

### Import
You can import corp lists with the generic site import formula

Example:
```shell script
terraform import sigsci_site_list.test site_short_name:id
```