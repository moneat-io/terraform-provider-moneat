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

var _ resource.Resource = &EscalationPolicyResource{}
var _ resource.ResourceWithImportState = &EscalationPolicyResource{}

type EscalationPolicyResource struct {
	client *apiclient.Client
}

type EscalationPolicyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func NewEscalationPolicyResource() resource.Resource {
	return &EscalationPolicyResource{}
}

func (r *EscalationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_escalation_policy"
}

func (r *EscalationPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat escalation policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the escalation policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the escalation policy.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the escalation policy.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *EscalationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EscalationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EscalationPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateEscalationPolicyRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.Description = plan.Description.ValueString()
	}

	policy, err := r.client.CreateEscalationPolicy(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating escalation policy", err.Error())
		return
	}

	plan.ID = types.StringValue(policy.ID)
	plan.Name = types.StringValue(policy.Name)
	plan.Description = types.StringValue(policy.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EscalationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EscalationPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.client.GetEscalationPolicy(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading escalation policy", err.Error())
		return
	}

	state.Name = types.StringValue(policy.Name)
	state.Description = types.StringValue(policy.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EscalationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EscalationPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state EscalationPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateEscalationPolicyRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.Description = plan.Description.ValueString()
	}

	policy, err := r.client.UpdateEscalationPolicy(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating escalation policy", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(policy.Name)
	plan.Description = types.StringValue(policy.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EscalationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EscalationPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEscalationPolicy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting escalation policy", err.Error())
		return
	}
}

func (r *EscalationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	policy, err := r.client.GetEscalationPolicy(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing escalation policy", err.Error())
		return
	}

	state := EscalationPolicyResourceModel{
		ID:          types.StringValue(policy.ID),
		Name:        types.StringValue(policy.Name),
		Description: types.StringValue(policy.Description),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
