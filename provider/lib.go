package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"sort"
	"strconv"
	"strings"
)

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
	var detectionsSet = make([]interface{}, len(detections))
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

		//t := (*string)(nil)
		var t string
		if s, ok := cast["type"].(string); ok {
			t = s
		}

		var tn string
		if s, ok := cast["tag_name"].(string); ok {
			tn = s
		}

		alerts = append(alerts, sigsci.Alert{
			AlertUpdateBody: sigsci.AlertUpdateBody{
				LongName:          cast["long_name"].(string),
				Interval:          cast["interval"].(int),
				Threshold:         cast["threshold"].(int),
				SkipNotifications: cast["skip_notifications"].(bool),
				Enabled:           cast["enabled"].(bool),
				Action:            cast["action"].(string),
			},
			Type:    t,
			TagName: tn,
		})
	}
	return alerts
}

func flattenAlerts(alerts []sigsci.Alert) []interface{} {
	var alertsSet = make([]interface{}, len(alerts))
	for i, alert := range alerts {
		alertsSet[i] = map[string]interface{}{
			"id":                 alert.ID,
			"long_name":          alert.LongName,
			"interval":           alert.Interval,
			"threshold":          alert.Threshold,
			"skip_notifications": alert.SkipNotifications,
			"enabled":            alert.Enabled,
			"action":             alert.Action,
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
		a.LongName == b.LongName
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
	for i, _ := range old {
		if old[i].Name != new[i].Name {
			return false
		}
		if old[i].Value != new[i].Value {
			return false
		}
	}
	return true
}

func intArrContains(slice []int, val int) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
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
			return conditions
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
	var conditionsMap = make([]interface{}, len(conditions), len(conditions))
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
		a := sigsci.Action{
			Type: castElement["type"].(string),
		}
		actions = append(actions, a)
	}
	return actions
}

func flattenRuleActions(actions []sigsci.Action) []interface{} {
	var actionsMap = make([]interface{}, len(actions), len(actions))
	for i, action := range actions {
		actionMap := map[string]interface{}{
			"type": action.Type,
		}
		actionsMap[i] = actionMap
	}

	return actionsMap
}

//func resourceCorpImport(corpId string) (corp string, id string, err error) {
//	parts := strings.SplitN(corpId, ":", 2)
//
//	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
//		return "", "", fmt.Errorf("unexpected format of ID (%s), expected corp:id", corpId)
//	}
//
//	return parts[0], parts[1], nil
//}

func resourceSiteImport(corpSiteId string) (site string, id string, err error) {
	parts := strings.SplitN(corpSiteId, ":", 3)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected site:id", corpSiteId)
	}

	return parts[0], parts[1], nil
}
