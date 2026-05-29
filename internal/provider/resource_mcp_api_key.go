package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &McpAPIKeyResource{}
var _ resource.ResourceWithImportState = &McpAPIKeyResource{}

type McpAPIKeyResource struct {
	client *apiclient.Client
}

type McpAPIKeyResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Key              types.String `tfsdk:"key"`
	KeyPrefix        types.String `tfsdk:"key_prefix"`
	EnabledTools     types.List   `tfsdk:"enabled_tools"`
	EnabledResources types.List   `tfsdk:"enabled_resources"`
	ExpiresInDays    types.Int64  `tfsdk:"expires_in_days"`
	ExpiresAt        types.String `tfsdk:"expires_at"`
}

func NewMcpAPIKeyResource() resource.Resource {
	return &McpAPIKeyResource{}
}

func (r *McpAPIKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mcp_api_key"
}

func (r *McpAPIKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat MCP API key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the MCP API key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the MCP API key.",
				Required:    true,
			},
			"key": schema.StringAttribute{
				Description: "The generated MCP API key value. Only available after creation.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_prefix": schema.StringAttribute{
				Description: "The display prefix for the generated key.",
				Computed:    true,
			},
			"enabled_tools": schema.ListAttribute{
				Description: "The MCP tool names this key can call.",
				Required:    true,
				ElementType: types.StringType,
			},
			"enabled_resources": schema.ListAttribute{
				Description: "The MCP resource URIs this key can read.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default: listdefault.StaticValue(
					types.ListValueMust(types.StringType, []attr.Value{}),
				),
			},
			"expires_in_days": schema.Int64Attribute{
				Description: "Optional relative expiration in days from create or update time.",
				Optional:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "The absolute expiration timestamp returned by the API.",
				Computed:    true,
			},
		},
	}
}

func (r *McpAPIKeyResource) Configure(
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

func (r *McpAPIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan McpAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enabledTools, enabledResources, ok := mcpAPIKeyLists(ctx, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	apiReq := apiclient.CreateMcpAPIKeyRequest{
		Name:             plan.Name.ValueString(),
		EnabledTools:     enabledTools,
		EnabledResources: enabledResources,
		ExpiresInDays:    intPointerFromInt64(plan.ExpiresInDays),
	}
	key, err := r.client.CreateMcpAPIKey(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating MCP API key", err.Error())
		return
	}

	mapMcpAPIKeyToState(ctx, &resp.Diagnostics, &plan, key, plan.Key, plan.ExpiresInDays)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *McpAPIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state McpAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "MCP API key ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading MCP API key", err.Error())
		return
	}
	key, err := r.client.GetMcpAPIKey(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading MCP API key", err.Error())
		return
	}

	mapMcpAPIKeyToState(ctx, &resp.Diagnostics, &state, key, state.Key, state.ExpiresInDays)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *McpAPIKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan McpAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(plan.ID, "MCP API key ID")
	if err != nil {
		resp.Diagnostics.AddError("Error updating MCP API key", err.Error())
		return
	}
	enabledTools, enabledResources, ok := mcpAPIKeyLists(ctx, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	apiReq := apiclient.UpdateMcpAPIKeyRequest{
		Name:             plan.Name.ValueString(),
		EnabledTools:     enabledTools,
		EnabledResources: enabledResources,
		ExpiresInDays:    intPointerFromInt64(plan.ExpiresInDays),
	}
	key, err := r.client.UpdateMcpAPIKey(id, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating MCP API key", err.Error())
		return
	}

	mapMcpAPIKeyToState(ctx, &resp.Diagnostics, &plan, key, plan.Key, plan.ExpiresInDays)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *McpAPIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state McpAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "MCP API key ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting MCP API key", err.Error())
		return
	}
	err = r.client.DeleteMcpAPIKey(id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting MCP API key", err.Error())
		return
	}
}

func (r *McpAPIKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := parseTerraformID(types.StringValue(req.ID), "MCP API key ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing MCP API key", err.Error())
		return
	}
	key, err := r.client.GetMcpAPIKey(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing MCP API key", err.Error())
		return
	}

	state := McpAPIKeyResourceModel{}
	mapMcpAPIKeyToState(ctx, &resp.Diagnostics, &state, key, types.StringNull(), types.Int64Null())
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func mcpAPIKeyLists(
	ctx context.Context,
	model McpAPIKeyResourceModel,
	diags *diag.Diagnostics,
) ([]string, []string, bool) {
	var enabledTools []string
	var enabledResources []string
	diags.Append(model.EnabledTools.ElementsAs(ctx, &enabledTools, false)...)
	diags.Append(model.EnabledResources.ElementsAs(ctx, &enabledResources, false)...)
	if diags.HasError() {
		return nil, nil, false
	}
	return enabledTools, enabledResources, true
}

func mapMcpAPIKeyToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *McpAPIKeyResourceModel,
	key *apiclient.McpAPIKey,
	existingKey types.String,
	expiresInDays types.Int64,
) {
	model.ID = terraformID(key.ID)
	model.Name = types.StringValue(key.Name)
	model.KeyPrefix = types.StringValue(key.KeyPrefix)
	if key.Key != "" {
		model.Key = types.StringValue(key.Key)
	} else {
		model.Key = existingKey
	}
	model.ExpiresInDays = expiresInDays
	if key.ExpiresAt == "" {
		model.ExpiresAt = types.StringNull()
	} else {
		model.ExpiresAt = types.StringValue(key.ExpiresAt)
	}

	tools, toolDiags := types.ListValueFrom(ctx, types.StringType, key.EnabledTools)
	diags.Append(toolDiags...)
	resources, resourceDiags := types.ListValueFrom(ctx, types.StringType, key.EnabledResources)
	diags.Append(resourceDiags...)
	if diags.HasError() {
		return
	}
	model.EnabledTools = tools
	model.EnabledResources = resources
}

func intPointerFromInt64(value types.Int64) *int {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	result := int(value.ValueInt64())
	return &result
}
