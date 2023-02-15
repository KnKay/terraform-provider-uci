package uci

import (
	"context"

	"github.com/digineo/go-uci"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// We will get useful information of the WAN module.
// This can be used to configure other devices for things like VPN.
type networkDataSourceModel struct {
	ID  types.Int64 `tfsdk:"id"`
	WAN wanModel    `tfsdk:"wan"`
}

type wanModel struct {
	IP        types.String `tfsdk:"ip"`
	INTERFACE types.String `tfsdk:"interface"`
	PROTO     types.String `tfsdk:"proto"`
	NETMASK   types.String `tfsdk:"netmask"`
	GATEWAY   types.String `tfsdk:"gateway"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &networkDataSource{}
	_ datasource.DataSourceWithConfigure = &networkDataSource{}
)

// systemDataSource is the data source implementation.
type networkDataSource struct {
	client *uci.SshTree
}

func NewNetworkDataSource() datasource.DataSource {
	return &networkDataSource{}
}

// Metadata returns the data source type name.
func (d *networkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

// Schema defines the schema for the data source.
func (d *networkDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"wan": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"ip": schema.StringAttribute{
						Computed: true,
					},
					"interface": schema.StringAttribute{
						Computed: true,
					},
					"proto": schema.StringAttribute{
						Computed: true,
					},
					"netmask": schema.StringAttribute{
						Computed: true,
					},
					"gateway": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *networkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state networkDataSourceModel

	err := d.client.LoadConfig("network", true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Load config",
			err.Error(),
		)
		return
	}
	state.ID = types.Int64Value(int64(1))

	proto, exist := d.client.Get("network", "wan", "proto")
	if !exist {
		resp.Diagnostics.AddError(
			"Unable to get wan ip",
			err.Error(),
		)
		return
	}
	state.WAN.PROTO = types.StringValue(proto[0])

	interf, exist := d.client.Get("network", "wan", "device")
	if !exist {
		resp.Diagnostics.AddError(
			"Unable to get wan ip",
			err.Error(),
		)
		return
	}
	state.WAN.INTERFACE = types.StringValue(interf[0])

	ip, exist := d.client.Get("network", "wan", "ip")
	if !exist {
		resp.Diagnostics.AddWarning(
			"Unable to get wan ip",
			err.Error(),
		)
	}
	if len(ip) == 0 {
		state.WAN.IP = types.StringValue("unknown")
	} else {
		state.WAN.IP = types.StringValue(ip[0])
	}

	netmask, exist := d.client.Get("network", "wan", "ip")
	if !exist {
		resp.Diagnostics.AddWarning(
			"Unable to get wan ip",
			err.Error(),
		)
	}
	if len(ip) == 0 {
		state.WAN.NETMASK = types.StringValue("unknown")
	} else {
		state.WAN.NETMASK = types.StringValue(netmask[0])
	}

	gateway, exist := d.client.Get("network", "wan", "ip")
	if !exist {
		resp.Diagnostics.AddWarning(
			"Unable to get wan ip",
			err.Error(),
		)
	}
	if len(ip) == 0 {
		state.WAN.GATEWAY = types.StringValue("unknown")
	} else {
		state.WAN.GATEWAY = types.StringValue(gateway[0])
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *networkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uci.SshTree)
}
