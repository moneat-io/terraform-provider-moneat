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

var _ resource.Resource = &StatusPageIncidentResource{}
var _ resource.ResourceWithImportState = &StatusPageIncidentResource{}

type StatusPageIncidentResource struct {
	client *apiclient.Client
}

type StatusPageIncidentResourceModel struct {
	ID          types.String `tfsdk:"id"`
	PageID      types.String `tfsdk:"page_id"`
	Name        types.String `tfsdk:"name"`
	Status      types.String `tfsdk:"status"`
	ImpactLevel types.String `tfsdk:"impact_level"`
	Message     types.String `tfsdk:"message"`
}

func NewStatusPageIncidentResource() resource.Resource {
	return &StatusPageIncidentResource{}
}

func (r *StatusPageIncidentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page_incident"
}

func (r *StatusPageIncidentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat status page incident.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the incident.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"page_id": schema.StringAttribute{
				Description: "The ID of the status page.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the incident.",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the incident (investigating, identified, monitoring, resolved).",
				Required:    true,
			},
			"impact_level": schema.StringAttribute{
				Description: "The impact level (none, minor, major, critical).",
				Required:    true,
			},
			"message": schema.StringAttribute{
				Description: "The incident message.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *StatusPageIncidentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StatusPageIncidentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StatusPageIncidentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateStatusPageIncidentRequest{
		Name:        plan.Name.ValueString(),
		Status:      plan.Status.ValueString(),
		ImpactLevel: plan.ImpactLevel.ValueString(),
	}
	if !plan.Message.IsNull() && !plan.Message.IsUnknown() {
		apiReq.Message = plan.Message.ValueString()
	}

	incident, err := r.client.CreateStatusPageIncident(plan.PageID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating status page incident", err.Error())
		return
	}

	plan.ID = types.StringValue(incident.ID)
	plan.Name = types.StringValue(incident.Name)
	plan.Status = types.StringValue(incident.Status)
	plan.ImpactLevel = types.StringValue(incident.ImpactLevel)
	plan.Message = types.StringValue(incident.Message)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageIncidentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StatusPageIncidentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	incident, err := r.client.GetStatusPageIncident(state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading status page incident", err.Error())
		return
	}

	state.Name = types.StringValue(incident.Name)
	state.Status = types.StringValue(incident.Status)
	state.ImpactLevel = types.StringValue(incident.ImpactLevel)
	state.Message = types.StringValue(incident.Message)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StatusPageIncidentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StatusPageIncidentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state StatusPageIncidentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateStatusPageIncidentRequest{
		Name:        plan.Name.ValueString(),
		Status:      plan.Status.ValueString(),
		ImpactLevel: plan.ImpactLevel.ValueString(),
	}
	if !plan.Message.IsNull() && !plan.Message.IsUnknown() {
		apiReq.Message = plan.Message.ValueString()
	}

	incident, err := r.client.UpdateStatusPageIncident(state.PageID.ValueString(), state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating status page incident", err.Error())
		return
	}

	plan.ID = state.ID
	plan.PageID = state.PageID
	plan.Name = types.StringValue(incident.Name)
	plan.Status = types.StringValue(incident.Status)
	plan.ImpactLevel = types.StringValue(incident.ImpactLevel)
	plan.Message = types.StringValue(incident.Message)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageIncidentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StatusPageIncidentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStatusPageIncident(state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting status page incident", err.Error())
		return
	}
}

func (r *StatusPageIncidentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: page_id/incident_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: page_id/incident_id",
		)
		return
	}

	pageID := parts[0]
	incidentID := parts[1]

	incident, err := r.client.GetStatusPageIncident(pageID, incidentID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing status page incident", err.Error())
		return
	}

	state := StatusPageIncidentResourceModel{
		ID:          types.StringValue(incident.ID),
		PageID:      types.StringValue(pageID),
		Name:        types.StringValue(incident.Name),
		Status:      types.StringValue(incident.Status),
		ImpactLevel: types.StringValue(incident.ImpactLevel),
		Message:     types.StringValue(incident.Message),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
