### Example Usage

```hcl-terraform
resource "sigsci_corp_rule" "test" {
  site_short_names = [sigsci_site.my-site.short_name]
  type             = "signal"
  corp_scope       = "specificSites"
  enabled          = true
  group_operator   = "any"
  signal           = "SQLI"
  reason           = "Example corp rule"
  expiration       = ""

  conditions {
    type     = "single"
    field    = "ip"
    operator = "equals"
    value    = "1.2.3.4"
  }
  conditions {
    type     = "single"
    field    = "ip"
    operator = "equals"
    value    = "1.2.3.5"
  }
  actions {
    type = "excludeSignal"
  }
}
```

### Argument Reference
 - `site_short_names` - (Required) Sites with the rule available. Rules with a global corpScope will return '[]'.
 - `type`  Type of rule (request, signal exclusion)
 - `corp_scope` -  Whether the rule is applied to all sites or to specific sites. (global, specificSites)
 - `enabled` -   enabled or disabled
 - `group_operator` -   Conditions that must be matched when evaluating the request (all, any)
 - `signal`  -   The signal id of the signal being excluded
 - `reason`  -   Description of the rule
 - `expiration` -   Date the rule will automatically be disabled. If rule is always enabled, will return empty string (RFC3339 date time)
 - `conditions` -   Conditions on which the rule should trigger. May be recursively nest up to 3 times.
   - `type` - (Required) (group, single)
   - `groupOperator` -  type: group - Conditions that must be matched when evaluating the request (all, any)
   - `field` -  type: single - (scheme, method, path, useragent, domain, ip, responseCode, agentname, paramname, paramvalue, country, name, valueString, valueIp, signalType)
   - `operator` -  type: single - (equals, doesNotEqual, contains, doesNotContain, like, notLike, exists, doesNotExist, inList, notInList)
   - `value` -  type: single - See request fields (https://docs.signalsciences.net/using-signal-sciences/features/rules/#request-fields)
   - `conditions` -  Conditions on which this condition should trigger. Can recursively add this 3 deep.
 - `actions` - Action to take when triggered
   - `type` - (block, allow, excludeSignal) 

### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource
 
 ### Import
 You can import corp rules with the generic corp import formula
 
Example: 
```shell script
terraform import sigsci_corp_rule.test id 
```