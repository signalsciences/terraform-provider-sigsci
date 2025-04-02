package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceCorpUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceCorpUserCreate,
		Read:   resourceCorpUserRead,
		Update: resourceCorpUserUpdate,
		Delete: resourceCorpUserDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "user",
			},
			"memberships": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceCorpUserCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	client := pm.Client
	corp := pm.Corp

	email := d.Get("email").(string)
	role := d.Get("role").(string)

	// Build site memberships
	rawMemberships := d.Get("memberships").([]interface{})
	var memberships []sigsci.SiteMembership
	for _, m := range rawMemberships {
		mem := m.(map[string]interface{})
		memberships = append(memberships, sigsci.NewSiteMembership(
			mem["site_name"].(string),
			sigsci.Role(mem["role"].(string)),
		))
	}

	invite := sigsci.NewCorpUserInvite(sigsci.Role(role), memberships)

	if _, err := client.InviteUser(corp, email, invite); err != nil {
		return fmt.Errorf("error inviting corp user: %w", err)
	}

	d.SetId(email)
	return resourceCorpUserRead(d, m)
}

func resourceCorpUserRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	client := pm.Client
	corp := pm.Corp

	email := d.Id()
	user, err := client.GetCorpUser(corp, email)
	if err != nil {
		d.SetId("")
    fmt.Println(err)
		return nil
	}

	d.Set("email", user.Email)
	d.Set("role", user.Role)

	// Flatten memberships
	var memberships []map[string]interface{}
	for siteName, role := range user.Memberships {
		memberships = append(memberships, map[string]interface{}{
			"site_name": siteName,
			"role":      role,
		})
	}
	d.Set("memberships", memberships)

	return nil
}

func resourceCorpUserUpdate(d *schema.ResourceData, m interface{}) error {
	// The Signal Sciences API does not support updating user roles or memberships directly.
	// To change a user's role or memberships, they must be removed and re-invited with the new settings.
	return fmt.Errorf("updating corp users is not supported by the Signal Sciences API")
}

func resourceCorpUserDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	client := pm.Client
	corp := pm.Corp

	email := d.Id()

	if err := client.DeleteCorpUser(corp, email); err != nil {
		return fmt.Errorf("error deleting corp user: %w", err)
	}

	d.SetId("")
	return nil
}
