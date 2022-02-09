package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()

	testAccProviders = map[string]*schema.Provider{
		"kafkamanager": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"kafkamanager": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}
func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("KAFKAMANAGER_URL"); v == "" {
		t.Fatal("KAFKAMANAGER_URL must be set for acceptance tests")
	}
	if v := os.Getenv("KAFKAMANAGER_KEY"); v == "" {
		t.Fatal("KAFKAMANAGER_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("KAFKAMANAGER_OKTA_GROUPS"); v == "" {
		t.Fatal("KAFKAMANAGER_OKTA_GROUPS must be set for acceptance tests")
	}
	if v := os.Getenv("KAFKAMANAGER_SUPPLIER"); v == "" {
		t.Fatal("KAFKAMANAGER_SUPPLIER must be set for acceptance tests")
	}
	if v := os.Getenv("KAFKAMANAGER_USER_ID"); v == "" {
		t.Fatal("KAFKAMANAGER_USER_ID must be set for acceptance tests")
	}

	err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}
