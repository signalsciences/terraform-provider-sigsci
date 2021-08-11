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
 - `type`  Type of rule (request, signal, multival, templatedSignal).
 - `enabled` - (Required) enabled or disabled
 - `group_operator` -   Conditions that must be matched when evaluating the request (all, any)
 - `signal`  -   The signal id of the signal being excluded or tagged. Only used for type=signal
 - `reason`  -   Description of the rule
 - `expiration` - (Required) Date the rule will automatically be disabled. If rule is always enabled, will return empty string (RFC3339 date time)
 - `conditions` -   Conditions on which the rule should trigger. May be recursively nest up to 3 times.
   - `type` - (Required) (group, single)
   - `group_operator` -  type: group - Conditions that must be matched when evaluating the request (all, any)
   - `field` -  type: single - (scheme, method, path, useragent, domain, ip, responseCode, agentname, paramname, paramvalue, country, name, valueString, valueIp, signalType, queryParameter)
   - `operator` -  type: single - (equals, doesNotEqual, contains, doesNotContain, like, notLike, exists, doesNotExist, inList, notInList)
   - `value` -  type: single - See request fields (https://docs.signalsciences.net/using-signal-sciences/features/rules/#request-fields)
   - `conditions` -  Conditions on which this condition should trigger. Can recursively add this 3 deep.
 - `actions` - Action to take when triggered
   - `type` - (block, allow, excludeSignal, addSignal). A RateLimit rule has valid values of (logRequest, blockSignal)
   - `signal` - id of signal to be tagged with or excluded
 - `rate_limit` - Enable rate limiting on this rule
   -  `threshold` - Number of requests to count before rate limiting is activated
   -  `interval` -  Length of time in minutes the threshold should be measured for (default: 1, options: 1, 5, 10)
   -  `duration` -  Length of time in seconds to enforce the rule for once activated (default: 600, minimum: 300, maximum: 3600)
### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource
 


### Templated Signals
We have curated a list of templates for common rules, the full list of available signals is available below. 

For these you must specify type = "templatedSignal". 
Note that they will show up in the site "Templated Rules" page in the Console.

```hcl-terraform
resource "sigsci_site_rule" "testsignal" {
  site_short_name = sigsci_site.my-site.short_name
  type            = "templatedSignal"
  group_operator  = "all"
  enabled         = true
  reason          = "Example site rule update"
  signal          = "PW-RESET-ATTEMPT"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "method"
    operator = "equals"
    value    = "POST"
  }

  conditions {
    type     = "single"
    field    = "path"
    operator = "equals"
    value    = "/change-password"
  }

  conditions {
    type           = "multival"
    group_operator = "all"
    conditions {
      field    = "name"
      operator = "equals"
      type     = "single"
      value    = "foo"
    }
  }
}
```

```javascript
// These are all of the valid values for signal
signals = ["2FA-CHANGED","2FA-DISABLED","ADDRESS-CHANGED","CC-VAL-ATTEMPT","CC-VAL-FAILURE", "CC-VAL-SUCCESS",
   "EMAIL-CHANGED","EMAIL-VALIDATION","GC-VAL-ATTEMPT","GC-VAL-FAILURE", "GC-VAL-SUCCESS","INFO-VIEWED",
   "INVITE-ATTEMPT","INVITE-FAILURE","INVITE-SUCCESS", "KBA-CHANGED","MESSAGE-SENT","PW-CHANGED","PW-RESET-ATTEMPT",
   "RSRC-ID-ENUM-ATTEMPT", "RSRC-ID-ENUM-FAILURE","RSRC-ID-ENUM-SUCCESS","RSRC-ID-ENUM-SUCCESS","USER-ID-ENUM-ATTEMPT", 
   "USER-ID-ENUM-FAILURE","USER-ID-ENUM-SUCCESS","USER-ID-ENUM-SUCCESS","WRONG-API-CLIENT"]
```
If you do not see the signal you want in this list, check out the [Templated Rules page](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_templated_rule.md) for some additional templates

### Import
You can import site rules with the generic site import formula

Example:
```shell script
terraform import sigsci_site_rule.test site_short_name:id 
```
