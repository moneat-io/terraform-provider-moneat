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

var _ resource.Resource = &SilencePeriodResource{}
var _ resource.ResourceWithImportState = &SilencePeriodResource{}

type SilencePeriodResource struct {
	client *apiclient.Client
}

type SilencePeriodResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	StartTime types.String `tfsdk:"start_time"`
	EndTime   types.String `tfsdk:"end_time"`
	Reason    types.String `tfsdk:"reason"`
}

func NewSilencePeriodResource() resource.Resource {
	return &SilencePeriodResource{}
}

func (r *SilencePeriodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_silence_period"
}

func (r *SilencePeriodResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat alert silence period.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the silence period.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the silence period.",
				Required:    true,
			},
			"start_time": schema.StringAttribute{
				Description: "The start time of the silence period (ISO 8601).",
				Required:    true,
			},
			"end_time": schema.StringAttribute{
				Description: "The end time of the silence period (ISO 8601).",
				Required:    true,
			},
			"reason": schema.StringAttribute{
				Description: "The reason for the silence period.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SilencePeriodResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SilencePeriodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SilencePeriodResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateSilencePeriodRequest{
		Name:      plan.Name.ValueString(),
		StartTime: plan.StartTime.ValueString(),
		EndTime:   plan.EndTime.ValueString(),
	}
	if !plan.Reason.IsNull() && !plan.Reason.IsUnknown() {
		apiReq.Reason = plan.Reason.ValueString()
	}

	period, err := r.client.CreateSilencePeriod(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating silence period", err.Error())
		return
	}

	plan.ID = types.StringValue(period.ID)
	plan.Name = types.StringValue(period.Name)
	plan.StartTime = types.StringValue(period.StartTime)
	plan.EndTime = types.StringValue(period.EndTime)
	plan.Reason = types.StringValue(period.Reason)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SilencePeriodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SilencePeriodResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	period, err := r.client.GetSilencePeriod(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading silence period", err.Error())
		return
	}

	state.Name = types.StringValue(period.Name)
	state.StartTime = types.StringValue(period.StartTime)
	state.EndTime = types.StringValue(period.EndTime)
	state.Reason = types.StringValue(period.Reason)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SilencePeriodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SilencePeriodResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SilencePeriodResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateSilencePeriodRequest{
		Name:      plan.Name.ValueString(),
		StartTime: plan.StartTime.ValueString(),
		EndTime:   plan.EndTime.ValueString(),
	}
	if !plan.Reason.IsNull() && !plan.Reason.IsUnknown() {
		apiReq.Reason = plan.Reason.ValueString()
	}

	period, err := r.client.UpdateSilencePeriod(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating silence period", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(period.Name)
	plan.StartTime = types.StringValue(period.StartTime)
	plan.EndTime = types.StringValue(period.EndTime)
	plan.Reason = types.StringValue(period.Reason)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SilencePeriodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SilencePeriodResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSilencePeriod(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting silence period", err.Error())
		return
	}
}

func (r *SilencePeriodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	period, err := r.client.GetSilencePeriod(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing silence period", err.Error())
		return
	}

	state := SilencePeriodResourceModel{
		ID:        types.StringValue(period.ID),
		Name:      types.StringValue(period.Name),
		StartTime: types.StringValue(period.StartTime),
		EndTime:   types.StringValue(period.EndTime),
		Reason:    types.StringValue(period.Reason),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
