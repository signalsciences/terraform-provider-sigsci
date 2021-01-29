### Example Usage

```hcl-terraform
resource "sigsci_site_integration" "test_integration" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "slack"
  url             = "https://wat.slack.com"
  events          = ["listCreated"]
}
```
|Warning: You must terraform apply with the option parallelism=1 when using this resource or risk data inconsistencies! [See the FAQ.](https://github.com/signalsciences/terraform-provider-sigsci/blob/master/docs/guides/FAQ.md)|
|---|

### Argument Reference
- `site_short_name` - (Required) Identifying name of the site
- `type` - (Required) Type of integration. One of (mailingList, slack, microsoftTeams)
- `url` -  (Required) Integration Url
- `events` - Optional. Array of event types to fire on. Visit https://docs.signalsciences.net/integrations to find out which events the service you are connecting allows.

### Attributes Reference
In addition to all arguments, the following fields are also available
- `id` - the identifier of the resource
- `name` - the name of the integration

### Import
You can import site integrations with the generic site import formula

Example:
```shell script
terraform import sigsci_site_integration.test_integration site_short_name:id
```