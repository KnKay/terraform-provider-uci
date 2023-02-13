package uci

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSystemResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "uci_system" "test" {
  hostname = "test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("uci_system.test", "hostname", "test"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("uci_system.test", "id"),
					resource.TestCheckResourceAttrSet("uci_system.test", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "uci_system.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "uci_system" "test" {
	hostname = "OpentWrt"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("uci_system.test", "hostname", "OpentWrt"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
