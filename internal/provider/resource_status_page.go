package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &StatusPageResource{}
var _ resource.ResourceWithImportState = &StatusPageResource{}

// StatusPageResource defines the resource implementation.
type StatusPageResource struct {
	client *apiclient.Client
}

// StatusPageResourceModel describes the resource data model.
type StatusPageResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
	IsPublic    types.Bool   `tfsdk:"is_public"`
}

// NewStatusPageResource returns a new resource factory function.
func NewStatusPageResource() resource.Resource {
	return &StatusPageResource{}
}

func (r *StatusPageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (r *StatusPageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat status page.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the status page.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the status page.",
				Required:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The URL slug of the status page.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the status page.",
				Optional:    true,
				Computed:    true,
			},
			"is_public": schema.BoolAttribute{
				Description: "Whether the status page is publicly accessible.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *StatusPageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StatusPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StatusPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateStatusPageRequest{
		Name:     plan.Name.ValueString(),
		Slug:     plan.Slug.ValueString(),
		IsPublic: plan.IsPublic.ValueBool(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.Description = plan.Description.ValueString()
	}

	page, err := r.client.CreateStatusPage(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating status page", err.Error())
		return
	}

	plan.ID = types.StringValue(page.ID)
	plan.Name = types.StringValue(page.Name)
	plan.Slug = types.StringValue(page.Slug)
	plan.Description = types.StringValue(page.Description)
	plan.IsPublic = types.BoolValue(page.IsPublic)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StatusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	page, err := r.client.GetStatusPage(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading status page", err.Error())
		return
	}

	state.Name = types.StringValue(page.Name)
	state.Slug = types.StringValue(page.Slug)
	state.Description = types.StringValue(page.Description)
	state.IsPublic = types.BoolValue(page.IsPublic)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StatusPageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StatusPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state StatusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateStatusPageRequest{
		Name:     plan.Name.ValueString(),
		Slug:     plan.Slug.ValueString(),
		IsPublic: plan.IsPublic.ValueBool(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.Description = plan.Description.ValueString()
	}

	page, err := r.client.UpdateStatusPage(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating status page", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(page.Name)
	plan.Slug = types.StringValue(page.Slug)
	plan.Description = types.StringValue(page.Description)
	plan.IsPublic = types.BoolValue(page.IsPublic)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StatusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStatusPage(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting status page", err.Error())
		return
	}
}

func (r *StatusPageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	page, err := r.client.GetStatusPage(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing status page", err.Error())
		return
	}

	state := StatusPageResourceModel{
		ID:          types.StringValue(page.ID),
		Name:        types.StringValue(page.Name),
		Slug:        types.StringValue(page.Slug),
		Description: types.StringValue(page.Description),
		IsPublic:    types.BoolValue(page.IsPublic),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
