### Example Usage

```hcl-terraform
resource "sigsci_site_redaction" "test" {
  site_short_name = sigsci_site.my-site.short_name
  field           = "redacted"
  redaction_type  = 0
}
```
|Warning: You must terraform apply with the option parallelism=1 when using this resource or risk data inconsistencies! [See the FAQ.](https://github.com/signalsciences/terraform-provider-sigsci/blob/master/docs/guides/FAQ.md)|
|---|

### Argument Reference
- `site_short_name` - (Required) Identifying name of the site
- `field` - (Required) Field name
- `redaction_type` - (Required) Type of redaction (0: Request Parameter, 1: Request Header, 2: Response Header)

### Import
You can import corp lists with the generic site import formula

Example:
```shell script
terraform import sigsci_site_redaction.test site_short_name:id
```