package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &WorkflowConnectionResource{}
var _ resource.ResourceWithImportState = &WorkflowConnectionResource{}

type WorkflowConnectionResource struct {
	client *apiclient.Client
}

type WorkflowConnectionResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Type           types.String `tfsdk:"type"`
	Name           types.String `tfsdk:"name"`
	IdentifierTags types.Map    `tfsdk:"identifier_tags"`
	Secret         types.String `tfsdk:"secret"`
	LastFour       types.String `tfsdk:"last_four"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func NewWorkflowConnectionResource() resource.Resource {
	return &WorkflowConnectionResource{}
}

func (r *WorkflowConnectionResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_workflow_connection"
}

func (r *WorkflowConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a vaulted Moneat workflow connection secret.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the workflow connection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The connection type, such as jira, github, pagerduty, or servicenow.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the workflow connection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"identifier_tags": schema.MapAttribute{
				Description: "Tags used by workflow connection groups to select a connection.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"secret": schema.StringAttribute{
				Description: "The connection secret. Changing this value rotates the stored credential.",
				Required:    true,
				Sensitive:   true,
			},
			"last_four": schema.StringAttribute{
				Description: "The last four characters of the stored secret, when available.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The connection creation timestamp.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "The connection update timestamp.",
				Computed:    true,
			},
		},
	}
}

func (r *WorkflowConnectionResource) Configure(
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

func (r *WorkflowConnectionResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan WorkflowConnectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags, ok := workflowConnectionTags(ctx, plan, resp)
	if !ok {
		return
	}
	connection, err := r.client.CreateWorkflowConnection(apiclient.CreateWorkflowConnectionRequest{
		Type:           plan.Type.ValueString(),
		Name:           plan.Name.ValueString(),
		IdentifierTags: tags,
		Secret:         plan.Secret.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating workflow connection", err.Error())
		return
	}

	mapWorkflowConnectionToState(ctx, &resp.Diagnostics, &plan, connection, plan.Secret)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "workflow connection ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading workflow connection", err.Error())
		return
	}
	connection, err := r.client.GetWorkflowConnection(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading workflow connection", err.Error())
		return
	}

	mapWorkflowConnectionToState(ctx, &resp.Diagnostics, &state, connection, state.Secret)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowConnectionResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan WorkflowConnectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(plan.ID, "workflow connection ID")
	if err != nil {
		resp.Diagnostics.AddError("Error rotating workflow connection", err.Error())
		return
	}
	connection, err := r.client.RotateWorkflowConnection(id, apiclient.RotateWorkflowConnectionRequest{
		Secret: plan.Secret.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error rotating workflow connection", err.Error())
		return
	}

	mapWorkflowConnectionToState(ctx, &resp.Diagnostics, &plan, connection, plan.Secret)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowConnectionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state WorkflowConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "workflow connection ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workflow connection", err.Error())
		return
	}
	err = r.client.DeleteWorkflowConnection(id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workflow connection", err.Error())
		return
	}
}

func (r *WorkflowConnectionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := parseTerraformID(types.StringValue(req.ID), "workflow connection ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow connection", err.Error())
		return
	}
	connection, err := r.client.GetWorkflowConnection(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow connection", err.Error())
		return
	}

	state := WorkflowConnectionResourceModel{Secret: types.StringNull()}
	mapWorkflowConnectionToState(ctx, &resp.Diagnostics, &state, connection, types.StringNull())
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func workflowConnectionTags(
	ctx context.Context,
	model WorkflowConnectionResourceModel,
	resp *resource.CreateResponse,
) (map[string]string, bool) {
	tags := map[string]string{}
	if model.IdentifierTags.IsNull() || model.IdentifierTags.IsUnknown() {
		return tags, true
	}
	resp.Diagnostics.Append(model.IdentifierTags.ElementsAs(ctx, &tags, false)...)
	return tags, !resp.Diagnostics.HasError()
}

func mapWorkflowConnectionToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *WorkflowConnectionResourceModel,
	connection *apiclient.WorkflowConnection,
	existingSecret types.String,
) {
	model.ID = terraformID(connection.ID)
	model.Type = types.StringValue(connection.Type)
	model.Name = types.StringValue(connection.Name)
	model.Secret = existingSecret
	model.LastFour = optionalString(connection.LastFour)
	model.CreatedAt = optionalString(connection.CreatedAt)
	model.UpdatedAt = optionalString(connection.UpdatedAt)
	tags, tagDiags := types.MapValueFrom(ctx, types.StringType, connection.IdentifierTags)
	diags.Append(tagDiags...)
	model.IdentifierTags = tags
}

func optionalString(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}
