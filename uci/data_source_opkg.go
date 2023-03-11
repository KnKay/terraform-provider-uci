package uci

import (
	"context"
	"strings"

	"github.com/KnKay/terraform-provider-uci/internal/ssh_helper"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type opkgDataSourceModel struct {
	ID       types.String  `tfsdk:"id"`
	Packages []opkgPackage `tfsdk:"packages"`
}

type opkgPackage struct {
	Name    types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
}

var (
	_ datasource.DataSource              = &opkgDataSource{}
	_ datasource.DataSourceWithConfigure = &opkgDataSource{}
)

// systemDataSource is the data source implementation.
type opkgDataSource struct {
	client *ssh_helper.SshClient
}

func NewOpkgDataSource() datasource.DataSource {
	return &opkgDataSource{}
}

func (d *opkgDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_opkg"
}

// Schema defines the schema for the data source.
func (d *opkgDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of packages.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"packages": schema.ListNestedAttribute{
				Description: "List of packages.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the packages",
							Computed:    true,
						},
						"version": schema.StringAttribute{
							Description: "Version of the packages",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *opkgDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uciConnection).Ssh
}

// Read refreshes the Terraform state with the latest data.
func (d *opkgDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state opkgDataSourceModel
	state.ID = types.StringValue("1")
	lines, err := d.client.RunCommand("opkg list")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Load config",
			err.Error(),
		)
		return
	}
	for _, line := range strings.Split(lines, "\n") {
		if line != "" {
			info := strings.Split(line, " - ")
			state.Packages = append(state.Packages, opkgPackage{Name: types.StringValue(info[0]), Version: types.StringValue(info[1])})
		}

	}
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
