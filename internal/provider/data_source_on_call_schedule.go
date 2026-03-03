package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &OnCallScheduleDataSource{}

type OnCallScheduleDataSource struct {
	client *apiclient.Client
}

type OnCallScheduleDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Timezone     types.String `tfsdk:"timezone"`
	RotationType types.String `tfsdk:"rotation_type"`
}

func NewOnCallScheduleDataSource() datasource.DataSource {
	return &OnCallScheduleDataSource{}
}

func (d *OnCallScheduleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_on_call_schedule"
}

func (d *OnCallScheduleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Moneat on-call schedule by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the on-call schedule.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the on-call schedule.",
				Computed:    true,
			},
			"timezone": schema.StringAttribute{
				Description: "The timezone for the schedule.",
				Computed:    true,
			},
			"rotation_type": schema.StringAttribute{
				Description: "The rotation type (daily, weekly, custom).",
				Computed:    true,
			},
		},
	}
}

func (d *OnCallScheduleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallScheduleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OnCallScheduleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schedule, err := d.client.GetOnCallSchedule(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading on-call schedule", err.Error())
		return
	}

	config.Name = types.StringValue(schedule.Name)
	config.Timezone = types.StringValue(schedule.Timezone)
	config.RotationType = types.StringValue(schedule.RotationType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
