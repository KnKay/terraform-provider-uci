package pkg

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDatasourceSystem(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckHashicupsOrderDestroy,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.system.test", "hostname", strings.TrimSpace(string("openwrt"))),
				),
			},
		},
	})
}
