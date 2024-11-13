package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TODO implement sweepers for everyone
func TestAccResourceSiteSignalTagCRUD(t *testing.T) {
	t.Parallel()
	resourceName := "sigsci_site_signal_tag.test_tag"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "sigsci_site_signal_tag" "test_tag" {
                      site_short_name = "%s" 
                      name            = "My new tag"
                      description     = "test description"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "My new tag"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttrSet(resourceName, "configurable"),
					resource.TestCheckResourceAttrSet(resourceName, "informational"),
					resource.TestCheckResourceAttrSet(resourceName, "needs_response"),
				),
			},
			{
				Config: fmt.Sprintf(`
                     resource "sigsci_site_signal_tag" "test_tag" {
                      site_short_name = "%s" 
                      name            = "My new tag"
                      description     = "test description 2"
				}`, testSite),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_short_name", testSite),
					resource.TestCheckResourceAttr(resourceName, "name", "My new tag"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description 2"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: fmt.Sprintf("%s:", testSite),
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateCheck:    testAccImportStateCheckFunction(1),
			},
		},
	})
}

func TestResourceSiteSignalTagNameValidation(t *testing.T) {
	cases := []struct {
		value    string
		expected bool
	}{
		{"s", true},
		{"valid-name", false},
		{"this-name-is-way-too-long-for-the-validation-rules", true},
	}

	resource := resourceSiteSignalTag()
	nameSchema := resource.Schema["name"]

	for _, tc := range cases {
		_, errors := nameSchema.ValidateFunc(tc.value, "name")

		if tc.expected && len(errors) == 0 {
			t.Errorf("Expected an error for value '%s', but got none", tc.value)
		}

		if !tc.expected && len(errors) > 0 {
			t.Errorf("Did not expect an error for value '%s', but got: %v", tc.value, errors)
		}
	}
}

func TestResourceSiteSignalTagDescriptionValidation(t *testing.T) {
	cases := []struct {
		value    string
		expected bool
	}{
		{"", false},
		{"valid-description", false},
		{"this-is-way-too-long-for-the-validation-rules-and-this-is-way-too-long-for-the-validation-rules-and-this-is-way-too-long-for-the-validation-rules", true},
	}

	resource := resourceSiteSignalTag()
	nameSchema := resource.Schema["description"]

	for _, tc := range cases {
		_, errors := nameSchema.ValidateFunc(tc.value, "description")

		if tc.expected && len(errors) == 0 {
			t.Errorf("Expected an error for value '%s', but got none", tc.value)
		}

		if !tc.expected && len(errors) > 0 {
			t.Errorf("Did not expect an error for value '%s', but got: %v", tc.value, errors)
		}
	}
}
