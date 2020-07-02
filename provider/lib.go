package provider

import (
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
