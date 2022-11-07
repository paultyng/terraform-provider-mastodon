package provider

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mattn/go-mastodon"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &MastodonProvider{}
var _ provider.ProviderWithMetadata = &MastodonProvider{}

// MastodonProvider defines the provider implementation.
type MastodonProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MastodonProviderModel describes the provider data model.
type MastodonProviderModel struct {
	Server       types.String `tfsdk:"server"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	Insecure     types.Bool   `tfsdk:"allow_insecure"`
}

func (p *MastodonProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mastodon"
	resp.Version = p.version
}

func (p *MastodonProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"server": {
				MarkdownDescription: "The server to connect to.",
				Required:            true,
				Type:                types.StringType,
			},
			"client_id": {
				MarkdownDescription: "The client ID of the application.",
				Required:            true,
				Type:                types.StringType,
			},
			"client_secret": {
				MarkdownDescription: "The client secret of the application.",
				Required:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
			"username": {
				MarkdownDescription: "The user with which to login.",
				Required:            true,
				Type:                types.StringType,
			},
			"password": {
				MarkdownDescription: "The password of the user.",
				Required:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},

			"allow_insecure": {
				MarkdownDescription: "Allow invalid certificates on the Mastodon server.",
				Optional:            true,
				Type:                types.BoolType,
			},
		},
	}, nil
}

func (p *MastodonProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MastodonProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := mastodon.NewClient(&mastodon.Config{
		Server:       data.Server.ValueString(),
		ClientID:     data.ClientID.ValueString(),
		ClientSecret: data.ClientSecret.ValueString(),
	})

	insecure := data.Insecure.ValueBool()

	client.Client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
	}

	err := client.Authenticate(ctx, data.Username.ValueString(), data.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to authenticate", err.Error())
		return
	}

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *MastodonProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFollowResource,
	}
}

func (p *MastodonProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MastodonProvider{
			version: version,
		}
	}
}
