package pkg

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func TestProviderConfigure(t *testing.T) {
	testAccPreCheck(t)
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OPENWRT_HOST"); v == "" {
		t.Fatal("TF_VAR_VSPHERE_USER must be set for acceptance tests")
	}

	if v := os.Getenv("OPENWRT_USERNAME"); v == "" {
		t.Fatal("TF_VAR_VSPHERE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("OPENWRT_PASSWORD"); v == "" {
		t.Fatal("TF_VAR_VSPHERE_SERVER must be set for acceptance tests")
	}
}
