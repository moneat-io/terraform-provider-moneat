package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &HostDataSource{}

type HostDataSource struct {
	client *apiclient.Client
}

type HostDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Hostname types.String `tfsdk:"hostname"`
	OS       types.String `tfsdk:"os"`
	Platform types.String `tfsdk:"platform"`
	Status   types.String `tfsdk:"status"`
}

func NewHostDataSource() datasource.DataSource {
	return &HostDataSource{}
}

func (d *HostDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (d *HostDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Moneat monitored host by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the host.",
				Required:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "The hostname.",
				Computed:    true,
			},
			"os": schema.StringAttribute{
				Description: "The operating system.",
				Computed:    true,
			},
			"platform": schema.StringAttribute{
				Description: "The platform.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The host status.",
				Computed:    true,
			},
		},
	}
}

func (d *HostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *apiclient.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *HostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config HostDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, err := d.client.GetHost(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading host", err.Error())
		return
	}

	config.Hostname = types.StringValue(host.Hostname)
	config.OS = types.StringValue(host.OS)
	config.Platform = types.StringValue(host.Platform)
	config.Status = types.StringValue(host.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
