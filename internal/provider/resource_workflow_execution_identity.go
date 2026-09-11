package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &WorkflowExecutionIdentityResource{}
var _ resource.ResourceWithImportState = &WorkflowExecutionIdentityResource{}
var _ resource.ResourceWithValidateConfig = &WorkflowExecutionIdentityResource{}

type WorkflowExecutionIdentityResource struct {
	client *apiclient.Client
}

type WorkflowExecutionIdentityResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	WorkflowID         types.String `tfsdk:"workflow_id"`
	Type               types.String `tfsdk:"type"`
	ServicePrincipalID types.String `tfsdk:"service_principal_id"`
	WorkflowVersion    types.Int64  `tfsdk:"workflow_version"`
}

type workflowExecutionIdentityState struct {
	Type               string `json:"type"`
	ServicePrincipalID string `json:"service_principal_id,omitempty"`
}

func NewWorkflowExecutionIdentityResource() resource.Resource {
	return &WorkflowExecutionIdentityResource{}
}

func (r *WorkflowExecutionIdentityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_execution_identity"
}

func (r *WorkflowExecutionIdentityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the durable execution identity selected for workflow runs.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"workflow_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"type": schema.StringAttribute{
				Description: "initiator, owner, or service_principal.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("initiator"),
			},
			"service_principal_id": schema.StringAttribute{Optional: true, Computed: true},
			"workflow_version":     schema.Int64Attribute{Computed: true},
		},
	}
}

func (r *WorkflowExecutionIdentityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *apiclient.Client, got: %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *WorkflowExecutionIdentityResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config WorkflowExecutionIdentityResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := validateExecutionIdentity(config); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("service_principal_id"), "Invalid execution identity", err.Error())
	}
}

func (r *WorkflowExecutionIdentityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkflowExecutionIdentityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := validateExecutionIdentity(plan); err != nil {
		resp.Diagnostics.AddError("Invalid workflow execution identity", err.Error())
		return
	}
	workflowID, err := parseTerraformUUID(plan.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	workflow, err := r.updateIdentity(workflowID, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error setting workflow execution identity", err.Error())
		return
	}
	mapExecutionIdentityToState(&plan, workflowID, workflow)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowExecutionIdentityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowExecutionIdentityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(state.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	workflow, err := r.client.GetWorkflow(workflowID)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading workflow execution identity", err.Error())
		return
	}
	mapExecutionIdentityToState(&state, workflowID, workflow)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowExecutionIdentityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkflowExecutionIdentityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := validateExecutionIdentity(plan); err != nil {
		resp.Diagnostics.AddError("Invalid workflow execution identity", err.Error())
		return
	}
	workflowID, err := parseTerraformUUID(plan.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	workflow, err := r.updateIdentity(workflowID, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error updating workflow execution identity", err.Error())
		return
	}
	mapExecutionIdentityToState(&plan, workflowID, workflow)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowExecutionIdentityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkflowExecutionIdentityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(state.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	workflow, err := r.client.GetWorkflow(workflowID)
	if err != nil {
		if apiclient.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error reading workflow", err.Error())
		return
	}
	identity, _ := json.Marshal(workflowExecutionIdentityState{Type: "initiator"})
	request := workflowUpdateRequest(workflow)
	request.ExecutionIdentity = identity
	updated, err := r.client.UpdateWorkflow(workflowID, request)
	if err != nil {
		resp.Diagnostics.AddError("Error resetting workflow execution identity", err.Error())
		return
	}
	if _, err := preserveWorkflowPublication(r.client, workflow, updated, workflow.Published); err != nil {
		resp.Diagnostics.AddError("Error republishing workflow after execution identity reset", err.Error())
	}
}

func (r *WorkflowExecutionIdentityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	workflowID, err := parseTerraformUUID(types.StringValue(req.ID), "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	workflow, err := r.client.GetWorkflow(workflowID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow execution identity", err.Error())
		return
	}
	state := WorkflowExecutionIdentityResourceModel{}
	mapExecutionIdentityToState(&state, workflowID, workflow)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowExecutionIdentityResource) updateIdentity(workflowID string, plan WorkflowExecutionIdentityResourceModel) (*apiclient.Workflow, error) {
	if err := validateExecutionIdentity(plan); err != nil {
		return nil, err
	}
	identity := workflowExecutionIdentityState{Type: plan.Type.ValueString(), ServicePrincipalID: plan.ServicePrincipalID.ValueString()}
	raw, err := json.Marshal(identity)
	if err != nil {
		return nil, err
	}
	workflow, err := r.client.GetWorkflow(workflowID)
	if err != nil {
		return nil, err
	}
	request := workflowUpdateRequest(workflow)
	request.ExecutionIdentity = raw
	updated, err := r.client.UpdateWorkflow(workflowID, request)
	if err != nil {
		return nil, err
	}
	return preserveWorkflowPublication(r.client, workflow, updated, workflow.Published)
}

func validateExecutionIdentity(model WorkflowExecutionIdentityResourceModel) error {
	if model.ServicePrincipalID.IsNull() || model.ServicePrincipalID.IsUnknown() || model.ServicePrincipalID.ValueString() == "" {
		return nil
	}
	if model.Type.IsUnknown() || model.Type.IsNull() || model.Type.ValueString() == "service_principal" {
		return nil
	}
	return fmt.Errorf("service_principal_id is only valid when type is service_principal")
}

func mapExecutionIdentityToState(model *WorkflowExecutionIdentityResourceModel, workflowID string, workflow *apiclient.Workflow) {
	model.ID = types.StringValue(workflowID)
	model.WorkflowID = types.StringValue(workflowID)
	model.WorkflowVersion = types.Int64Value(int64(workflow.Version))
	var identity workflowExecutionIdentityState
	if len(workflow.ExecutionIdentity) > 0 {
		_ = json.Unmarshal(workflow.ExecutionIdentity, &identity)
	}
	if identity.Type == "" {
		identity.Type = "initiator"
	}
	model.Type = types.StringValue(identity.Type)
	model.ServicePrincipalID = optionalString(identity.ServicePrincipalID)
}
