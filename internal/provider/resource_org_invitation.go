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

var _ resource.Resource = &OrgInvitationResource{}
var _ resource.ResourceWithImportState = &OrgInvitationResource{}

type OrgInvitationResource struct {
	client *apiclient.Client
}

type OrgInvitationResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

func NewOrgInvitationResource() resource.Resource {
	return &OrgInvitationResource{}
}

func (r *OrgInvitationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_invitation"
}

func (r *OrgInvitationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat organization invitation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the invitation.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Description: "The email address to invite.",
				Required:    true,
			},
			"role": schema.StringAttribute{
				Description: "The role for the invited member (admin, member, viewer).",
				Required:    true,
			},
		},
	}
}

func (r *OrgInvitationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrgInvitationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrgInvitationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateOrgInvitationRequest{
		Email: plan.Email.ValueString(),
		Role:  plan.Role.ValueString(),
	}

	invitation, err := r.client.CreateOrgInvitation(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating organization invitation", err.Error())
		return
	}

	plan.ID = types.StringValue(invitation.ID)
	plan.Email = types.StringValue(invitation.Email)
	plan.Role = types.StringValue(invitation.Role)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrgInvitationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrgInvitationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invitation, err := r.client.GetOrgInvitation(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading organization invitation", err.Error())
		return
	}

	state.Email = types.StringValue(invitation.Email)
	state.Role = types.StringValue(invitation.Role)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrgInvitationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrgInvitationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OrgInvitationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateOrgInvitationRequest{
		Email: plan.Email.ValueString(),
		Role:  plan.Role.ValueString(),
	}

	invitation, err := r.client.UpdateOrgInvitation(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating organization invitation", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Email = types.StringValue(invitation.Email)
	plan.Role = types.StringValue(invitation.Role)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrgInvitationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrgInvitationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrgInvitation(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting organization invitation", err.Error())
		return
	}
}

func (r *OrgInvitationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	invitation, err := r.client.GetOrgInvitation(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing organization invitation", err.Error())
		return
	}

	state := OrgInvitationResourceModel{
		ID:    types.StringValue(invitation.ID),
		Email: types.StringValue(invitation.Email),
		Role:  types.StringValue(invitation.Role),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
