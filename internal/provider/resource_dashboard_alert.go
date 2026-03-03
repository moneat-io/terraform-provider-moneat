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

var _ resource.Resource = &DashboardAlertResource{}
var _ resource.ResourceWithImportState = &DashboardAlertResource{}

type DashboardAlertResource struct {
	client *apiclient.Client
}

type DashboardAlertResourceModel struct {
	ID               types.String  `tfsdk:"id"`
	DashboardID      types.String  `tfsdk:"dashboard_id"`
	WidgetID         types.String  `tfsdk:"widget_id"`
	Name             types.String  `tfsdk:"name"`
	Condition        types.String  `tfsdk:"condition"`
	Threshold        types.Float64 `tfsdk:"threshold"`
	DurationSeconds  types.Int64   `tfsdk:"duration_seconds"`
	IncidentSeverity types.String  `tfsdk:"incident_severity"`
}

func NewDashboardAlertResource() resource.Resource {
	return &DashboardAlertResource{}
}

func (r *DashboardAlertResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard_alert"
}

func (r *DashboardAlertResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat dashboard alert.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the dashboard alert.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dashboard_id": schema.StringAttribute{
				Description: "The ID of the dashboard.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"widget_id": schema.StringAttribute{
				Description: "The ID of the widget to alert on.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the alert.",
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
			"incident_severity": schema.StringAttribute{
				Description: "The severity of incidents created by this alert.",
				Required:    true,
			},
		},
	}
}

func (r *DashboardAlertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DashboardAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DashboardAlertResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateDashboardAlertRequest{
		WidgetID:         plan.WidgetID.ValueString(),
		Name:             plan.Name.ValueString(),
		Condition:        plan.Condition.ValueString(),
		Threshold:        plan.Threshold.ValueFloat64(),
		DurationSeconds:  plan.DurationSeconds.ValueInt64(),
		IncidentSeverity: plan.IncidentSeverity.ValueString(),
	}

	alert, err := r.client.CreateDashboardAlert(plan.DashboardID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating dashboard alert", err.Error())
		return
	}

	plan.ID = types.StringValue(alert.ID)
	plan.WidgetID = types.StringValue(alert.WidgetID)
	plan.Name = types.StringValue(alert.Name)
	plan.Condition = types.StringValue(alert.Condition)
	plan.Threshold = types.Float64Value(alert.Threshold)
	plan.DurationSeconds = types.Int64Value(alert.DurationSeconds)
	plan.IncidentSeverity = types.StringValue(alert.IncidentSeverity)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DashboardAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DashboardAlertResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alert, err := r.client.GetDashboardAlert(state.DashboardID.ValueString(), state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading dashboard alert", err.Error())
		return
	}

	state.WidgetID = types.StringValue(alert.WidgetID)
	state.Name = types.StringValue(alert.Name)
	state.Condition = types.StringValue(alert.Condition)
	state.Threshold = types.Float64Value(alert.Threshold)
	state.DurationSeconds = types.Int64Value(alert.DurationSeconds)
	state.IncidentSeverity = types.StringValue(alert.IncidentSeverity)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DashboardAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DashboardAlertResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state DashboardAlertResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateDashboardAlertRequest{
		WidgetID:         plan.WidgetID.ValueString(),
		Name:             plan.Name.ValueString(),
		Condition:        plan.Condition.ValueString(),
		Threshold:        plan.Threshold.ValueFloat64(),
		DurationSeconds:  plan.DurationSeconds.ValueInt64(),
		IncidentSeverity: plan.IncidentSeverity.ValueString(),
	}

	alert, err := r.client.UpdateDashboardAlert(state.DashboardID.ValueString(), state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating dashboard alert", err.Error())
		return
	}

	plan.ID = state.ID
	plan.DashboardID = state.DashboardID
	plan.WidgetID = types.StringValue(alert.WidgetID)
	plan.Name = types.StringValue(alert.Name)
	plan.Condition = types.StringValue(alert.Condition)
	plan.Threshold = types.Float64Value(alert.Threshold)
	plan.DurationSeconds = types.Int64Value(alert.DurationSeconds)
	plan.IncidentSeverity = types.StringValue(alert.IncidentSeverity)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DashboardAlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DashboardAlertResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDashboardAlert(state.DashboardID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting dashboard alert", err.Error())
		return
	}
}

func (r *DashboardAlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: dashboard_id/alert_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: dashboard_id/alert_id",
		)
		return
	}

	dashboardID := parts[0]
	alertID := parts[1]

	alert, err := r.client.GetDashboardAlert(dashboardID, alertID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing dashboard alert", err.Error())
		return
	}

	state := DashboardAlertResourceModel{
		ID:               types.StringValue(alert.ID),
		DashboardID:      types.StringValue(dashboardID),
		WidgetID:         types.StringValue(alert.WidgetID),
		Name:             types.StringValue(alert.Name),
		Condition:        types.StringValue(alert.Condition),
		Threshold:        types.Float64Value(alert.Threshold),
		DurationSeconds:  types.Int64Value(alert.DurationSeconds),
		IncidentSeverity: types.StringValue(alert.IncidentSeverity),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
