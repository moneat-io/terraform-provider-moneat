package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &SystemAlertResource{}
var _ resource.ResourceWithImportState = &SystemAlertResource{}

// SystemAlertResource defines the resource implementation.
type SystemAlertResource struct {
	client *apiclient.Client
}

// SystemAlertResourceModel describes the resource data model.
type SystemAlertResourceModel struct {
	ID              types.String  `tfsdk:"id"`
	SystemID        types.String  `tfsdk:"system_id"`
	Metric          types.String  `tfsdk:"metric"`
	Condition       types.String  `tfsdk:"condition"`
	Threshold       types.Float64 `tfsdk:"threshold"`
	DurationSeconds types.Int64   `tfsdk:"duration_seconds"`
	Enabled         types.Bool    `tfsdk:"enabled"`
}

// NewSystemAlertResource returns a new resource factory function.
func NewSystemAlertResource() resource.Resource {
	return &SystemAlertResource{}
}

func (r *SystemAlertResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_alert"
}

func (r *SystemAlertResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat system alert for host monitoring.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the system alert.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"system_id": schema.StringAttribute{
				Description: "The ID of the monitored system/host.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metric": schema.StringAttribute{
				Description: "The metric to alert on (cpu, memory, disk, network, load).",
				Required:    true,
			},
			"condition": schema.StringAttribute{
				Description: "The comparison condition (gt, lt, gte, lte, eq).",
				Required:    true,
			},
			"threshold": schema.Float64Attribute{
				Description: "The threshold value for the alert.",
				Required:    true,
			},
			"duration_seconds": schema.Int64Attribute{
				Description: "Duration in seconds the condition must be true before alerting.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the alert is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *SystemAlertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SystemAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SystemAlertResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateSystemAlertRequest{
		Metric:          plan.Metric.ValueString(),
		Condition:       plan.Condition.ValueString(),
		Threshold:       plan.Threshold.ValueFloat64(),
		DurationSeconds: plan.DurationSeconds.ValueInt64(),
		Enabled:         plan.Enabled.ValueBool(),
	}

	alert, err := r.client.CreateSystemAlert(plan.SystemID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating system alert", err.Error())
		return
	}

	plan.ID = types.StringValue(alert.ID)
	plan.Metric = types.StringValue(alert.Metric)
	plan.Condition = types.StringValue(alert.Condition)
	plan.Threshold = types.Float64Value(alert.Threshold)
	plan.DurationSeconds = types.Int64Value(alert.DurationSeconds)
	plan.Enabled = types.BoolValue(alert.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SystemAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SystemAlertResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alert, err := r.client.GetSystemAlert(state.SystemID.ValueString(), state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading system alert", err.Error())
		return
	}

	state.Metric = types.StringValue(alert.Metric)
	state.Condition = types.StringValue(alert.Condition)
	state.Threshold = types.Float64Value(alert.Threshold)
	state.DurationSeconds = types.Int64Value(alert.DurationSeconds)
	state.Enabled = types.BoolValue(alert.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SystemAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SystemAlertResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SystemAlertResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateSystemAlertRequest{
		Metric:          plan.Metric.ValueString(),
		Condition:       plan.Condition.ValueString(),
		Threshold:       plan.Threshold.ValueFloat64(),
		DurationSeconds: plan.DurationSeconds.ValueInt64(),
		Enabled:         plan.Enabled.ValueBool(),
	}

	alert, err := r.client.UpdateSystemAlert(state.SystemID.ValueString(), state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating system alert", err.Error())
		return
	}

	plan.ID = state.ID
	plan.SystemID = state.SystemID
	plan.Metric = types.StringValue(alert.Metric)
	plan.Condition = types.StringValue(alert.Condition)
	plan.Threshold = types.Float64Value(alert.Threshold)
	plan.DurationSeconds = types.Int64Value(alert.DurationSeconds)
	plan.Enabled = types.BoolValue(alert.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SystemAlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SystemAlertResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSystemAlert(state.SystemID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting system alert", err.Error())
		return
	}
}

func (r *SystemAlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: system_id/alert_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: system_id/alert_id",
		)
		return
	}

	systemID := parts[0]
	alertID := parts[1]

	alert, err := r.client.GetSystemAlert(systemID, alertID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing system alert", err.Error())
		return
	}

	state := SystemAlertResourceModel{
		ID:              types.StringValue(alert.ID),
		SystemID:        types.StringValue(systemID),
		Metric:          types.StringValue(alert.Metric),
		Condition:       types.StringValue(alert.Condition),
		Threshold:       types.Float64Value(alert.Threshold),
		DurationSeconds: types.Int64Value(alert.DurationSeconds),
		Enabled:         types.BoolValue(alert.Enabled),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
