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

var _ resource.Resource = &DashboardResource{}
var _ resource.ResourceWithImportState = &DashboardResource{}

// DashboardResource defines the resource implementation.
type DashboardResource struct {
	client *apiclient.Client
}

// DashboardResourceModel describes the resource data model.
type DashboardResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// NewDashboardResource returns a new resource factory function.
func NewDashboardResource() resource.Resource {
	return &DashboardResource{}
}

func (r *DashboardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func (r *DashboardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat dashboard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the dashboard.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the dashboard.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the dashboard.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *DashboardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DashboardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateDashboardRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.Description = plan.Description.ValueString()
	}

	dashboard, err := r.client.CreateDashboard(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating dashboard", err.Error())
		return
	}

	plan.ID = types.StringValue(dashboard.ID)
	plan.Name = types.StringValue(dashboard.Name)
	plan.Description = types.StringValue(dashboard.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DashboardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dashboard, err := r.client.GetDashboard(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading dashboard", err.Error())
		return
	}

	state.Name = types.StringValue(dashboard.Name)
	state.Description = types.StringValue(dashboard.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DashboardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state DashboardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateDashboardRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.Description = plan.Description.ValueString()
	}

	dashboard, err := r.client.UpdateDashboard(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating dashboard", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(dashboard.Name)
	plan.Description = types.StringValue(dashboard.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DashboardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDashboard(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting dashboard", err.Error())
		return
	}
}

func (r *DashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	dashboard, err := r.client.GetDashboard(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing dashboard", err.Error())
		return
	}

	state := DashboardResourceModel{
		ID:          types.StringValue(dashboard.ID),
		Name:        types.StringValue(dashboard.Name),
		Description: types.StringValue(dashboard.Description),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
