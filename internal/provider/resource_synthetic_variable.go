package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &SyntheticVariableResource{}
var _ resource.ResourceWithImportState = &SyntheticVariableResource{}

type SyntheticVariableResource struct {
	client *apiclient.Client
}

type SyntheticVariableResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Value    types.String `tfsdk:"value"`
	IsSecret types.Bool   `tfsdk:"is_secret"`
}

func NewSyntheticVariableResource() resource.Resource {
	return &SyntheticVariableResource{}
}

func (r *SyntheticVariableResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_synthetic_variable"
}

func (r *SyntheticVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat synthetic test variable.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the synthetic variable.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The variable name.",
				Required:    true,
			},
			"value": schema.StringAttribute{
				Description: "The variable value.",
				Required:    true,
				Sensitive:   true,
			},
			"is_secret": schema.BoolAttribute{
				Description: "Whether the API should mask this variable value in responses.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *SyntheticVariableResource) Configure(
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

func (r *SyntheticVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SyntheticVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.SyntheticVariableRequest{
		Name:     plan.Name.ValueString(),
		Value:    plan.Value.ValueString(),
		IsSecret: plan.IsSecret.ValueBool(),
	}
	variable, err := r.client.CreateSyntheticVariable(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating synthetic variable", err.Error())
		return
	}

	mapSyntheticVariableToState(&plan, variable, plan.Value)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SyntheticVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SyntheticVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "synthetic variable ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading synthetic variable", err.Error())
		return
	}
	variable, err := r.client.GetSyntheticVariable(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading synthetic variable", err.Error())
		return
	}

	mapSyntheticVariableToState(&state, variable, state.Value)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SyntheticVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SyntheticVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(plan.ID, "synthetic variable ID")
	if err != nil {
		resp.Diagnostics.AddError("Error updating synthetic variable", err.Error())
		return
	}
	apiReq := apiclient.SyntheticVariableRequest{
		Name:     plan.Name.ValueString(),
		Value:    plan.Value.ValueString(),
		IsSecret: plan.IsSecret.ValueBool(),
	}
	variable, err := r.client.UpdateSyntheticVariable(id, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating synthetic variable", err.Error())
		return
	}

	mapSyntheticVariableToState(&plan, variable, plan.Value)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SyntheticVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SyntheticVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "synthetic variable ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting synthetic variable", err.Error())
		return
	}
	err = r.client.DeleteSyntheticVariable(id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting synthetic variable", err.Error())
		return
	}
}

func (r *SyntheticVariableResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := parseTerraformID(types.StringValue(req.ID), "synthetic variable ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing synthetic variable", err.Error())
		return
	}
	variable, err := r.client.GetSyntheticVariable(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing synthetic variable", err.Error())
		return
	}

	state := SyntheticVariableResourceModel{}
	mapSyntheticVariableToState(&state, variable, types.StringValue(variable.Value))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func mapSyntheticVariableToState(
	model *SyntheticVariableResourceModel,
	variable *apiclient.SyntheticVariable,
	existingValue types.String,
) {
	model.ID = terraformID(variable.ID)
	model.Name = types.StringValue(variable.Name)
	model.IsSecret = types.BoolValue(variable.IsSecret)
	if variable.IsSecret && variable.Value == "********" {
		model.Value = existingValue
	} else {
		model.Value = types.StringValue(variable.Value)
	}
}
