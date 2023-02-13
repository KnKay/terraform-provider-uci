package uci

import (
	"context"

	"github.com/digineo/go-uci"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &systemDataSource{}
	_ datasource.DataSourceWithConfigure = &systemDataSource{}
)

// NewsystemDataSource is a helper function to simplify the provider implementation.
func NewSystemDataSource() datasource.DataSource {
	return &systemDataSource{}
}

// systemDataSource is the data source implementation.
type systemDataSource struct {
	client *uci.SshTree
}

// DataSourceModel maps the data source schema data.
type systemDataSourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	Hostname types.String `tfsdk:"hostname"`
}

// Metadata returns the data source type name.
func (d *systemDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

// Schema defines the schema for the data source.
func (d *systemDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"hostname": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *systemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state systemDataSourceModel

	err := d.client.LoadConfig("system", true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Load config",
			err.Error(),
		)
		return
	}
	state.ID = types.Int64Value(int64(1))
	hostname, exist := d.client.Get("system", "@system[0]", "hostname")
	if !exist {
		resp.Diagnostics.AddError(
			"Unable to get hostname",
			err.Error(),
		)
		return
	}
	state.Hostname = types.StringValue(hostname[0])

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *systemDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uci.SshTree)
}
