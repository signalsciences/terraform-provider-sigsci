### Example Usage

```hcl-terraform
resource "sigsci_site_templated_rule" "test_template_rule" {
  site_short_name = sigsci_site.my-site.short_name
  name            = "LOGINATTEMPT"
  detections {
    enabled = "true"
    fields {
      name  = "path"
      value = "/login/*"
    }
  }

  alerts {
    long_name          = "alert 1"
    interval           = 60
    threshold          = 10
    skip_notifications = true
    enabled            = true
    action             = "info"
    block_duration_seconds = sigsci_site.my-site.block_duration_seconds
  }

  alerts {
    long_name          = "alert 2"
    interval           = 60
    threshold          = 1
    skip_notifications = false
    enabled            = false
    action             = "info"
    block_duration_seconds = 64000
  }
}
```
|Warning: You must terraform apply with the option parallelism=1 when using this resource or risk data inconsistencies! [See the FAQ.](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/guides/FAQ.md)|
|---|

### Argument Reference
- `site_short_name` - (Required) Identifying name of the site
- `name` - (Required) Field name. A list of available names is below.
- `detections` - (Required) Type of redaction (0: Request Parameter, 1: Request Header, 2: Response Header)
  - `enabled` - (Required) A flag to toggle this detection
  - `fields` - (Required) detection fields that should trigger an alert
    - `name` - field name
    - `value` - field value
- `alerts` -  Type of redaction (0: Request Parameter, 1: Request Header, 2: Response Header)
  - `long_name` - (Required) A human readable description of the alert. Must be between 3 and 25 characters.
  - `skip_notifications` - (Required) A flag to disable external notifications - slack, webhooks, emails, etc.
  - `enabled` - (Required) A flag to toggle this alert.
  - `action` - (Required) A flag that describes what happens when the alert is triggered. 'info' creates an incident in the dashboard. 'flagged' creates an incident and blocks traffic for 24 hours.
  - `block_duration_seconds` - (Required) Duration to block for in seconds
  - `interval` - The number of minutes of past traffic to examine. Must be 0, 1, 10 or 60.
  - `threshold` - The number of occurrences of the tag in the interval needed to trigger the alert.

### Attributes Reference
In addition to all arguments, the following fields are also available
 - `id` - the identifier of the resource
 - `name` - Name of templated rule

### Available Templates
```javascript
//These are the valid values for the 'name' field.  
signals = ["AWS-SSRF", "CVE-2017-5638", "CVE-2017-7269", "CVE-2017-9805", "CVE-2018-9206", "CVE-2018-11776", 
  "CVE-2019-0192", "CVE-2019-0193", "CVE-2019-0232", "CVE-2019-2725", "CVE-2019-3396", "CVE-2019-5418", "CVE-2019-6340",
  "CVE-2019-14234", "CVE-2019-16759", "CVE-2021-26855", "LOGINATTEMPT", "LOGINFAILURE", "LOGINSUCCESS", "REGATTEMPT", 
  "REGFAILURE", "REGSUCCESS"]
```

If you do not see the signal you want in this list, check out the [Rules page](https://github.com/signalsciences/terraform-provider-sigsci/blob/main/docs/resources/site_rule.md) for some additional templates

### Import
You can import corp lists with the generic site import formula

Example:
```shell script
terraform import sigsci_site_templated_rule.test site_short_name:id
```
