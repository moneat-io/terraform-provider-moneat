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

var _ resource.Resource = &OnCallOverrideResource{}
var _ resource.ResourceWithImportState = &OnCallOverrideResource{}

type OnCallOverrideResource struct {
	client *apiclient.Client
}

type OnCallOverrideResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ScheduleID types.String `tfsdk:"schedule_id"`
	UserID     types.String `tfsdk:"user_id"`
	StartTime  types.String `tfsdk:"start_time"`
	EndTime    types.String `tfsdk:"end_time"`
	Reason     types.String `tfsdk:"reason"`
}

func NewOnCallOverrideResource() resource.Resource {
	return &OnCallOverrideResource{}
}

func (r *OnCallOverrideResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_on_call_override"
}

func (r *OnCallOverrideResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat on-call schedule override.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the override.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_id": schema.StringAttribute{
				Description: "The ID of the on-call schedule.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The ID of the user for the override.",
				Required:    true,
			},
			"start_time": schema.StringAttribute{
				Description: "The start time of the override (ISO 8601).",
				Required:    true,
			},
			"end_time": schema.StringAttribute{
				Description: "The end time of the override (ISO 8601).",
				Required:    true,
			},
			"reason": schema.StringAttribute{
				Description: "The reason for the override.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *OnCallOverrideResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OnCallOverrideResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OnCallOverrideResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateOnCallOverrideRequest{
		UserID:    plan.UserID.ValueString(),
		StartTime: plan.StartTime.ValueString(),
		EndTime:   plan.EndTime.ValueString(),
	}
	if !plan.Reason.IsNull() && !plan.Reason.IsUnknown() {
		apiReq.Reason = plan.Reason.ValueString()
	}

	override, err := r.client.CreateOnCallOverride(plan.ScheduleID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating on-call override", err.Error())
		return
	}

	plan.ID = types.StringValue(override.ID)
	plan.UserID = types.StringValue(override.UserID)
	plan.StartTime = types.StringValue(override.StartTime)
	plan.EndTime = types.StringValue(override.EndTime)
	plan.Reason = types.StringValue(override.Reason)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OnCallOverrideResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OnCallOverrideResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	override, err := r.client.GetOnCallOverride(state.ScheduleID.ValueString(), state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading on-call override", err.Error())
		return
	}

	state.UserID = types.StringValue(override.UserID)
	state.StartTime = types.StringValue(override.StartTime)
	state.EndTime = types.StringValue(override.EndTime)
	state.Reason = types.StringValue(override.Reason)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OnCallOverrideResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OnCallOverrideResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OnCallOverrideResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateOnCallOverrideRequest{
		UserID:    plan.UserID.ValueString(),
		StartTime: plan.StartTime.ValueString(),
		EndTime:   plan.EndTime.ValueString(),
	}
	if !plan.Reason.IsNull() && !plan.Reason.IsUnknown() {
		apiReq.Reason = plan.Reason.ValueString()
	}

	override, err := r.client.UpdateOnCallOverride(state.ScheduleID.ValueString(), state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating on-call override", err.Error())
		return
	}

	plan.ID = state.ID
	plan.ScheduleID = state.ScheduleID
	plan.UserID = types.StringValue(override.UserID)
	plan.StartTime = types.StringValue(override.StartTime)
	plan.EndTime = types.StringValue(override.EndTime)
	plan.Reason = types.StringValue(override.Reason)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OnCallOverrideResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OnCallOverrideResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOnCallOverride(state.ScheduleID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting on-call override", err.Error())
		return
	}
}

func (r *OnCallOverrideResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: schedule_id/override_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: schedule_id/override_id",
		)
		return
	}

	scheduleID := parts[0]
	overrideID := parts[1]

	override, err := r.client.GetOnCallOverride(scheduleID, overrideID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing on-call override", err.Error())
		return
	}

	state := OnCallOverrideResourceModel{
		ID:         types.StringValue(override.ID),
		ScheduleID: types.StringValue(scheduleID),
		UserID:     types.StringValue(override.UserID),
		StartTime:  types.StringValue(override.StartTime),
		EndTime:    types.StringValue(override.EndTime),
		Reason:     types.StringValue(override.Reason),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
