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

var _ resource.Resource = &AuthTokenResource{}
var _ resource.ResourceWithImportState = &AuthTokenResource{}

type AuthTokenResource struct {
	client *apiclient.Client
}

type AuthTokenResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Token  types.String `tfsdk:"token"`
	Scopes types.List   `tfsdk:"scopes"`
}

func NewAuthTokenResource() resource.Resource {
	return &AuthTokenResource{}
}

func (r *AuthTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_token"
}

func (r *AuthTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat authentication token.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the auth token.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the auth token.",
				Required:    true,
			},
			"token": schema.StringAttribute{
				Description: "The generated token value. Only available after creation.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scopes": schema.ListAttribute{
				Description: "The list of scopes for the token.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *AuthTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AuthTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AuthTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var scopes []string
	resp.Diagnostics.Append(plan.Scopes.ElementsAs(ctx, &scopes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateAuthTokenRequest{
		Name:   plan.Name.ValueString(),
		Scopes: scopes,
	}

	token, err := r.client.CreateAuthToken(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating auth token", err.Error())
		return
	}

	plan.ID = types.StringValue(token.ID)
	plan.Name = types.StringValue(token.Name)
	plan.Token = types.StringValue(token.Token)

	scopeList, diags := types.ListValueFrom(ctx, types.StringType, token.Scopes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Scopes = scopeList

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AuthTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AuthTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.GetAuthToken(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading auth token", err.Error())
		return
	}

	state.Name = types.StringValue(token.Name)
	// Token value is only available at creation time; preserve existing state
	if token.Token != "" {
		state.Token = types.StringValue(token.Token)
	}

	scopeList, diags := types.ListValueFrom(ctx, types.StringType, token.Scopes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Scopes = scopeList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AuthTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AuthTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AuthTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var scopes []string
	resp.Diagnostics.Append(plan.Scopes.ElementsAs(ctx, &scopes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateAuthTokenRequest{
		Name:   plan.Name.ValueString(),
		Scopes: scopes,
	}

	token, err := r.client.UpdateAuthToken(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating auth token", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Token = state.Token // Preserve token from state
	plan.Name = types.StringValue(token.Name)

	scopeList, diags := types.ListValueFrom(ctx, types.StringType, token.Scopes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Scopes = scopeList

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AuthTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AuthTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAuthToken(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting auth token", err.Error())
		return
	}
}

func (r *AuthTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	token, err := r.client.GetAuthToken(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing auth token", err.Error())
		return
	}

	scopeList, diags := types.ListValueFrom(ctx, types.StringType, token.Scopes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := AuthTokenResourceModel{
		ID:     types.StringValue(token.ID),
		Name:   types.StringValue(token.Name),
		Token:  types.StringValue(token.Token),
		Scopes: scopeList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
