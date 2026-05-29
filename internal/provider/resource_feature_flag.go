package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &FeatureFlagResource{}
var _ resource.ResourceWithImportState = &FeatureFlagResource{}

type FeatureFlagResource struct {
	client *apiclient.Client
}

type FeatureFlagResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Key           types.String `tfsdk:"key"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ValueType     types.String `tfsdk:"value_type"`
	ClientVisible types.Bool   `tfsdk:"client_visible"`
	Tags          types.List   `tfsdk:"tags"`
	VariantsJSON  types.String `tfsdk:"variants_json"`
}

func NewFeatureFlagResource() resource.Resource {
	return &FeatureFlagResource{}
}

func (r *FeatureFlagResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_feature_flag"
}

func (r *FeatureFlagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat feature flag.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the feature flag.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The stable feature flag key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the feature flag.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Optional feature flag description.",
				Optional:    true,
			},
			"value_type": schema.StringAttribute{
				Description: "The value type: BOOLEAN, STRING, INTEGER, DOUBLE, or OBJECT.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_visible": schema.BoolAttribute{
				Description: "Whether client SDK keys can evaluate this flag.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"tags": schema.ListAttribute{
				Description: "Tags for organizing the feature flag.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default: listdefault.StaticValue(
					types.ListValueMust(types.StringType, []attr.Value{}),
				),
			},
			"variants_json": schema.StringAttribute{
				Description: "Feature flag variants as JSON. Use jsonencode with key, optional name, and value.",
				Required:    true,
			},
		},
	}
}

func (r *FeatureFlagResource) Configure(
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

func (r *FeatureFlagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FeatureFlagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags, variants, normalized, ok := featureFlagPlanPayload(ctx, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	apiReq := apiclient.CreateFeatureFlagRequest{
		Key:           plan.Key.ValueString(),
		Name:          plan.Name.ValueString(),
		ValueType:     plan.ValueType.ValueString(),
		ClientVisible: plan.ClientVisible.ValueBool(),
		Tags:          tags,
		Variants:      variants,
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.Description = plan.Description.ValueString()
	}

	flag, err := r.client.CreateFeatureFlag(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating feature flag", err.Error())
		return
	}

	mapFeatureFlagToState(ctx, &resp.Diagnostics, &plan, flag, normalized)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FeatureFlagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	flag, err := r.client.GetFeatureFlag(state.Key.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading feature flag", err.Error())
		return
	}

	mapFeatureFlagToState(ctx, &resp.Diagnostics, &state, flag, "")
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FeatureFlagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FeatureFlagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags, variants, normalized, ok := featureFlagPlanPayload(ctx, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	clientVisible := plan.ClientVisible.ValueBool()
	apiReq := apiclient.UpdateFeatureFlagRequest{
		Name:          plan.Name.ValueString(),
		ClientVisible: &clientVisible,
		Tags:          tags,
		Variants:      variants,
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.Description = plan.Description.ValueString()
	}

	flag, err := r.client.UpdateFeatureFlag(plan.Key.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating feature flag", err.Error())
		return
	}

	mapFeatureFlagToState(ctx, &resp.Diagnostics, &plan, flag, normalized)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FeatureFlagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFeatureFlag(state.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting feature flag", err.Error())
		return
	}
}

func (r *FeatureFlagResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	flag, err := r.client.GetFeatureFlag(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing feature flag", err.Error())
		return
	}

	state := FeatureFlagResourceModel{}
	mapFeatureFlagToState(ctx, &resp.Diagnostics, &state, flag, "")
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func featureFlagPlanPayload(
	ctx context.Context,
	model FeatureFlagResourceModel,
	diags *diag.Diagnostics,
) ([]string, []apiclient.FeatureFlagVariantRequest, string, bool) {
	var tags []string
	diags.Append(model.Tags.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil, nil, "", false
	}

	rawVariants, normalized, err := rawMessageFromJSONString(model.VariantsJSON.ValueString())
	if err != nil {
		diags.AddError("Invalid feature flag variants JSON", err.Error())
		return nil, nil, "", false
	}
	var variants []apiclient.FeatureFlagVariantRequest
	if err := json.Unmarshal(rawVariants, &variants); err != nil {
		diags.AddError("Invalid feature flag variants JSON", err.Error())
		return nil, nil, "", false
	}

	return tags, variants, normalized, true
}

func mapFeatureFlagToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *FeatureFlagResourceModel,
	flag *apiclient.FeatureFlag,
	variantsJSON string,
) {
	model.ID = terraformID(flag.ID)
	model.Key = types.StringValue(flag.Key)
	model.Name = types.StringValue(flag.Name)
	if flag.Description == "" {
		model.Description = types.StringNull()
	} else {
		model.Description = types.StringValue(flag.Description)
	}
	model.ValueType = types.StringValue(flag.ValueType)
	model.ClientVisible = types.BoolValue(flag.ClientVisible)

	tags, tagDiags := types.ListValueFrom(ctx, types.StringType, flag.Tags)
	diags.Append(tagDiags...)
	if diags.HasError() {
		return
	}
	model.Tags = tags

	if variantsJSON == "" {
		variantsJSON = featureFlagVariantsJSON(flag.Variants)
	}
	model.VariantsJSON = types.StringValue(variantsJSON)
}

func featureFlagVariantsJSON(variants []apiclient.FeatureFlagVariant) string {
	requestVariants := make([]apiclient.FeatureFlagVariantRequest, len(variants))
	for i, variant := range variants {
		requestVariants[i] = apiclient.FeatureFlagVariantRequest{
			Key:   variant.Key,
			Name:  variant.Name,
			Value: variant.Value,
		}
	}
	bytes, err := json.Marshal(requestVariants)
	if err != nil {
		return "[]"
	}
	return string(bytes)
}
