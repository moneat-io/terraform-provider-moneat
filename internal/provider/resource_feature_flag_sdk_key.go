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

var _ resource.Resource = &FeatureFlagSdkKeyResource{}
var _ resource.ResourceWithImportState = &FeatureFlagSdkKeyResource{}

type FeatureFlagSdkKeyResource struct {
	client *apiclient.Client
}

type FeatureFlagSdkKeyResourceModel struct {
	ID             types.String `tfsdk:"id"`
	EnvironmentKey types.String `tfsdk:"environment_key"`
	Name           types.String `tfsdk:"name"`
	KeyType        types.String `tfsdk:"key_type"`
	KeyPrefix      types.String `tfsdk:"key_prefix"`
	Key            types.String `tfsdk:"key"`
}

func NewFeatureFlagSdkKeyResource() resource.Resource {
	return &FeatureFlagSdkKeyResource{}
}

func (r *FeatureFlagSdkKeyResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_feature_flag_sdk_key"
}

func (r *FeatureFlagSdkKeyResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat feature flag SDK key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the SDK key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_key": schema.StringAttribute{
				Description: "The feature flag environment key this SDK key can evaluate.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the SDK key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key_type": schema.StringAttribute{
				Description: "The SDK key type, such as server or client.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key_prefix": schema.StringAttribute{
				Description: "The display prefix for the generated SDK key.",
				Computed:    true,
			},
			"key": schema.StringAttribute{
				Description: "The generated SDK key value. Only available after creation.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *FeatureFlagSdkKeyResource) Configure(
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

func (r *FeatureFlagSdkKeyResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan FeatureFlagSdkKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.FeatureFlagSdkKeyRequest{
		EnvironmentKey: plan.EnvironmentKey.ValueString(),
		Name:           plan.Name.ValueString(),
		KeyType:        plan.KeyType.ValueString(),
	}
	key, err := r.client.CreateFeatureFlagSdkKey(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating feature flag SDK key", err.Error())
		return
	}

	mapFeatureFlagSdkKeyToState(&plan, key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagSdkKeyResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state FeatureFlagSdkKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "feature flag SDK key ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading feature flag SDK key", err.Error())
		return
	}
	key, err := r.client.GetFeatureFlagSdkKey(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading feature flag SDK key", err.Error())
		return
	}

	existingKey := state.Key
	mapFeatureFlagSdkKeyToState(&state, key)
	if key.Key == "" {
		state.Key = existingKey
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FeatureFlagSdkKeyResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan FeatureFlagSdkKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagSdkKeyResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state FeatureFlagSdkKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "feature flag SDK key ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting feature flag SDK key", err.Error())
		return
	}
	err = r.client.DeleteFeatureFlagSdkKey(id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting feature flag SDK key", err.Error())
		return
	}
}

func (r *FeatureFlagSdkKeyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := parseTerraformID(types.StringValue(req.ID), "feature flag SDK key ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing feature flag SDK key", err.Error())
		return
	}
	key, err := r.client.GetFeatureFlagSdkKey(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing feature flag SDK key", err.Error())
		return
	}

	state := FeatureFlagSdkKeyResourceModel{}
	mapFeatureFlagSdkKeyToState(&state, key)
	if key.Key == "" {
		state.Key = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func mapFeatureFlagSdkKeyToState(model *FeatureFlagSdkKeyResourceModel, key *apiclient.FeatureFlagSdkKey) {
	model.ID = terraformID(key.ID)
	model.EnvironmentKey = types.StringValue(key.EnvironmentKey)
	model.Name = types.StringValue(key.Name)
	model.KeyType = types.StringValue(key.KeyType)
	model.KeyPrefix = types.StringValue(key.KeyPrefix)
	if key.Key != "" {
		model.Key = types.StringValue(key.Key)
	}
}
