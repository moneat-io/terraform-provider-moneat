package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &DashboardDataSource{}

type DashboardDataSource struct {
	client *apiclient.Client
}

func NewDashboardDataSource() datasource.DataSource {
	return &DashboardDataSource{}
}

func (d *DashboardDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func (d *DashboardDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Moneat dashboard by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the dashboard.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the dashboard.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the dashboard.",
				Computed:    true,
			},
		},
	}
}

func (d *DashboardDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DashboardDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dashboard, err := d.client.GetDashboard(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading dashboard", err.Error())
		return
	}

	config.Name = types.StringValue(dashboard.Name)
	config.Description = types.StringValue(dashboard.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
