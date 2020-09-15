### Example Usage

```hcl-terraform
resource "sigsci_site" "my-site" {
  short_name             = "manual_test"
  display_name           = "manual terraform test"
  block_duration_seconds = 86400
  block_http_code        = 406
  agent_anon_mode        = ""
  agent_level            = "block"
}
```

### Argument Reference
 - `short_name` - (Required) Identifying name of the site
 - `display_name` - (Required) Display name of the site
 - `block_duration_seconds` -  Duration to block an IP in seconds
 - `block_http_code` - HTTP response code to send when when traffic is being blocked
 - `agent_anon_mode` - Agent IP anonimization mode - 'EU' or '' (off)
 - `agent_level` -  Agent action level - 'block', 'log' or 'off'
 
 ### Import
You can import corp lists with the generic corp import formula
 
Example:
```shell script
terraform import sigsci_site.test id
```