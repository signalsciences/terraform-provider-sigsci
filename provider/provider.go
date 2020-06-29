package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider is the Signalsciences terraform provider, returns a terraform.ResourceProvider
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"corp": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Corp short name (id)",
			},
			"auth_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email to be used for authentication",
				DefaultFunc: schema.EnvDefaultFunc("SIGSCI_EMAIL", nil),
			},
			"auth_password": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("SIGSCI_PASSWORD", nil),
				Description:  "The password used to for authentication specify either the password or the token",
				Sensitive:    true,
				AtLeastOneOf: []string{"auth_password", "auth_token"},
			},
			"auth_token": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("SIGSCI_TOKEN", nil),
				Description:  "The token used for authentication specify either the password or the token",
				Sensitive:    true,
				AtLeastOneOf: []string{"auth_password", "auth_token"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sigsci_site": resourceSite(),
		},
	}
	provider.ConfigureFunc = providerConfigure()
	return provider
}

func providerConfigure() schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := Config{
			Email:    d.Get("auth_email").(string),
			Password: d.Get("auth_password").(string),
			APIToken: d.Get("auth_token").(string),
		}
		client, err := config.Client()
		if err != nil {
			return nil, err
		}
		return client, nil
	}
}
