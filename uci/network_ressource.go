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

type networkResourceModel struct {
	ID          types.String `tfsdk:"id"`
	WAN         wanModel     `tfsdk:"wan"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// systemResource is the resource implementation.
type networkResource struct {
	client *uci.SshTree
}

func NewNetworkRessource() resource.Resource {
	return &networkResource{}
}

func (r *networkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (r *networkResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uciConnection).Client
}

// Schema defines the schema for the resource.
func (r *networkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"wan": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"interface": schema.StringAttribute{
						Required: true,
					},
					"proto": schema.StringAttribute{
						Required: true,
					},
					"netmask": schema.StringAttribute{
						Optional: true,
					},
					"gateway": schema.StringAttribute{
						Optional: true,
					},
					"ip": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *networkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan networkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue("1")
	// Update configuration
	r.client.LoadConfig("network", true)
	r.client.Set("network", "wan", "device", plan.WAN.INTERFACE.ValueString())
	r.client.Set("network", "wan", "proto", plan.WAN.PROTO.ValueString())
	proto := plan.WAN.PROTO.ValueString()
	if proto != "dhcp" {
		r.client.Set("network", "wan", "ipaddr", plan.WAN.IP.ValueString())
		r.client.Set("network", "wan", "netmask", plan.WAN.NETMASK.ValueString())
		r.client.Set("network", "wan", "gateway", plan.WAN.GATEWAY.ValueString())
	}

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
func (r *networkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state networkResourceModel

	err := r.client.LoadConfig("network", true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Load config",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("1")

	proto, exist := r.client.Get("network", "wan", "proto")
	if !exist {
		resp.Diagnostics.AddError(
			"Unable to get wan ip",
			err.Error(),
		)
		return
	}
	state.WAN.PROTO = types.StringValue(proto[0])

	interf, exist := r.client.Get("network", "wan", "device")
	if !exist {
		resp.Diagnostics.AddError(
			"Unable to get wan ip",
			err.Error(),
		)
		return
	}
	state.WAN.INTERFACE = types.StringValue(interf[0])

	ip, exist := r.client.Get("network", "wan", "ipaddr")
	if !exist {
		resp.Diagnostics.AddWarning(
			"Unable to get wan ip",
			err.Error(),
		)
	}
	if len(ip) > 0 {
		state.WAN.IP = types.StringValue(ip[0])
	}

	netmask, exist := r.client.Get("network", "wan", "netmask")
	if !exist {
		resp.Diagnostics.AddWarning(
			"Unable to get wan ip",
			err.Error(),
		)
	}
	if len(ip) > 0 {
		state.WAN.NETMASK = types.StringValue(netmask[0])
	}

	gateway, exist := r.client.Get("network", "wan", "ipaddr")
	if !exist {
		resp.Diagnostics.AddWarning(
			"Unable to get wan ip",
			err.Error(),
		)
	}
	if len(ip) > 0 {
		state.WAN.GATEWAY = types.StringValue(gateway[0])
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *networkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// We call the create. As we know the provider is always creating a file. There is no update!
	creq := resource.CreateRequest{
		Config:       req.Config,
		Plan:         req.Plan,
		ProviderMeta: req.ProviderMeta,
	}
	r.Create(ctx, creq, (*resource.CreateResponse)(resp))
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// We can not delete this! It will destroy the system!
	resp.Diagnostics.AddWarning("You can not delete the system!", "Otherwise it will happen a very bad thing!")
}

func (r *networkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
