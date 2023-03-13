package uci

import (
	"context"
	"strings"
	"time"

	"github.com/KnKay/terraform-provider-uci/internal/ssh_helper"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &opkgResource{}
	_ resource.ResourceWithConfigure = &opkgResource{}
)

type opgRessourceModel struct {
	ID          types.String  `tfsdk:"id"`
	Packages    []opkgPackage `tfsdk:"packages"`
	LastUpdated types.String  `tfsdk:"last_updated"`
}

type opkgResource struct {
	client *ssh_helper.SshClient
}

func NewOpkgRessource() resource.Resource {
	return &opkgResource{}
}

func (r *opkgResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_opkg"
}

func (r *opkgResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uciConnection).Ssh
}

func (r *opkgResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of packages.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"packages": schema.ListNestedAttribute{
				Description: "List of packages.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the packages",
							Required:    true,
						},
						"version": schema.StringAttribute{
							Description: "Version of the packages",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *opkgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan opgRessourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.RunCommand("opkg update")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update opkg",
			err.Error(),
		)
		return
	}

	pkgs := ""
	for _, pkg := range plan.Packages {
		pkgs = pkgs + " " + pkg.Name.String()
	}

	_, err = r.client.RunCommand("opkg install " + pkgs)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update opkg",
			err.Error(),
		)
		return
	}
	plan.ID = types.StringValue("1")

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *opkgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// We get the state prior read. This seems to be the plan
	var plan opgRessourceModel
	req.State.Get(ctx, &plan)

	var state opgRessourceModel
	state.ID = types.StringValue("1")
	lines, err := r.client.RunCommand("opkg list-installed")
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
			for _, pkg := range plan.Packages {
				if strings.Contains(pkg.Name.String(), info[0]) {
					state.Packages = append(state.Packages, opkgPackage{Name: types.StringValue(info[0])})
				}
			}
		}
	}
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *opkgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// We call the create. As we know the provider is always creating a file. There is no update!
	creq := resource.CreateRequest{
		Config:       req.Config,
		Plan:         req.Plan,
		ProviderMeta: req.ProviderMeta,
	}
	r.Create(ctx, creq, (*resource.CreateResponse)(resp))
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *opkgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state opgRessourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pkgs := ""
	for _, pkg := range state.Packages {
		pkgs = pkgs + " " + pkg.Name.String()
	}

	_, err := r.client.RunCommand("opkg install " + pkgs)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update opkg",
			err.Error(),
		)
		return
	}
}

func (r *opkgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
