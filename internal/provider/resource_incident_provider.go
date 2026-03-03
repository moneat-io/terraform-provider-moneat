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

var _ resource.Resource = &IncidentProviderResource{}
var _ resource.ResourceWithImportState = &IncidentProviderResource{}

type IncidentProviderResource struct {
	client *apiclient.Client
}

type IncidentProviderResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Type   types.String `tfsdk:"type"`
	Config types.String `tfsdk:"config"`
}

func NewIncidentProviderResource() resource.Resource {
	return &IncidentProviderResource{}
}

func (r *IncidentProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incident_provider"
}

func (r *IncidentProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat incident provider.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the incident provider.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the incident provider.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the incident provider.",
				Required:    true,
			},
			"config": schema.StringAttribute{
				Description: "The configuration JSON for the incident provider.",
				Required:    true,
			},
		},
	}
}

func (r *IncidentProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IncidentProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IncidentProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateIncidentProviderRequest{
		Name:   plan.Name.ValueString(),
		Type:   plan.Type.ValueString(),
		Config: plan.Config.ValueString(),
	}

	provider, err := r.client.CreateIncidentProvider(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating incident provider", err.Error())
		return
	}

	plan.ID = types.StringValue(provider.ID)
	plan.Name = types.StringValue(provider.Name)
	plan.Type = types.StringValue(provider.Type)
	plan.Config = types.StringValue(provider.Config)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IncidentProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IncidentProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	provider, err := r.client.GetIncidentProvider(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading incident provider", err.Error())
		return
	}

	state.Name = types.StringValue(provider.Name)
	state.Type = types.StringValue(provider.Type)
	state.Config = types.StringValue(provider.Config)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *IncidentProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IncidentProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IncidentProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateIncidentProviderRequest{
		Name:   plan.Name.ValueString(),
		Type:   plan.Type.ValueString(),
		Config: plan.Config.ValueString(),
	}

	provider, err := r.client.UpdateIncidentProvider(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating incident provider", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(provider.Name)
	plan.Type = types.StringValue(provider.Type)
	plan.Config = types.StringValue(provider.Config)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IncidentProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IncidentProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIncidentProvider(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting incident provider", err.Error())
		return
	}
}

func (r *IncidentProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	provider, err := r.client.GetIncidentProvider(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing incident provider", err.Error())
		return
	}

	state := IncidentProviderResourceModel{
		ID:     types.StringValue(provider.ID),
		Name:   types.StringValue(provider.Name),
		Type:   types.StringValue(provider.Type),
		Config: types.StringValue(provider.Config),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
