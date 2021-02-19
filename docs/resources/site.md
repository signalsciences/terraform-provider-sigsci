### Example Usage

```hcl-terraform
resource "sigsci_site" "my-site" {
  short_name             = "manual_test"
  display_name           = "manual terraform test"
  block_duration_seconds = 86400
  agent_anon_mode        = ""
  agent_level            = "block"
}
```

### Argument Reference
 - `short_name` - (Required) Identifying name of the site
 - `display_name` - (Required) Display name of the site
 - `block_duration_seconds` -  Duration to block an IP in seconds
 - `agent_anon_mode` - Agent IP anonimization mode - 'EU' or '' (off)
 - `agent_level` -  Agent action level - 'block', 'log' or 'off'
 
 ### Import
You can import corp lists with the generic corp import formula

### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource
 - `http_block_code` - HTTP response code to send when traffic is being blocked
 - `primary_agent_key` - Primary agent key containing secret and access keys 
   - `name`
   - `secret_key`
   - `access_key`
 
Example:
```shell script
terraform import sigsci_site.test id
```
