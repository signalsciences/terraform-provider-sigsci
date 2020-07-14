package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/signalsciences/go-sigsci"
)

// Provider is the Signalsciences terraform provider, returns a terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"corp": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Corp short name (id)",
				DefaultFunc: schema.EnvDefaultFunc("SIGSCI_CORP", nil),
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email to be used for authentication",
				DefaultFunc: schema.EnvDefaultFunc("SIGSCI_EMAIL", nil),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("SIGSCI_PASSWORD", nil),
				Description:  "The password used to for authentication specify either the password or the token",
				Sensitive:    true,
				AtLeastOneOf: []string{"password", "auth_token"},
			},
			"auth_token": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("SIGSCI_TOKEN", nil),
				Description:  "The token used for authentication specify either the password or the token",
				Sensitive:    true,
				AtLeastOneOf: []string{"password", "auth_token"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sigsci_site":            resourceSite(),
			"sigsci_site_list":       resourceSiteList(),
			"sigsci_site_signal_tag": resourceSiteSignalTag(),
			"sigsci_site_redaction":  resourceSiteRedaction(),
			"sigsci_site_alert":      resourceSiteAlert(),
		},
	}
	provider.ConfigureFunc = providerConfigure()
	return provider
}

func providerConfigure() schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := Config{
			Email:    d.Get("email").(string),
			Password: d.Get("password").(string),
			APIToken: d.Get("auth_token").(string),
		}
		client, err := config.Client()
		if err != nil {
			return nil, err
		}
		return providerMetadata{
			Corp:   d.Get("corp").(string),
			Client: client.(sigsci.Client),
		}, nil
	}
}
