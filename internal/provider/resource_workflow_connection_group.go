package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &WorkflowConnectionGroupResource{}
var _ resource.ResourceWithImportState = &WorkflowConnectionGroupResource{}

type WorkflowConnectionGroupResource struct {
	client *apiclient.Client
}

type WorkflowConnectionGroupResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ConnectionType      types.String `tfsdk:"connection_type"`
	MemberConnectionIDs types.List   `tfsdk:"member_connection_ids"`
	SelectionStrategy   types.String `tfsdk:"selection_strategy"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
}

func NewWorkflowConnectionGroupResource() resource.Resource {
	return &WorkflowConnectionGroupResource{}
}

func (r *WorkflowConnectionGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_workflow_connection_group"
}

func (r *WorkflowConnectionGroupResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat workflow connection group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the workflow connection group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the workflow connection group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connection_type": schema.StringAttribute{
				Description: "The shared connection type for group members.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"member_connection_ids": schema.ListAttribute{
				Description: "Workflow connection IDs included in the group.",
				Required:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"selection_strategy": schema.StringAttribute{
				Description: "Selection strategy for resolving a member connection.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("first_match"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "The group creation timestamp.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "The group update timestamp.",
				Computed:    true,
			},
		},
	}
}

func (r *WorkflowConnectionGroupResource) Configure(
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

func (r *WorkflowConnectionGroupResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan WorkflowConnectionGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	memberIDs, ok := workflowConnectionGroupMembers(ctx, &resp.Diagnostics, plan)
	if !ok {
		return
	}
	group, err := r.client.CreateWorkflowConnectionGroup(apiclient.CreateWorkflowConnectionGroupRequest{
		Name:                plan.Name.ValueString(),
		ConnectionType:      plan.ConnectionType.ValueString(),
		MemberConnectionIDs: memberIDs,
		SelectionStrategy:   plan.SelectionStrategy.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating workflow connection group", err.Error())
		return
	}

	mapWorkflowConnectionGroupToState(ctx, &resp.Diagnostics, &plan, group)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowConnectionGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowConnectionGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "workflow connection group ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading workflow connection group", err.Error())
		return
	}
	group, err := r.client.GetWorkflowConnectionGroup(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading workflow connection group", err.Error())
		return
	}

	mapWorkflowConnectionGroupToState(ctx, &resp.Diagnostics, &state, group)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowConnectionGroupResource) Update(
	ctx context.Context,
	_ resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError(
		"Workflow connection groups cannot be updated in place",
		"Change group attributes to replace the resource.",
	)
}

func (r *WorkflowConnectionGroupResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state WorkflowConnectionGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "workflow connection group ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workflow connection group", err.Error())
		return
	}
	err = r.client.DeleteWorkflowConnectionGroup(id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workflow connection group", err.Error())
		return
	}
}

func (r *WorkflowConnectionGroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := parseTerraformID(types.StringValue(req.ID), "workflow connection group ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow connection group", err.Error())
		return
	}
	group, err := r.client.GetWorkflowConnectionGroup(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow connection group", err.Error())
		return
	}

	state := WorkflowConnectionGroupResourceModel{}
	mapWorkflowConnectionGroupToState(ctx, &resp.Diagnostics, &state, group)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func workflowConnectionGroupMembers(
	ctx context.Context,
	diags *diag.Diagnostics,
	model WorkflowConnectionGroupResourceModel,
) ([]int, bool) {
	var rawMemberIDs []int64
	diags.Append(model.MemberConnectionIDs.ElementsAs(ctx, &rawMemberIDs, false)...)
	if diags.HasError() {
		return nil, false
	}
	memberIDs := make([]int, 0, len(rawMemberIDs))
	for _, id := range rawMemberIDs {
		memberIDs = append(memberIDs, int(id))
	}
	return memberIDs, true
}

func mapWorkflowConnectionGroupToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *WorkflowConnectionGroupResourceModel,
	group *apiclient.WorkflowConnectionGroup,
) {
	model.ID = terraformID(group.ID)
	model.Name = types.StringValue(group.Name)
	model.ConnectionType = types.StringValue(group.ConnectionType)
	model.SelectionStrategy = types.StringValue(group.SelectionStrategy)
	model.CreatedAt = optionalString(group.CreatedAt)
	model.UpdatedAt = optionalString(group.UpdatedAt)
	memberIDs := make([]int64, 0, len(group.MemberConnectionIDs))
	for _, id := range group.MemberConnectionIDs {
		memberIDs = append(memberIDs, int64(id))
	}
	members, memberDiags := types.ListValueFrom(ctx, types.Int64Type, memberIDs)
	diags.Append(memberDiags...)
	model.MemberConnectionIDs = members
}
