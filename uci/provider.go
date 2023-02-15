package uci

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"

	gouci "github.com/digineo/go-uci"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &uciProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &uciProvider{}
}

type uciProvider struct{}

// hashicupsProviderModel maps provider schema data to a Go type.
type uciProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *uciProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uci"
}

// Schema defines the provider-level schema for configuration data.
func (p *uciProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a uciAPI client for data sources and resources.
func (p *uciProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config uciProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown uciAPI Host",
			"The provider cannot create the uci API client as there is an unknown configuration value for the uci API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the UCI_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown uci API Username",
			"The provider cannot create the uci API client as there is an unknown configuration value for the uci API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the UCI_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown uci API Password",
			"The provider cannot create the uci API client as there is an unknown configuration value for the uci API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the UCI_PASSWORD environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("UCI_HOST")
	username := os.Getenv("UCI_USERNAME")
	password := os.Getenv("UCI_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing uci API Host",
			"The provider cannot create the uci API client as there is a missing or empty value for the uci API host. "+
				"Set the host value in the configuration or use the HASHICUPS_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing uci API Username",
			"The provider cannot create the uci API client as there is a missing or empty value for the uci API username. "+
				"Set the username value in the configuration or use the HASHICUPS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing uci API Password",
			"The provider cannot create the uci API client as there is a missing or empty value for the uci API password. "+
				"Set the password value in the configuration or use the HASHICUPS_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	//Create ssh config
	conf := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// Non-production only
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := gouci.NewSshTree(conf, host)
	// // Create a new client using the configuration values
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create uci API Client",
			"An unexpected error occurred when creating the uci API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"uci Client Error: "+err.Error(),
		)
		return
	}
	// Make the uci client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

}

// DataSources defines the data sources implemented in the provider.
func (p *uciProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSystemDataSource,
		NewNetworkDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *uciProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSystemResource,
	}
}
