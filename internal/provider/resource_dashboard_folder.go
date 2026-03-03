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

var _ resource.Resource = &DashboardFolderResource{}
var _ resource.ResourceWithImportState = &DashboardFolderResource{}

type DashboardFolderResource struct {
	client *apiclient.Client
}

type DashboardFolderResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewDashboardFolderResource() resource.Resource {
	return &DashboardFolderResource{}
}

func (r *DashboardFolderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard_folder"
}

func (r *DashboardFolderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat dashboard folder.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the dashboard folder.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the dashboard folder.",
				Required:    true,
			},
		},
	}
}

func (r *DashboardFolderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DashboardFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DashboardFolderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateDashboardFolderRequest{
		Name: plan.Name.ValueString(),
	}

	folder, err := r.client.CreateDashboardFolder(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating dashboard folder", err.Error())
		return
	}

	plan.ID = types.StringValue(folder.ID)
	plan.Name = types.StringValue(folder.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DashboardFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DashboardFolderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folder, err := r.client.GetDashboardFolder(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading dashboard folder", err.Error())
		return
	}

	state.Name = types.StringValue(folder.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DashboardFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DashboardFolderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state DashboardFolderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateDashboardFolderRequest{
		Name: plan.Name.ValueString(),
	}

	folder, err := r.client.UpdateDashboardFolder(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating dashboard folder", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(folder.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DashboardFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DashboardFolderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDashboardFolder(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting dashboard folder", err.Error())
		return
	}
}

func (r *DashboardFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	folder, err := r.client.GetDashboardFolder(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing dashboard folder", err.Error())
		return
	}

	state := DashboardFolderResourceModel{
		ID:   types.StringValue(folder.ID),
		Name: types.StringValue(folder.Name),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
