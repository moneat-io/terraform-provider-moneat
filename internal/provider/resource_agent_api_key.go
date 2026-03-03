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

var _ resource.Resource = &AgentAPIKeyResource{}
var _ resource.ResourceWithImportState = &AgentAPIKeyResource{}

type AgentAPIKeyResource struct {
	client *apiclient.Client
}

type AgentAPIKeyResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Key  types.String `tfsdk:"key"`
}

func NewAgentAPIKeyResource() resource.Resource {
	return &AgentAPIKeyResource{}
}

func (r *AgentAPIKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_api_key"
}

func (r *AgentAPIKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat agent API key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the agent API key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the agent API key.",
				Required:    true,
			},
			"key": schema.StringAttribute{
				Description: "The generated API key value. Only available after creation.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *AgentAPIKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AgentAPIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgentAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateAgentAPIKeyRequest{
		Name: plan.Name.ValueString(),
	}

	key, err := r.client.CreateAgentAPIKey(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating agent API key", err.Error())
		return
	}

	plan.ID = types.StringValue(key.ID)
	plan.Name = types.StringValue(key.Name)
	plan.Key = types.StringValue(key.Key)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AgentAPIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgentAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GetAgentAPIKey(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading agent API key", err.Error())
		return
	}

	state.Name = types.StringValue(key.Name)
	if key.Key != "" {
		state.Key = types.StringValue(key.Key)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AgentAPIKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AgentAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateAgentAPIKeyRequest{
		Name: plan.Name.ValueString(),
	}

	key, err := r.client.UpdateAgentAPIKey(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating agent API key", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Key = state.Key // Preserve key from state
	plan.Name = types.StringValue(key.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AgentAPIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgentAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAgentAPIKey(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting agent API key", err.Error())
		return
	}
}

func (r *AgentAPIKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	key, err := r.client.GetAgentAPIKey(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing agent API key", err.Error())
		return
	}

	state := AgentAPIKeyResourceModel{
		ID:   types.StringValue(key.ID),
		Name: types.StringValue(key.Name),
		Key:  types.StringValue(key.Key),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
