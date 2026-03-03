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

var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client *apiclient.Client
}

// ProjectResourceModel describes the resource data model.
type ProjectResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Platform  types.String `tfsdk:"platform"`
	Framework types.String `tfsdk:"framework"`
}

// NewProjectResource returns a new resource factory function.
func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

func (r *ProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the project.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the project.",
				Required:    true,
			},
			"platform": schema.StringAttribute{
				Description: "The platform of the project (e.g., python, javascript, go, kotlin).",
				Required:    true,
			},
			"framework": schema.StringAttribute{
				Description: "The framework used by the project (e.g., django, react, gin).",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *ProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateProjectRequest{
		Name:     plan.Name.ValueString(),
		Platform: plan.Platform.ValueString(),
	}
	if !plan.Framework.IsNull() && !plan.Framework.IsUnknown() {
		apiReq.Framework = plan.Framework.ValueString()
	}

	project, err := r.client.CreateProject(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", err.Error())
		return
	}

	plan.ID = types.StringValue(project.ID)
	plan.Name = types.StringValue(project.Name)
	plan.Platform = types.StringValue(project.Platform)
	if project.Framework != "" {
		plan.Framework = types.StringValue(project.Framework)
	} else {
		plan.Framework = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetProject(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	state.Name = types.StringValue(project.Name)
	state.Platform = types.StringValue(project.Platform)
	if project.Framework != "" {
		state.Framework = types.StringValue(project.Framework)
	} else {
		state.Framework = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateProjectRequest{
		Name:     plan.Name.ValueString(),
		Platform: plan.Platform.ValueString(),
	}
	if !plan.Framework.IsNull() && !plan.Framework.IsUnknown() {
		apiReq.Framework = plan.Framework.ValueString()
	}

	project, err := r.client.UpdateProject(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(project.Name)
	plan.Platform = types.StringValue(project.Platform)
	if project.Framework != "" {
		plan.Framework = types.StringValue(project.Framework)
	} else {
		plan.Framework = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting project", err.Error())
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	project, err := r.client.GetProject(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing project", err.Error())
		return
	}

	state := ProjectResourceModel{
		ID:       types.StringValue(project.ID),
		Name:     types.StringValue(project.Name),
		Platform: types.StringValue(project.Platform),
	}
	if project.Framework != "" {
		state.Framework = types.StringValue(project.Framework)
	} else {
		state.Framework = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
