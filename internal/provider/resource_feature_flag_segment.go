package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &FeatureFlagSegmentResource{}
var _ resource.ResourceWithImportState = &FeatureFlagSegmentResource{}

type FeatureFlagSegmentResource struct {
	client *apiclient.Client
}

type FeatureFlagSegmentResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Key            types.String `tfsdk:"key"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	ConditionsJSON types.String `tfsdk:"conditions_json"`
}

func NewFeatureFlagSegmentResource() resource.Resource {
	return &FeatureFlagSegmentResource{}
}

func (r *FeatureFlagSegmentResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_feature_flag_segment"
}

func (r *FeatureFlagSegmentResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat feature flag segment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the feature flag segment.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The stable segment key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the segment.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Optional segment description.",
				Optional:    true,
			},
			"conditions_json": schema.StringAttribute{
				Description: "Segment conditions as JSON.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(`{"all":[]}`),
			},
		},
	}
}

func (r *FeatureFlagSegmentResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
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

func (r *FeatureFlagSegmentResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan FeatureFlagSegmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	segment, normalized, err := r.upsert(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating feature flag segment", err.Error())
		return
	}

	mapFeatureFlagSegmentToState(&plan, segment, normalized)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagSegmentResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state FeatureFlagSegmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	segment, err := r.client.GetFeatureFlagSegment(state.Key.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading feature flag segment", err.Error())
		return
	}

	mapFeatureFlagSegmentToState(&state, segment, rawMessageString(segment.Conditions))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FeatureFlagSegmentResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan FeatureFlagSegmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	segment, normalized, err := r.upsert(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error updating feature flag segment", err.Error())
		return
	}

	mapFeatureFlagSegmentToState(&plan, segment, normalized)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagSegmentResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state FeatureFlagSegmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFeatureFlagSegment(state.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting feature flag segment", err.Error())
		return
	}
}

func (r *FeatureFlagSegmentResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	segment, err := r.client.GetFeatureFlagSegment(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing feature flag segment", err.Error())
		return
	}

	state := FeatureFlagSegmentResourceModel{}
	mapFeatureFlagSegmentToState(&state, segment, rawMessageString(segment.Conditions))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FeatureFlagSegmentResource) upsert(
	model FeatureFlagSegmentResourceModel,
) (*apiclient.FeatureFlagSegment, string, error) {
	conditions, normalized, err := rawMessageFromJSONString(model.ConditionsJSON.ValueString())
	if err != nil {
		return nil, "", err
	}

	apiReq := apiclient.FeatureFlagSegmentRequest{
		Key:        model.Key.ValueString(),
		Name:       model.Name.ValueString(),
		Conditions: conditions,
	}
	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		apiReq.Description = model.Description.ValueString()
	}

	segment, err := r.client.UpsertFeatureFlagSegment(apiReq)
	if err != nil {
		return nil, "", err
	}
	return segment, normalized, nil
}

func mapFeatureFlagSegmentToState(
	model *FeatureFlagSegmentResourceModel,
	segment *apiclient.FeatureFlagSegment,
	conditionsJSON string,
) {
	model.ID = terraformID(segment.ID)
	model.Key = types.StringValue(segment.Key)
	model.Name = types.StringValue(segment.Name)
	if segment.Description == "" {
		model.Description = types.StringNull()
	} else {
		model.Description = types.StringValue(segment.Description)
	}
	model.ConditionsJSON = types.StringValue(conditionsJSON)
}
