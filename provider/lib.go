package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
	"strings"
)

func corpSiteToID(corp, site string) string {
	return strings.Join([]string{corp, site}, ":")
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
