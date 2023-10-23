package provider

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func suppressEquivalentTrimSpaceDiffs(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func suppressRequestLoggingDefaultDiffs(k, old, new string, d *schema.ResourceData) bool {
	if new == "" && old == "sampled" {
		return true
	}
	return false
}

type providerMetadata struct {
	Corp   string
	Client sigsci.Client
}

func flattenStringArray(entries []string) []interface{} {
	interfaceArray := make([]interface{}, len(entries))
	for i, val := range entries {
		interfaceArray[i] = val
	}
	return interfaceArray
}

func expandStringArray(entries *schema.Set) []string {
	listOfEntries := entries.List()
	strArray := make([]string, len(listOfEntries))
	for i, e := range listOfEntries {
		strArray[i] = e.(string)
	}
	return strArray
}

func flattenDetections(detections []sigsci.Detection) []interface{} {
	detectionsSet := make([]interface{}, len(detections))
	for i, detection := range detections {
		fieldSet := make([]interface{}, len(detection.Fields))
		for j, field := range detection.Fields {
			if _, ok := field.Value.(float64); ok {
				field.Value = fmt.Sprintf("%.0f", field.Value.(float64))
			}
			fieldMap := map[string]interface{}{
				"name":  field.Name,
				"value": field.Value,
			}
			fieldSet[j] = fieldMap
		}
		detectionMap := map[string]interface{}{
			"id":      detection.ID,
			"name":    detection.Name,
			"enabled": detection.Enabled,
			"fields":  fieldSet,
		}
		detectionsSet[i] = detectionMap
	}
	return detectionsSet
}

func expandDetections(entries *schema.Set) []sigsci.Detection {
	var detections []sigsci.Detection
	for _, e := range entries.List() {
		cast := e.(map[string]interface{})
		fieldsI := cast["fields"].(*schema.Set)
		var fields []sigsci.ConfiguredDetectionField
		for _, v := range fieldsI.List() {
			castV := v.(map[string]interface{})
			if _, ok := castV["value"].(float64); ok {
				castV["value"] = int(castV["value"].(float64))
			}

			//switch castV["value"].(type) {
			//case float64:
			//	castV["value"] = int(castV["value"].(float64))
			//case float32:
			//	castV["value"] = int(castV["value"].(float32))
			//}
			fields = append(fields, sigsci.ConfiguredDetectionField{
				Name:  castV["name"].(string),
				Value: castV["value"],
			})
		}

		detections = append(detections, sigsci.Detection{
			DetectionUpdateBody: sigsci.DetectionUpdateBody{
				ID:      cast["id"].(string),
				Name:    cast["name"].(string),
				Enabled: cast["enabled"].(bool),
				Fields:  fields,
			},
		})
	}
	return detections
}

func expandAlerts(entries *schema.Set) []sigsci.Alert {
	var alerts []sigsci.Alert

	for _, e := range entries.List() {
		cast := e.(map[string]interface{})

		var typ string
		if s, ok := cast["type"].(string); ok {
			typ = s
		}

		var tagName string
		if s, ok := cast["tag_name"].(string); ok {
			tagName = s
		}

		var blockDurationSeconds int
		if b, ok := cast["block_duration_seconds"].(int); ok {
			blockDurationSeconds = b
		}

		alerts = append(alerts, sigsci.Alert{
			AlertUpdateBody: sigsci.AlertUpdateBody{
				LongName:             cast["long_name"].(string),
				Interval:             cast["interval"].(int),
				Threshold:            cast["threshold"].(int),
				SkipNotifications:    cast["skip_notifications"].(bool),
				Enabled:              cast["enabled"].(bool),
				Action:               cast["action"].(string),
				BlockDurationSeconds: blockDurationSeconds,
			},
			Type:    typ,
			TagName: tagName,
		})
	}
	return alerts
}

func flattenAlerts(alerts []sigsci.Alert) []interface{} {
	alertsSet := make([]interface{}, len(alerts))
	for i, alert := range alerts {
		alertsSet[i] = map[string]interface{}{
			"id":                     alert.ID,
			"long_name":              alert.LongName,
			"interval":               alert.Interval,
			"threshold":              alert.Threshold,
			"skip_notifications":     alert.SkipNotifications,
			"enabled":                alert.Enabled,
			"action":                 alert.Action,
			"block_duration_seconds": alert.BlockDurationSeconds,
		}
	}
	return alertsSet
}

func getListAdditionsDeletions(existlist, newlist []string) (additions []string, deletions []string) {
	if len(existlist) == 0 {
		return newlist, []string{}
	}
	setExist := make(map[string]string, len(existlist))
	for _, exE := range existlist {
		setExist[exE] = ""
	}
	add := []string{}
	for _, nwE := range newlist {
		if _, ok := setExist[nwE]; !ok {
			add = append(add, nwE)
		}
	}
	setNew := make(map[string]string, len(newlist))
	for _, nwE := range newlist {
		setNew[nwE] = ""
	}
	del := []string{}
	for _, exE := range existlist {
		if _, ok := setNew[exE]; !ok {
			del = append(del, exE)
		}
	}

	return add, del
}

func diffTemplateAlerts(orig, new []sigsci.Alert) ([]sigsci.Alert, []sigsci.Alert, []sigsci.Alert) {
	return calcAlertAdds(orig, new), calcAlertUpdates(orig, new), calcAlertDeletes(orig, new)
}

func calcAlertAdds(old, new []sigsci.Alert) []sigsci.Alert {
	var alertAdds []sigsci.Alert
	for _, newA := range new {
		exists := false
		for _, oldA := range old {
			if oldA.ID == newA.ID {
				exists = true
			}
		}

		if !exists {
			alertAdds = append(alertAdds, newA)
		}
	}
	return alertAdds
}

func alertEquals(a, b sigsci.Alert) bool {
	return a.Enabled == b.Enabled &&
		a.Action == b.Action &&
		a.SkipNotifications == b.SkipNotifications &&
		a.Interval == b.Interval &&
		a.Threshold == b.Threshold &&
		a.LongName == b.LongName &&
		a.BlockDurationSeconds == b.BlockDurationSeconds
}

func calcAlertUpdates(old, new []sigsci.Alert) []sigsci.Alert {
	var alertUpdates []sigsci.Alert
	for _, oldA := range old {
		for _, newA := range new {
			if oldA.ID == newA.ID {
				if !alertEquals(oldA, newA) {
					alertUpdates = append(alertUpdates, newA)
				}
			}
		}
	}
	return alertUpdates
}

func calcAlertDeletes(old, new []sigsci.Alert) []sigsci.Alert {
	var alertDels []sigsci.Alert
	for _, oldA := range old {
		exists := false
		for _, newA := range new {
			if newA.ID == oldA.ID {
				exists = true
			}
		}

		if !exists {
			alertDels = append(alertDels, oldA)
		}
	}
	return alertDels
}

func diffTemplateDetections(template string, orig, new []sigsci.Detection) ([]sigsci.Detection, []sigsci.Detection, []sigsci.Detection) {
	return calcDetectionAdds(template, orig, new), calcDetectionUpdates(template, orig, new), calcDetectionDeletes(orig, new)
}

func calcDetectionAdds(templateID string, old, new []sigsci.Detection) []sigsci.Detection {
	var detectionAdds []sigsci.Detection
	for _, newD := range new {
		exists := false
		for _, oldD := range old {
			if newD.Name == oldD.Name {
				exists = true
			}
		}
		if !exists {
			newD.Name = templateID
			// Convert strings to numbers
			if newD.Fields == nil {
				newD.Fields = []sigsci.ConfiguredDetectionField{}
			}
			for i, f := range newD.Fields {
				if v, err := strconv.Atoi(f.Value.(string)); err == nil {
					newD.Fields[i].Value = v
				}
			}
			detectionAdds = append(detectionAdds, newD)
		}
	}
	return detectionAdds
}

func calcDetectionDeletes(old, new []sigsci.Detection) []sigsci.Detection {
	var detectionDeletes []sigsci.Detection
	for _, oldD := range old {
		exists := false
		for _, newD := range new {
			if oldD.Name == newD.Name {
				exists = true
			}
		}
		if !exists {
			detectionDeletes = append(detectionDeletes, oldD)
		}
	}
	return detectionDeletes
}

func calcDetectionUpdates(templateID string, old, new []sigsci.Detection) []sigsci.Detection {
	var detectionUpdates []sigsci.Detection
	for _, oldD := range old {
		for _, newD := range new {
			if oldD.Name == newD.Name {
				if oldD.Enabled != newD.Enabled || !detectionFieldsEqual(newD.Fields, oldD.Fields) {
					newD.Name = templateID
					detectionUpdates = append(detectionUpdates, newD)
				}
			}
		}
	}
	return detectionUpdates
}

func detectionFieldsEqual(old, new []sigsci.ConfiguredDetectionField) bool {
	if len(old) != len(new) {
		return false
	}
	sort.Slice(old, func(i, j int) bool {
		return old[i].Name < old[j].Name
	})
	sort.Slice(new, func(i, j int) bool {
		return new[i].Name < new[j].Name
	})
	for i := range old {
		if old[i].Name != new[i].Name {
			return false
		}
		if old[i].Value != new[i].Value {
			return false
		}
	}
	return true
}

func existsInInt(needle int, haystack ...int) bool {
	for _, i := range haystack {
		if i == needle {
			return true
		}
	}
	return false
}

func existsInString(needle string, haystack ...string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func existsInRange(needle int, min, max int) bool {
	if needle >= min && needle <= max {
		return true
	}
	return false
}

func validStringLength(needle string, min, max int) bool {
	length := len(needle)

	return length >= min && length <= max
}

func expandRuleConditions(conditionsResource *schema.Set) []sigsci.Condition {
	var conditions []sigsci.Condition
	for _, genericElement := range conditionsResource.List() {
		castElement := genericElement.(map[string]interface{})
		if _, ok := castElement["conditions"]; !ok {
			c := sigsci.Condition{
				Type:          castElement["type"].(string),
				Field:         castElement["field"].(string),
				Operator:      castElement["operator"].(string),
				GroupOperator: castElement["group_operator"].(string),
				Value:         castElement["value"].(string),
			}
			conditions = append(conditions, c)
			continue
		}
		c := sigsci.Condition{
			Type:          castElement["type"].(string),
			Field:         castElement["field"].(string),
			Operator:      castElement["operator"].(string),
			Value:         castElement["value"].(string),
			GroupOperator: castElement["group_operator"].(string),
			Conditions:    expandRuleConditions(castElement["conditions"].(*schema.Set)),
		}
		conditions = append(conditions, c)
	}
	return conditions
}

func flattenRuleConditions(conditions []sigsci.Condition) []interface{} {
	conditionsMap := make([]interface{}, len(conditions))
	for i, condition := range conditions {
		conditionMap := map[string]interface{}{
			"type":           condition.Type,
			"field":          condition.Field,
			"operator":       condition.Operator,
			"value":          condition.Value,
			"group_operator": condition.GroupOperator,
		}
		if len(condition.Conditions) != 0 {
			conditionMap["conditions"] = flattenRuleConditions(condition.Conditions)
		}
		conditionsMap[i] = conditionMap
	}
	return conditionsMap
}

func expandRuleActions(actionsResource *schema.Set) []sigsci.Action {
	var actions []sigsci.Action
	for _, genericElement := range actionsResource.List() {
		castElement := genericElement.(map[string]interface{})
		var signal string

		if castElement["signal"] != nil {
			signal = castElement["signal"].(string)
		}
		var responseCode int
		if castElement["response_code"] != nil {
			responseCode = castElement["response_code"].(int)
		}

		var redirectURL string
		if castElement["redirect_url"] != nil {
			redirectURL = castElement["redirect_url"].(string)
		}

		a := sigsci.Action{
			Type:         castElement["type"].(string),
			Signal:       signal,
			ResponseCode: responseCode,
			RedirectURL:  redirectURL,
		}
		actions = append(actions, a)
	}
	return actions
}

func expandRuleRateLimit(rateLimitResource map[string]interface{}) *sigsci.RateLimit {
	var threshold, interval, duration int
	var err error
	if val, ok := rateLimitResource["threshold"]; ok {
		threshold, err = strconv.Atoi(val.(string))
		if err != nil {
			return nil
		}
	} else {
		return nil
	}
	if val, ok := rateLimitResource["interval"]; ok {
		interval, err = strconv.Atoi(val.(string))
		if err != nil {
			return nil
		}
	}
	if val, ok := rateLimitResource["duration"]; ok {
		duration, err = strconv.Atoi(val.(string))
		if err != nil {
			return nil
		}
	}

	return &sigsci.RateLimit{
		Threshold: threshold,
		Interval:  interval,
		Duration:  duration,
	}
}

func flattenRuleRateLimit(rateLimit *sigsci.RateLimit) map[string]string {
	if rateLimit == nil {
		return nil
	}
	return map[string]string{
		"threshold": strconv.Itoa(rateLimit.Threshold),
		"interval":  strconv.Itoa(rateLimit.Interval),
		"duration":  strconv.Itoa(rateLimit.Duration),
	}
}

func flattenRuleActions(actions []sigsci.Action, customResponseCode bool) []interface{} {
	actionsMap := make([]interface{}, len(actions))
	for i, action := range actions {

		actionMap := map[string]interface{}{
			"type":   action.Type,
			"signal": action.Signal,
		}
		// customResponseCode is enabled for site rules but disabled for corp rules
		// this boolean flag reflects the differences and flattens objects accordingly
		if customResponseCode {
			// response code is set to 0 by sigsci api when action.type != "block"
			// for types such as "allow" or "logRequest", response code is irrelevant and hence not provided in API response
			actionMap["response_code"] = action.ResponseCode

			if action.ResponseCode == 301 || action.ResponseCode == 302 {
				actionMap["redirect_url"] = action.RedirectURL
			}
		}
		actionsMap[i] = actionMap
	}

	return actionsMap
}

func resourceSiteImport(siteID string) (site string, id string, err error) {
	parts := strings.SplitN(siteID, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected site:id", siteID)
	}

	return parts[0], parts[1], nil
}

var siteImporter = schema.ResourceImporter{
	State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		site, id, err := resourceSiteImport(d.Id())
		if err != nil {
			return nil, err
		}
		d.Set("site_short_name", site)
		d.SetId(id)
		return []*schema.ResourceData{d}, nil
	},
}

var KnownSingleConditionFields = []string{
	"scheme", "method", "path", "useragent", "domain", "ip", "responseCode", "agentname",
	"paramname", "paramvalue", "country", "name", "valueString", "valueIp", "signalType",
	"value",
}

var KnownMultivalConditionFields = []string{
	"signal", "requestHeader", "queryParameter", "postParameter", "requestCookie",
	"responseHeader",
}

func validateConditionField(val interface{}, key string) ([]string, []error) {
	if existsInString(val.(string), KnownSingleConditionFields...) || existsInString(val.(string), KnownMultivalConditionFields...) {
		return nil, nil
	}

	return []string{fmt.Sprintf("received %q for conditions.field. This is not necessarily an error, but we only know about the following values. If this is a new value, please open a PR to get it added.\n(%s)", val.(string), strings.Join(append(KnownSingleConditionFields, KnownMultivalConditionFields...), ", "))}, nil
}

func validateActionResponseCode(val interface{}, key string) ([]string, []error) {
	// response code needs to be 301, 302, or 400-599
	code := val.(int)
	if (code >= 301 && code <= 302) || (code >= 400 && code <= 599) {
		return nil, nil
	}
	rangeError := fmt.Errorf("received action responseCode '%d'. should be in 301, 302, or 400-599", code)
	return nil, []error{rangeError}
}

func validateActionRedirectURL(val interface{}, key string) ([]string, []error) {
	url := val.(string)
	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "/") {
		return nil, nil
	}
	rangeError := fmt.Errorf("received redirect url '%s'. must start with either 'http' or '/'", url)
	return nil, []error{rangeError}
}

func validateRegion(val interface{}, key string) ([]string, []error) {
	// https://docs.fastly.com/signalsciences/api/#_corps__corpName__cloudwafInstances_post
	regionList := []string{
		"us-east-1",
		"us-west-1",
		"af-south-1",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-north-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"sa-east-1",
		"us-east-2",
		"us-west-2",
	}

	if existsInString(val.(string), regionList...) {
		return nil, nil
	}

	return nil, []error{fmt.Errorf("received region name %q is invalid. should be in (%s)", val.(string), strings.Join(regionList, ", "))}
}
