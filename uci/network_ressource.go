package uci

import "github.com/hashicorp/terraform-plugin-framework/types"

type networkResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	WAN_INTERFACE types.String `tfsdk:"wan_interface"`
	WAN_PROTO     types.String `tfsdk:"wan_proto"`
	WAN_IP        types.String `tfsdk:"wan_ip"`
	WAN_NETMASK   types.String `tfsdk:"wan_netmasl"`
}
