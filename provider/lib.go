package provider

import "strings"

func corpSiteToId(corp, site string) string {
	return strings.Join([]string{corp, site}, ":")
}

func idToCorpSite(corpsite string) (corp, site string) {
	split := strings.SplitN(corpsite, ":", 2)
	return split[0], split[1]
}
