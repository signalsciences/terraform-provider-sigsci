### Example Usage

```hcl-terraform
resource "sigsci_corp_integration" "test_corp_integration" {
  type   = "slack"
  url    = "https://signalsciences.slack.com"
  events = ["newSite", "enableSSO"]
}
```
|Warning: You must terraform apply with the option parallelism=1 when using this resource or risk data inconsistencies! [See the FAQ.](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/guides/FAQ.md)|
|---|

### Argument Reference
- `type` - (Required) Type of integration. One of (mailingList, slack, microsoftTeams)
- `url` -  (Required) Integration Url
- `events` - Optional. Array of event types to fire on. Visit https://docs.signalsciences.net/integrations to find out which events the service you are connecting allows. 

### Attributes Reference
In addition to all arguments, the following fields are also available
- `id` - the identifier of the resource
- `name` - the name of the integration

### Import
You can import corp integrations with the generic corp import formula

Example:
```shell script
terraform import sigsci_corp_integration.test_corp_integration id
```