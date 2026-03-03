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

var _ resource.Resource = &OrgMemberResource{}
var _ resource.ResourceWithImportState = &OrgMemberResource{}

type OrgMemberResource struct {
	client *apiclient.Client
}

type OrgMemberResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

func NewOrgMemberResource() resource.Resource {
	return &OrgMemberResource{}
}

func (r *OrgMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_member"
}

func (r *OrgMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat organization member. Members cannot be created directly; " +
			"use moneat_org_invitation to invite members, then import existing members.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the organization member.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Description: "The email address of the member.",
				Computed:    true,
			},
			"role": schema.StringAttribute{
				Description: "The role of the member (admin, member, viewer).",
				Required:    true,
			},
		},
	}
}

func (r *OrgMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrgMemberResource) Create(_ context.Context, _ resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Create not supported",
		"Organization members cannot be created directly. Use moneat_org_invitation to invite members, "+
			"then import existing members with 'terraform import moneat_org_member.<name> <member_id>'.",
	)
}

func (r *OrgMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrgMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members, err := r.client.ListOrgMembers()
	if err != nil {
		resp.Diagnostics.AddError("Error reading organization member", err.Error())
		return
	}

	memberID := state.ID.ValueString()
	for _, m := range members {
		if m.ID == memberID {
			state.Email = types.StringValue(m.Email)
			state.Role = types.StringValue(m.Role)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// Not found - remove from state
	resp.State.RemoveResource(ctx)
}

func (r *OrgMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrgMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OrgMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateOrgMemberRequest{
		Role: plan.Role.ValueString(),
	}

	member, err := r.client.UpdateOrgMember(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating organization member", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Email = types.StringValue(member.Email)
	plan.Role = types.StringValue(member.Role)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrgMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrgMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrgMember(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting organization member", err.Error())
		return
	}
}

func (r *OrgMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	members, err := r.client.ListOrgMembers()
	if err != nil {
		resp.Diagnostics.AddError("Error importing organization member", err.Error())
		return
	}

	for _, m := range members {
		if m.ID == req.ID {
			state := OrgMemberResourceModel{
				ID:    types.StringValue(m.ID),
				Email: types.StringValue(m.Email),
				Role:  types.StringValue(m.Role),
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	resp.Diagnostics.AddError(
		"Member not found",
		fmt.Sprintf("No organization member found with ID: %s", req.ID),
	)
}
