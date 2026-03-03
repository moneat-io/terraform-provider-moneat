package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &OnCallScheduleResource{}
var _ resource.ResourceWithImportState = &OnCallScheduleResource{}

type OnCallScheduleResource struct {
	client *apiclient.Client
}

type OnCallScheduleResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Timezone     types.String `tfsdk:"timezone"`
	RotationType types.String `tfsdk:"rotation_type"`
}

func NewOnCallScheduleResource() resource.Resource {
	return &OnCallScheduleResource{}
}

func (r *OnCallScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_on_call_schedule"
}

func (r *OnCallScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat on-call schedule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the on-call schedule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the on-call schedule.",
				Required:    true,
			},
			"timezone": schema.StringAttribute{
				Description: "The timezone for the schedule (e.g., America/New_York).",
				Required:    true,
			},
			"rotation_type": schema.StringAttribute{
				Description: "The rotation type (daily, weekly, custom).",
				Required:    true,
			},
		},
	}
}

func (r *OnCallScheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OnCallScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OnCallScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateOnCallScheduleRequest{
		Name:         plan.Name.ValueString(),
		Timezone:     plan.Timezone.ValueString(),
		RotationType: plan.RotationType.ValueString(),
	}

	schedule, err := r.client.CreateOnCallSchedule(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating on-call schedule", err.Error())
		return
	}

	plan.ID = types.StringValue(schedule.ID)
	plan.Name = types.StringValue(schedule.Name)
	plan.Timezone = types.StringValue(schedule.Timezone)
	plan.RotationType = types.StringValue(schedule.RotationType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OnCallScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OnCallScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schedule, err := r.client.GetOnCallSchedule(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading on-call schedule", err.Error())
		return
	}

	state.Name = types.StringValue(schedule.Name)
	state.Timezone = types.StringValue(schedule.Timezone)
	state.RotationType = types.StringValue(schedule.RotationType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OnCallScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OnCallScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OnCallScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateOnCallScheduleRequest{
		Name:         plan.Name.ValueString(),
		Timezone:     plan.Timezone.ValueString(),
		RotationType: plan.RotationType.ValueString(),
	}

	schedule, err := r.client.UpdateOnCallSchedule(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating on-call schedule", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(schedule.Name)
	plan.Timezone = types.StringValue(schedule.Timezone)
	plan.RotationType = types.StringValue(schedule.RotationType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OnCallScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OnCallScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOnCallSchedule(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting on-call schedule", err.Error())
		return
	}
}

func (r *OnCallScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	schedule, err := r.client.GetOnCallSchedule(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing on-call schedule", err.Error())
		return
	}

	state := OnCallScheduleResourceModel{
		ID:           types.StringValue(schedule.ID),
		Name:         types.StringValue(schedule.Name),
		Timezone:     types.StringValue(schedule.Timezone),
		RotationType: types.StringValue(schedule.RotationType),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
