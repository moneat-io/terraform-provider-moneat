package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &ProjectTargetResource{}
var _ resource.ResourceWithImportState = &ProjectTargetResource{}

type ProjectTargetResource struct {
	client *apiclient.Client
}

type ProjectTargetResourceModel struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	Target    types.String `tfsdk:"target"`
	DSN       types.String `tfsdk:"dsn"`
}

func NewProjectTargetResource() resource.Resource {
	return &ProjectTargetResource{}
}

func (r *ProjectTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_target"
}

func (r *ProjectTargetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat project target and DSN.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The project target import identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The Moneat project ID or resource ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target": schema.StringAttribute{
				Description: "The target platform or service name for the generated DSN.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dsn": schema.StringAttribute{
				Description: "The generated Sentry-compatible DSN for this target.",
				Computed:    true,
			},
		},
	}
}

func (r *ProjectTargetResource) Configure(
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

func (r *ProjectTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectTargetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.AddProjectTarget(
		plan.ProjectID.ValueString(),
		apiclient.AddProjectTargetRequest{Target: plan.Target.ValueString()},
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating project target", err.Error())
		return
	}

	mapProjectTargetToState(&plan, key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectTargetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GetProjectTarget(state.ProjectID.ValueString(), state.Target.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading project target", err.Error())
		return
	}

	mapProjectTargetToState(&state, key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProjectTargetResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan ProjectTargetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectTargetResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
	// The API does not expose project target deletion.
}

func (r *ProjectTargetResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: project_id/target",
		)
		return
	}

	key, err := r.client.GetProjectTarget(parts[0], parts[1])
	if err != nil {
		resp.Diagnostics.AddError("Error importing project target", err.Error())
		return
	}

	state := ProjectTargetResourceModel{
		ProjectID: types.StringValue(parts[0]),
		Target:    types.StringValue(parts[1]),
	}
	mapProjectTargetToState(&state, key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func mapProjectTargetToState(model *ProjectTargetResourceModel, key *apiclient.ProjectKey) {
	model.ID = types.StringValue(model.ProjectID.ValueString() + "/" + key.PlatformTarget)
	model.Target = types.StringValue(key.PlatformTarget)
	model.DSN = types.StringValue(key.DSN)
}
