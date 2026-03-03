package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &UptimeMonitorsDataSource{}

type UptimeMonitorsDataSource struct {
	client *apiclient.Client
}

type UptimeMonitorDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	URL             types.String `tfsdk:"url"`
	Type            types.String `tfsdk:"type"`
	IntervalSeconds types.Int64  `tfsdk:"interval_seconds"`
	Paused          types.Bool   `tfsdk:"paused"`
}

type UptimeMonitorsDataSourceModel struct {
	Monitors []UptimeMonitorDataSourceModel `tfsdk:"monitors"`
}

func NewUptimeMonitorsDataSource() datasource.DataSource {
	return &UptimeMonitorsDataSource{}
}

func (d *UptimeMonitorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uptime_monitors"
}

func (d *UptimeMonitorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list all Moneat uptime monitors.",
		Attributes: map[string]schema.Attribute{
			"monitors": schema.ListNestedAttribute{
				Description: "List of uptime monitors.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the uptime monitor.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the uptime monitor.",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "The URL or address being monitored.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of monitor.",
							Computed:    true,
						},
						"interval_seconds": schema.Int64Attribute{
							Description: "The check interval in seconds.",
							Computed:    true,
						},
						"paused": schema.BoolAttribute{
							Description: "Whether the monitor is paused.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *UptimeMonitorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UptimeMonitorsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	monitors, err := d.client.ListUptimeMonitors()
	if err != nil {
		resp.Diagnostics.AddError("Error listing uptime monitors", err.Error())
		return
	}

	var state UptimeMonitorsDataSourceModel
	for _, m := range monitors {
		state.Monitors = append(state.Monitors, UptimeMonitorDataSourceModel{
			ID:              types.StringValue(m.ID),
			Name:            types.StringValue(m.Name),
			URL:             types.StringValue(m.URL),
			Type:            types.StringValue(m.Type),
			IntervalSeconds: types.Int64Value(m.IntervalSeconds),
			Paused:          types.BoolValue(m.Paused),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
