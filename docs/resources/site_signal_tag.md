### Example Usage

```hcl-terraform
resource "sigsci_site_signal_tag" "test" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "My new signal tag"
  description     = "description of tag"
}
```

### Argument Reference
- `site_short_name` - (Required) Identifying name of the site
- `name` - (Required) Field name
- `description` -  Type of redaction (0: Request Parameter, 1: Request Header, 2: Response Header)

### Import
You can import corp lists with the generic site import formula

Example:
```shell script
terraform import sigsci_site_signal_tag.test site_short_name:id
```