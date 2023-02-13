package uci

import "github.com/hashicorp/terraform-plugin-framework/types"

// We will get useful information of the WAN module.
// This can be used to configure other devices for things like VPN.
type networkDataSourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	WAN_IP        types.String `tfsdk:"wan_ip"`
	WAN_INTERFACE types.String `tfsdk:"wan_interface"`
}
