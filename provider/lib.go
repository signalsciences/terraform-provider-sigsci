package provider

import (
	"github.com/signalsciences/go-sigsci"
	"strings"
)

func corpSiteToId(corp, site string) string {
	return strings.Join([]string{corp, site}, ":")
}

func idToCorpSite(corpsite string) (corp, site string) {
	split := strings.SplitN(corpsite, ":", 2)
	return split[0], split[1]
}

type ProviderMetadata struct {
	Corp   string
	Client sigsci.Client
}
