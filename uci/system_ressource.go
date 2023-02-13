package uci

import (
	"context"
	"time"

	"github.com/digineo/go-uci"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &systemResource{}
	_ resource.ResourceWithConfigure = &systemResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewSystemResource() resource.Resource {
	return &systemResource{}
}

// systemResource is the resource implementation.
type systemResource struct {
	client *uci.SshTree
}

type systemResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Hostname    types.String `tfsdk:"hostname"`
}

// Metadata returns the resource type name.
func (r *systemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

// Configure adds the provider configured client to the resource.
func (r *systemResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*uci.SshTree)
}

// Schema defines the schema for the resource.
func (r *systemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"hostname": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *systemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan systemResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue("1")
	// Update configuration
	r.client.LoadConfig("system", true)
	r.client.Set("system", "@system[0]", "hostname", plan.Hostname.ValueString())

	// Write Config
	err := r.client.Commit()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get hostname",
			err.Error(),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *systemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state systemResourceModel

	err := r.client.LoadConfig("system", true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Load config",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("1")
	hostname, exist := r.client.Get("system", "@system[0]", "hostname")
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

// Update updates the resource and sets the updated Terraform state on success.
func (r *systemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// We call the create. As we know the provider is always creating a file. There is no update!
	creq := resource.CreateRequest{
		Config:       req.Config,
		Plan:         req.Plan,
		ProviderMeta: req.ProviderMeta,
	}
	r.Create(ctx, creq, (*resource.CreateResponse)(resp))
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *systemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// We can not delete this! It will destroy the system!
	resp.Diagnostics.AddWarning("You can not delete the system!", "Otherwise it will happen a very bad thing!")
}

func (r *systemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
