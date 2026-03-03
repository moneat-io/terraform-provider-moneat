package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &DashboardsDataSource{}

type DashboardsDataSource struct {
	client *apiclient.Client
}

type DashboardDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type DashboardsDataSourceModel struct {
	Dashboards []DashboardDataSourceModel `tfsdk:"dashboards"`
}

func NewDashboardsDataSource() datasource.DataSource {
	return &DashboardsDataSource{}
}

func (d *DashboardsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboards"
}

func (d *DashboardsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list all Moneat dashboards.",
		Attributes: map[string]schema.Attribute{
			"dashboards": schema.ListNestedAttribute{
				Description: "List of dashboards.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the dashboard.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *DashboardsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DashboardsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	dashboards, err := d.client.ListDashboards()
	if err != nil {
		resp.Diagnostics.AddError("Error listing dashboards", err.Error())
		return
	}

	var state DashboardsDataSourceModel
	for _, db := range dashboards {
		state.Dashboards = append(state.Dashboards, DashboardDataSourceModel{
			ID:          types.StringValue(db.ID),
			Name:        types.StringValue(db.Name),
			Description: types.StringValue(db.Description),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
