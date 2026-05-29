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

var _ resource.Resource = &FeatureFlagEnvironmentResource{}
var _ resource.ResourceWithImportState = &FeatureFlagEnvironmentResource{}

type FeatureFlagEnvironmentResource struct {
	client *apiclient.Client
}

type FeatureFlagEnvironmentResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Version     types.Int64  `tfsdk:"version"`
}

func NewFeatureFlagEnvironmentResource() resource.Resource {
	return &FeatureFlagEnvironmentResource{}
}

func (r *FeatureFlagEnvironmentResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_feature_flag_environment"
}

func (r *FeatureFlagEnvironmentResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat feature flag environment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the feature flag environment.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The stable environment key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the environment.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Optional environment description.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.Int64Attribute{
				Description: "The environment configuration version.",
				Computed:    true,
			},
		},
	}
}

func (r *FeatureFlagEnvironmentResource) Configure(
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

func (r *FeatureFlagEnvironmentResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan FeatureFlagEnvironmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateFeatureFlagEnvironmentRequest{
		Key:  plan.Key.ValueString(),
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.Description = plan.Description.ValueString()
	}

	environment, err := r.client.CreateFeatureFlagEnvironment(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating feature flag environment", err.Error())
		return
	}

	mapFeatureFlagEnvironmentToState(&plan, environment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagEnvironmentResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state FeatureFlagEnvironmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environment, err := r.client.GetFeatureFlagEnvironment(state.Key.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading feature flag environment", err.Error())
		return
	}

	mapFeatureFlagEnvironmentToState(&state, environment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FeatureFlagEnvironmentResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan FeatureFlagEnvironmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureFlagEnvironmentResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
	// The API does not expose environment deletion.
}

func (r *FeatureFlagEnvironmentResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	environment, err := r.client.GetFeatureFlagEnvironment(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing feature flag environment", err.Error())
		return
	}

	state := FeatureFlagEnvironmentResourceModel{}
	mapFeatureFlagEnvironmentToState(&state, environment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func mapFeatureFlagEnvironmentToState(
	model *FeatureFlagEnvironmentResourceModel,
	environment *apiclient.FeatureFlagEnvironment,
) {
	model.ID = terraformID(environment.ID)
	model.Key = types.StringValue(environment.Key)
	model.Name = types.StringValue(environment.Name)
	if environment.Description == "" {
		model.Description = types.StringNull()
	} else {
		model.Description = types.StringValue(environment.Description)
	}
	model.Version = types.Int64Value(int64(environment.Version))
}
