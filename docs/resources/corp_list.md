### Example Usage

```hcl-terraform
resource "sigsci_corp_list" "test" {
  name        = "My corp list"
  type        = "ip"
  description = "Some IPs"
  entries = [
    "4.5.6.7",
    "2.3.4.5",
    "1.2.3.4",
  ]
}
```

### Argument Reference
 - `name` - (Required) Descriptive List name
 - `type` - (Required) List types (string, ip, country, wildcard)
 - `description` - (Optional) List description
 - `entries` - (Required) List entries
 
 ### Import
 You can import corp lists with the generic corp import formula
 
Example: 
```shell script
terraform import sigsci_corp_list.test id
```