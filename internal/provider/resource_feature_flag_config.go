package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &FeatureFlagConfigResource{}
var _ resource.ResourceWithImportState = &FeatureFlagConfigResource{}

type FeatureFlagConfigResource struct {
	client *apiclient.Client
}

type FeatureFlagConfigResourceModel struct {
	ID                types.String `tfsdk:"id"`
	FlagKey           types.String `tfsdk:"flag_key"`
	EnvironmentKey    types.String `tfsdk:"environment_key"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	DefaultVariantKey types.String `tfsdk:"default_variant_key"`
	OffVariantKey     types.String `tfsdk:"off_variant_key"`
	RulesJSON         types.String `tfsdk:"rules_json"`
	Version           types.Int64  `tfsdk:"version"`
}

func NewFeatureFlagConfigResource() resource.Resource {
	return &FeatureFlagConfigResource{}
}

func (r *FeatureFlagConfigResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_feature_flag_config"
}

func (r *FeatureFlagConfigResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat feature flag environment configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The feature flag config import identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"flag_key": schema.StringAttribute{
				Description: "The feature flag key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_key": schema.StringAttribute{
				Description: "The feature flag environment key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the flag is enabled in this environment.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"default_variant_key": schema.StringAttribute{
				Description: "The variant served when no rule matches.",
				Required:    true,
			},
			"off_variant_key": schema.StringAttribute{
				Description: "The variant served when the flag is disabled.",
				Required:    true,
			},
			"rules_json": schema.StringAttribute{
				Description: "Targeting rules as JSON.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(`{"rules":[]}`),
			},
			"version": schema.Int64Attribute{
				Description: "The feature flag config version.",
				Computed:    true,
			},
		},
	}
}

func (r *FeatureFlagConfigResource) Configure(
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

func (r *FeatureFlagConfigResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan FeatureFlagConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, normalized, err := r.updateConfig(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating feature flag config", err.Error())
		return
	}

	mapFeatureFlagConfigToState(&plan, config, normalized)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagConfigResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state FeatureFlagConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.GetFeatureFlagConfig(state.FlagKey.ValueString(), state.EnvironmentKey.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading feature flag config", err.Error())
		return
	}

	mapFeatureFlagConfigToState(&state, config, rawMessageString(config.Rules))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FeatureFlagConfigResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan FeatureFlagConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, normalized, err := r.updateConfig(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error updating feature flag config", err.Error())
		return
	}

	mapFeatureFlagConfigToState(&plan, config, normalized)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagConfigResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state FeatureFlagConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enabled := false
	_, err := r.client.UpdateFeatureFlagConfig(
		state.FlagKey.ValueString(),
		state.EnvironmentKey.ValueString(),
		apiclient.UpdateFeatureFlagConfigRequest{Enabled: &enabled},
	)
	if err != nil {
		resp.Diagnostics.AddError("Error disabling feature flag config", err.Error())
		return
	}
}

func (r *FeatureFlagConfigResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: flag_key/environment_key",
		)
		return
	}
	flagKey := parts[0]
	environmentKey := parts[1]

	config, err := r.client.GetFeatureFlagConfig(flagKey, environmentKey)
	if err != nil {
		resp.Diagnostics.AddError("Error importing feature flag config", err.Error())
		return
	}

	state := FeatureFlagConfigResourceModel{
		FlagKey:        types.StringValue(flagKey),
		EnvironmentKey: types.StringValue(environmentKey),
	}
	mapFeatureFlagConfigToState(&state, config, rawMessageString(config.Rules))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FeatureFlagConfigResource) updateConfig(
	model FeatureFlagConfigResourceModel,
) (*apiclient.FeatureFlagConfig, string, error) {
	rules, normalized, err := rawMessageFromJSONString(model.RulesJSON.ValueString())
	if err != nil {
		return nil, "", err
	}
	enabled := model.Enabled.ValueBool()
	apiReq := apiclient.UpdateFeatureFlagConfigRequest{
		Enabled:           &enabled,
		DefaultVariantKey: model.DefaultVariantKey.ValueString(),
		OffVariantKey:     model.OffVariantKey.ValueString(),
		Rules:             rules,
	}
	config, err := r.client.UpdateFeatureFlagConfig(
		model.FlagKey.ValueString(),
		model.EnvironmentKey.ValueString(),
		apiReq,
	)
	if err != nil {
		return nil, "", err
	}
	return config, normalized, nil
}

func mapFeatureFlagConfigToState(
	model *FeatureFlagConfigResourceModel,
	config *apiclient.FeatureFlagConfig,
	rulesJSON string,
) {
	model.ID = types.StringValue(model.FlagKey.ValueString() + "/" + config.EnvironmentKey)
	model.EnvironmentKey = types.StringValue(config.EnvironmentKey)
	model.Enabled = types.BoolValue(config.Enabled)
	model.DefaultVariantKey = types.StringValue(config.DefaultVariantKey)
	model.OffVariantKey = types.StringValue(config.OffVariantKey)
	model.RulesJSON = types.StringValue(rulesJSON)
	model.Version = types.Int64Value(int64(config.Version))
}
