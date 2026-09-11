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

var _ resource.Resource = &ConnectionGroupResource{}
var _ resource.ResourceWithImportState = &ConnectionGroupResource{}

type ConnectionGroupResource struct {
	client *apiclient.Client
}

type ConnectionGroupResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	ProviderID            types.String `tfsdk:"provider_id"`
	UseID                 types.String `tfsdk:"use_id"`
	MemberInstallationIDs types.List   `tfsdk:"member_installation_ids"`
	SelectionStrategy     types.String `tfsdk:"selection_strategy"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

func NewConnectionGroupResource() resource.Resource {
	return &ConnectionGroupResource{}
}

func (r *ConnectionGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_connection_group"
}

func (r *ConnectionGroupResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a UUID-addressed Moneat connector connection group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "UUID resource ID of the connection group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Display name of the connection group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_id": schema.StringAttribute{
				Description: "Connector provider catalog ID shared by group members.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"use_id": schema.StringAttribute{
				Description: "Connector use served by the group.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("workflow_actions"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"member_installation_ids": schema.ListAttribute{
				Description: "UUID installation IDs in deterministic failover order.",
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"selection_strategy": schema.StringAttribute{
				Description: "Selection strategy applied to matching installation tags.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("first_match"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{Computed: true},
			"updated_at": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *ConnectionGroupResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
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

func (r *ConnectionGroupResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan ConnectionGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	members, ok := connectionGroupMembers(ctx, &resp.Diagnostics, plan.MemberInstallationIDs)
	if !ok {
		return
	}
	group, err := r.client.CreateConnectorGroup(apiclient.CreateConnectorGroupRequest{
		Name:                  plan.Name.ValueString(),
		ProviderID:            plan.ProviderID.ValueString(),
		UseID:                 plan.UseID.ValueString(),
		MemberInstallationIDs: members,
		SelectionStrategy:     plan.SelectionStrategy.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating connection group", err.Error())
		return
	}
	mapConnectionGroupToState(ctx, &resp.Diagnostics, &plan, group)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConnectionGroupResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state ConnectionGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := parseTerraformUUID(state.ID, "connection group ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading connection group", err.Error())
		return
	}
	group, err := r.client.GetConnectorGroup(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading connection group", err.Error())
		return
	}
	mapConnectionGroupToState(ctx, &resp.Diagnostics, &state, group)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ConnectionGroupResource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError("Connection groups cannot be updated in place", "Change group attributes to replace the resource.")
}

func (r *ConnectionGroupResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state ConnectionGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := parseTerraformUUID(state.ID, "connection group ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting connection group", err.Error())
		return
	}
	if err := r.client.DeleteConnectorGroup(id); err != nil && !apiclient.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting connection group", err.Error())
	}
}

func (r *ConnectionGroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := parseTerraformUUID(types.StringValue(req.ID), "connection group ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing connection group", err.Error())
		return
	}
	group, err := r.client.GetConnectorGroup(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing connection group", err.Error())
		return
	}
	state := ConnectionGroupResourceModel{}
	mapConnectionGroupToState(ctx, &resp.Diagnostics, &state, group)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func connectionGroupMembers(
	ctx context.Context,
	diags *diag.Diagnostics,
	value types.List,
) ([]string, bool) {
	var members []string
	diags.Append(value.ElementsAs(ctx, &members, false)...)
	for index, member := range members {
		if _, err := parseTerraformUUID(types.StringValue(member), fmt.Sprintf("member installation ID at index %d", index)); err != nil {
			diags.AddError("Invalid connector installation ID", err.Error())
		}
	}
	return members, !diags.HasError()
}

func mapConnectionGroupToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *ConnectionGroupResourceModel,
	group *apiclient.ConnectorGroup,
) {
	model.ID = types.StringValue(group.ID)
	model.Name = types.StringValue(group.Name)
	model.ProviderID = types.StringValue(group.ProviderID)
	model.UseID = types.StringValue(group.UseID)
	model.SelectionStrategy = types.StringValue(group.SelectionStrategy)
	model.CreatedAt = optionalString(group.CreatedAt)
	model.UpdatedAt = optionalString(group.UpdatedAt)
	members, memberDiags := types.ListValueFrom(ctx, types.StringType, group.MemberInstallationIDs)
	diags.Append(memberDiags...)
	model.MemberInstallationIDs = members
}
