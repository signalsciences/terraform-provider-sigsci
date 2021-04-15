### Example Usage

```hcl-terraform
resource "sigsci_site_rule" "test" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "signal"
  group_operator  = "any"
  enabled         = true
  reason          = "Example site rule update"
  signal          = "SQLI"
  expiration      = ""

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
    conditions {
      type           = "multival"
      field          = "ip"
      operator       = "equals"
      group_operator = "all"
      value          = "1.2.3.8"
    }
  }

  actions {
    type = "excludeSignal"
    signal = "corp.signal_id" 
  }
}
```

### Argument Reference
 - `site_short_names` - (Required) Sites with the rule available. Rules with a global corpScope will return '[]'.
 - `type`  Type of rule (request, signal, multival)
 - `enabled` - (Required)  enabled or disabled
 - `group_operator` -   Conditions that must be matched when evaluating the request (all, any)
 - `signal`  -   The signal id of the signal being excluded or tagged. Only used for type=signal
 - `reason`  -   Description of the rule
 - `expiration` -  (Required) Date the rule will automatically be disabled. If rule is always enabled, will return empty string (RFC3339 date time)
 - `conditions` -   Conditions on which the rule should trigger. May be recursively nest up to 3 times.
   - `type` - (Required) (group, single)
   - `group_operator` -  type: group - Conditions that must be matched when evaluating the request (all, any)
   - `field` -  type: single - (scheme, method, path, useragent, domain, ip, responseCode, agentname, paramname, paramvalue, country, name, valueString, valueIp, signalType, queryParameter)
   - `operator` -  type: single - (equals, doesNotEqual, contains, doesNotContain, like, notLike, exists, doesNotExist, inList, notInList)
   - `value` -  type: single - See request fields (https://docs.signalsciences.net/using-signal-sciences/features/rules/#request-fields)
   - `conditions` -  Conditions on which this condition should trigger. Can recursively add this 3 deep.
 - `actions` - Action to take when triggered
   - `type` - (block, allow, excludeSignal) 
   - `signal` - id of signal to be tagged with or excluded

### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource
 
 ### Import
 You can import site rules with the generic site import formula
 
Example: 
```shell script
terraform import sigsci_site_rule.test site_short_name:id 
```
