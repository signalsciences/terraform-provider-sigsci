package provider

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	rand.Seed(time.Now().UnixNano())

	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sigsci": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := testAccProvider.InternalValidate(); err != nil {
		t.Fatalf("err %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SIGSCI_CORP"); v == "" {
		t.Fatal("SIGSCI_CORP must be set for acceptance tests")
	}

	if v := os.Getenv("SIGSCI_EMAIL"); v == "" {
		t.Fatal("SIGSCI_EMAIL must be set for acceptance tests")
	}

	pass := os.Getenv("SIGSCI_PASSWORD")
	token := os.Getenv("SIGSCI_TOKEN")

	if pass == "" && token == "" {
		t.Fatal("SIGSCI_PASSWORD or SIGSCI_TOKEN must be set for acceptance tests")
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var testSite = "test" // createTestSiteName()
