### Example Usage

```hcl-terraform
resource "sigsci_corp_signal_tag" "test" {
  short_name  = "example-signal-tag"
  description = "An example of a custom signal tag"
}
```

### Argument Reference
 - `short_name` - (Required) The display name of the signal tag
 - `description` -  Optional signal tag description

### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource
 - `configurable` - boolean flag for configurable
 - `informational` - boolean flag for informational 
 - `needs_response` - boolean flag indicating if the tag needs a response
 
 ### Import
You can import corp lists with the generic corp import formula
 
Example:
```shell script
terraform import sigsci_corp_signal_tag.test id
```