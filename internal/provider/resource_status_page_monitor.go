package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &StatusPageMonitorResource{}
var _ resource.ResourceWithImportState = &StatusPageMonitorResource{}

type StatusPageMonitorResource struct {
	client *apiclient.Client
}

type StatusPageMonitorResourceModel struct {
	ID          types.String `tfsdk:"id"`
	PageID      types.String `tfsdk:"page_id"`
	MonitorID   types.String `tfsdk:"monitor_id"`
	DisplayName types.String `tfsdk:"display_name"`
	SortOrder   types.Int64  `tfsdk:"sort_order"`
}

func NewStatusPageMonitorResource() resource.Resource {
	return &StatusPageMonitorResource{}
}

func (r *StatusPageMonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page_monitor"
}

func (r *StatusPageMonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a monitor linked to a Moneat status page.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the status page monitor link.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"page_id": schema.StringAttribute{
				Description: "The ID of the status page.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"monitor_id": schema.StringAttribute{
				Description: "The ID of the uptime monitor to link.",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name for this monitor on the status page.",
				Required:    true,
			},
			"sort_order": schema.Int64Attribute{
				Description: "The sort order for this monitor on the status page.",
				Required:    true,
			},
		},
	}
}

func (r *StatusPageMonitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StatusPageMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StatusPageMonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateStatusPageMonitorRequest{
		MonitorID:   plan.MonitorID.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		SortOrder:   plan.SortOrder.ValueInt64(),
	}

	monitor, err := r.client.CreateStatusPageMonitor(plan.PageID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating status page monitor", err.Error())
		return
	}

	plan.ID = types.StringValue(monitor.ID)
	plan.MonitorID = types.StringValue(monitor.MonitorID)
	plan.DisplayName = types.StringValue(monitor.DisplayName)
	plan.SortOrder = types.Int64Value(monitor.SortOrder)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StatusPageMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.client.GetStatusPageMonitor(state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading status page monitor", err.Error())
		return
	}

	state.MonitorID = types.StringValue(monitor.MonitorID)
	state.DisplayName = types.StringValue(monitor.DisplayName)
	state.SortOrder = types.Int64Value(monitor.SortOrder)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StatusPageMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StatusPageMonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state StatusPageMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateStatusPageMonitorRequest{
		DisplayName: plan.DisplayName.ValueString(),
		SortOrder:   plan.SortOrder.ValueInt64(),
	}

	monitor, err := r.client.UpdateStatusPageMonitor(state.PageID.ValueString(), state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating status page monitor", err.Error())
		return
	}

	plan.ID = state.ID
	plan.PageID = state.PageID
	plan.MonitorID = types.StringValue(monitor.MonitorID)
	plan.DisplayName = types.StringValue(monitor.DisplayName)
	plan.SortOrder = types.Int64Value(monitor.SortOrder)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StatusPageMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStatusPageMonitor(state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting status page monitor", err.Error())
		return
	}
}

func (r *StatusPageMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: page_id/monitor_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: page_id/monitor_id",
		)
		return
	}

	pageID := parts[0]
	monitorLinkID := parts[1]

	monitor, err := r.client.GetStatusPageMonitor(pageID, monitorLinkID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing status page monitor", err.Error())
		return
	}

	state := StatusPageMonitorResourceModel{
		ID:          types.StringValue(monitor.ID),
		PageID:      types.StringValue(pageID),
		MonitorID:   types.StringValue(monitor.MonitorID),
		DisplayName: types.StringValue(monitor.DisplayName),
		SortOrder:   types.Int64Value(monitor.SortOrder),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
