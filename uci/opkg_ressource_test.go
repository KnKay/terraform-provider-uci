package uci

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOpkgResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "uci_opkg" "test" {
  packages = [
	{
		name = "nano"
	}
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("uci_opkg.test", "packages.0.name", "nano"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      "uci_opkg.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// The last_updated attribute does not exist in the HashiCups
			// 	// API, therefore there is no value for it during import.
			// 	ImportStateVerifyIgnore: []string{"last_updated"},
			// },
			// Update and Read testing
			// 			{
			// 				Config: providerConfig + `
			// resource "uci_opkg" "test" {
			// 	packages = [
			// 		{
			// 			name = "nano"
			// 		}
			// 	]
			// }
			// `,
			// 				Check: resource.ComposeAggregateTestCheckFunc(
			// 					// Verify first order item updated
			// 					resource.TestCheckResourceAttr("uci_opkg.test", "packages.0.name", "nano"),
			// 				),
			// 			},
			// 			{
			// 				Config: providerConfig + `
			// resource "uci_opkg" "test" {
			// 	packages = [
			// 		{
			// 			name = "nano"
			// 		}
			// 	]
			// }
			// `,
			// 				Check: resource.ComposeAggregateTestCheckFunc(
			// 					// Verify first order item updated
			// 					resource.TestCheckResourceAttr("uci_opkg.test", "wan.interface", "eth0"),
			// 				),
			// 			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
