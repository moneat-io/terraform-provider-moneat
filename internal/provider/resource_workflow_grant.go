package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &WorkflowGrantResource{}
var _ resource.ResourceWithImportState = &WorkflowGrantResource{}

type WorkflowGrantResource struct {
	client *apiclient.Client
}

type WorkflowGrantResourceModel struct {
	ID             types.String `tfsdk:"id"`
	WorkflowID     types.String `tfsdk:"workflow_id"`
	UserID         types.String `tfsdk:"user_id"`
	Role           types.String `tfsdk:"role"`
	AllowedActions types.Set    `tfsdk:"allowed_actions"`
	GrantedBy      types.String `tfsdk:"granted_by"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func NewWorkflowGrantResource() resource.Resource {
	return &WorkflowGrantResource{}
}

func (r *WorkflowGrantResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_grant"
}

func (r *WorkflowGrantResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages one organization-scoped workflow access grant.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"workflow_id":     schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"user_id":         schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"role":            schema.StringAttribute{Required: true},
			"allowed_actions": schema.SetAttribute{Optional: true, Computed: true, ElementType: types.StringType},
			"granted_by":      schema.StringAttribute{Computed: true},
			"created_at":      schema.StringAttribute{Computed: true},
			"updated_at":      schema.StringAttribute{Computed: true},
		},
	}
}

func (r *WorkflowGrantResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkflowGrantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkflowGrantResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(plan.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	userID, err := parseTerraformUUID(plan.UserID, "workflow grant user ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow grant user ID", err.Error())
		return
	}
	actions, ok := workflowGrantActions(ctx, &resp.Diagnostics, plan.AllowedActions)
	if !ok {
		return
	}
	grant, err := r.client.UpsertWorkflowGrant(
		workflowID,
		apiclient.WorkflowGrantRequest{UserID: userID, Role: plan.Role.ValueString(), AllowedActions: actions},
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating workflow grant", err.Error())
		return
	}
	mapWorkflowGrantToState(ctx, &resp.Diagnostics, &plan, grant)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowGrantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowGrantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(state.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	grants, err := r.client.ListWorkflowGrants(workflowID)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading workflow grant", err.Error())
		return
	}
	for _, grant := range grants {
		if grant.ID == state.ID.ValueString() {
			mapWorkflowGrantToState(ctx, &resp.Diagnostics, &state, &grant)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}
	resp.State.RemoveResource(ctx)
}

func (r *WorkflowGrantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkflowGrantResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(plan.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	grantID, err := parseTerraformUUID(plan.ID, "workflow grant ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow grant ID", err.Error())
		return
	}
	userID, err := parseTerraformUUID(plan.UserID, "workflow grant user ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow grant user ID", err.Error())
		return
	}
	actions, ok := workflowGrantActions(ctx, &resp.Diagnostics, plan.AllowedActions)
	if !ok {
		return
	}
	grant, err := r.client.UpdateWorkflowGrant(
		workflowID,
		grantID,
		apiclient.WorkflowGrantRequest{UserID: userID, Role: plan.Role.ValueString(), AllowedActions: actions},
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating workflow grant", err.Error())
		return
	}
	mapWorkflowGrantToState(ctx, &resp.Diagnostics, &plan, grant)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowGrantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkflowGrantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(state.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	grantID, err := parseTerraformUUID(state.ID, "workflow grant ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow grant ID", err.Error())
		return
	}
	if err := r.client.DeleteWorkflowGrant(workflowID, grantID); err != nil && !apiclient.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting workflow grant", err.Error())
	}
}

func (r *WorkflowGrantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := splitCompositeID(req.ID)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid workflow grant ID", "Import IDs must be <workflow UUID>/<grant UUID>")
		return
	}
	workflowID, err := parseTerraformUUID(types.StringValue(parts[0]), "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	grantID, err := parseTerraformUUID(types.StringValue(parts[1]), "workflow grant ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow grant ID", err.Error())
		return
	}
	grants, err := r.client.ListWorkflowGrants(workflowID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow grant", err.Error())
		return
	}
	for _, grant := range grants {
		if grant.ID == grantID {
			state := WorkflowGrantResourceModel{}
			mapWorkflowGrantToState(ctx, &resp.Diagnostics, &state, &grant)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}
	resp.Diagnostics.AddError("Workflow grant not found", parts[1])
}

func mapWorkflowGrantToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *WorkflowGrantResourceModel,
	grant *apiclient.WorkflowGrant,
) {
	model.ID = types.StringValue(grant.ID)
	model.WorkflowID = types.StringValue(grant.WorkflowID)
	model.UserID = types.StringValue(grant.UserID)
	model.Role = types.StringValue(grant.Role)
	if grant.AllowedActions == nil {
		model.AllowedActions = types.SetNull(types.StringType)
	} else {
		actions, actionDiags := types.SetValueFrom(ctx, types.StringType, grant.AllowedActions)
		diags.Append(actionDiags...)
		model.AllowedActions = actions
	}
	model.GrantedBy = optionalString(grant.GrantedBy)
	model.CreatedAt = optionalString(grant.CreatedAt)
	model.UpdatedAt = optionalString(grant.UpdatedAt)
}

func workflowGrantActions(ctx context.Context, diags *diag.Diagnostics, value types.Set) ([]string, bool) {
	if value.IsNull() || value.IsUnknown() {
		return nil, true
	}
	actions := make([]string, 0)
	diags.Append(value.ElementsAs(ctx, &actions, false)...)
	return actions, !diags.HasError()
}

func splitCompositeID(value string) []string {
	for index := 0; index < len(value); index++ {
		if value[index] == '/' {
			return []string{value[:index], value[index+1:]}
		}
	}
	return nil
}
