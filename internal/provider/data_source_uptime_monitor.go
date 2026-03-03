package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &UptimeMonitorDataSource{}

type UptimeMonitorDataSource struct {
	client *apiclient.Client
}

func NewUptimeMonitorDataSource() datasource.DataSource {
	return &UptimeMonitorDataSource{}
}

func (d *UptimeMonitorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uptime_monitor"
}

func (d *UptimeMonitorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Moneat uptime monitor by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the uptime monitor.",
				Required:    true,
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
				Description: "The type of monitor (http, tcp, ping, push).",
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
	}
}

func (d *UptimeMonitorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UptimeMonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UptimeMonitorDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := d.client.GetUptimeMonitor(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading uptime monitor", err.Error())
		return
	}

	config.Name = types.StringValue(monitor.Name)
	config.URL = types.StringValue(monitor.URL)
	config.Type = types.StringValue(monitor.Type)
	config.IntervalSeconds = types.Int64Value(monitor.IntervalSeconds)
	config.Paused = types.BoolValue(monitor.Paused)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
