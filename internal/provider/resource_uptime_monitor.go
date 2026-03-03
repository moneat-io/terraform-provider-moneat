package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &UptimeMonitorResource{}
var _ resource.ResourceWithImportState = &UptimeMonitorResource{}

// UptimeMonitorResource defines the resource implementation.
type UptimeMonitorResource struct {
	client *apiclient.Client
}

// UptimeMonitorResourceModel describes the resource data model.
type UptimeMonitorResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	URL             types.String `tfsdk:"url"`
	Type            types.String `tfsdk:"type"`
	IntervalSeconds types.Int64  `tfsdk:"interval_seconds"`
	Paused          types.Bool   `tfsdk:"paused"`
}

// NewUptimeMonitorResource returns a new resource factory function.
func NewUptimeMonitorResource() resource.Resource {
	return &UptimeMonitorResource{}
}

func (r *UptimeMonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uptime_monitor"
}

func (r *UptimeMonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat uptime monitor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the uptime monitor.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the uptime monitor.",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL or address to monitor.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of monitor (http, tcp, ping, push).",
				Required:    true,
			},
			"interval_seconds": schema.Int64Attribute{
				Description: "The check interval in seconds.",
				Required:    true,
			},
			"paused": schema.BoolAttribute{
				Description: "Whether the monitor is paused.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *UptimeMonitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *apiclient.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *UptimeMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UptimeMonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateUptimeMonitorRequest{
		Name:            plan.Name.ValueString(),
		URL:             plan.URL.ValueString(),
		Type:            plan.Type.ValueString(),
		IntervalSeconds: plan.IntervalSeconds.ValueInt64(),
	}

	monitor, err := r.client.CreateUptimeMonitor(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating uptime monitor", err.Error())
		return
	}

	plan.ID = types.StringValue(monitor.ID)
	plan.Name = types.StringValue(monitor.Name)
	plan.URL = types.StringValue(monitor.URL)
	plan.Type = types.StringValue(monitor.Type)
	plan.IntervalSeconds = types.Int64Value(monitor.IntervalSeconds)
	plan.Paused = types.BoolValue(monitor.Paused)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UptimeMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UptimeMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.client.GetUptimeMonitor(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading uptime monitor", err.Error())
		return
	}

	state.Name = types.StringValue(monitor.Name)
	state.URL = types.StringValue(monitor.URL)
	state.Type = types.StringValue(monitor.Type)
	state.IntervalSeconds = types.Int64Value(monitor.IntervalSeconds)
	state.Paused = types.BoolValue(monitor.Paused)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UptimeMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UptimeMonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UptimeMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateUptimeMonitorRequest{
		Name:            plan.Name.ValueString(),
		URL:             plan.URL.ValueString(),
		Type:            plan.Type.ValueString(),
		IntervalSeconds: plan.IntervalSeconds.ValueInt64(),
		Paused:          plan.Paused.ValueBool(),
	}

	monitor, err := r.client.UpdateUptimeMonitor(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating uptime monitor", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(monitor.Name)
	plan.URL = types.StringValue(monitor.URL)
	plan.Type = types.StringValue(monitor.Type)
	plan.IntervalSeconds = types.Int64Value(monitor.IntervalSeconds)
	plan.Paused = types.BoolValue(monitor.Paused)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UptimeMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UptimeMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUptimeMonitor(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting uptime monitor", err.Error())
		return
	}
}

func (r *UptimeMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	monitor, err := r.client.GetUptimeMonitor(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing uptime monitor", err.Error())
		return
	}

	state := UptimeMonitorResourceModel{
		ID:              types.StringValue(monitor.ID),
		Name:            types.StringValue(monitor.Name),
		URL:             types.StringValue(monitor.URL),
		Type:            types.StringValue(monitor.Type),
		IntervalSeconds: types.Int64Value(monitor.IntervalSeconds),
		Paused:          types.BoolValue(monitor.Paused),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
