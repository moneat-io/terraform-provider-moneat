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

var _ resource.Resource = &IncidentRoutingRuleResource{}
var _ resource.ResourceWithImportState = &IncidentRoutingRuleResource{}

type IncidentRoutingRuleResource struct {
	client *apiclient.Client
}

type IncidentRoutingRuleResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ProviderID    types.String `tfsdk:"provider_id"`
	Name          types.String `tfsdk:"name"`
	Condition     types.String `tfsdk:"condition"`
	TargetService types.String `tfsdk:"target_service"`
}

func NewIncidentRoutingRuleResource() resource.Resource {
	return &IncidentRoutingRuleResource{}
}

func (r *IncidentRoutingRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incident_routing_rule"
}

func (r *IncidentRoutingRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat incident routing rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the routing rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"provider_id": schema.StringAttribute{
				Description: "The ID of the incident provider.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the routing rule.",
				Required:    true,
			},
			"condition": schema.StringAttribute{
				Description: "The condition expression for the routing rule.",
				Required:    true,
			},
			"target_service": schema.StringAttribute{
				Description: "The target service for matched incidents.",
				Required:    true,
			},
		},
	}
}

func (r *IncidentRoutingRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IncidentRoutingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IncidentRoutingRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateIncidentRoutingRuleRequest{
		Name:          plan.Name.ValueString(),
		Condition:     plan.Condition.ValueString(),
		TargetService: plan.TargetService.ValueString(),
	}

	rule, err := r.client.CreateIncidentRoutingRule(plan.ProviderID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating incident routing rule", err.Error())
		return
	}

	plan.ID = types.StringValue(rule.ID)
	plan.Name = types.StringValue(rule.Name)
	plan.Condition = types.StringValue(rule.Condition)
	plan.TargetService = types.StringValue(rule.TargetService)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IncidentRoutingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IncidentRoutingRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetIncidentRoutingRule(state.ProviderID.ValueString(), state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading incident routing rule", err.Error())
		return
	}

	state.Name = types.StringValue(rule.Name)
	state.Condition = types.StringValue(rule.Condition)
	state.TargetService = types.StringValue(rule.TargetService)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *IncidentRoutingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IncidentRoutingRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IncidentRoutingRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateIncidentRoutingRuleRequest{
		Name:          plan.Name.ValueString(),
		Condition:     plan.Condition.ValueString(),
		TargetService: plan.TargetService.ValueString(),
	}

	rule, err := r.client.UpdateIncidentRoutingRule(state.ProviderID.ValueString(), state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating incident routing rule", err.Error())
		return
	}

	plan.ID = state.ID
	plan.ProviderID = state.ProviderID
	plan.Name = types.StringValue(rule.Name)
	plan.Condition = types.StringValue(rule.Condition)
	plan.TargetService = types.StringValue(rule.TargetService)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IncidentRoutingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IncidentRoutingRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIncidentRoutingRule(state.ProviderID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting incident routing rule", err.Error())
		return
	}
}

func (r *IncidentRoutingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: provider_id/rule_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: provider_id/rule_id",
		)
		return
	}

	providerID := parts[0]
	ruleID := parts[1]

	rule, err := r.client.GetIncidentRoutingRule(providerID, ruleID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing incident routing rule", err.Error())
		return
	}

	state := IncidentRoutingRuleResourceModel{
		ID:            types.StringValue(rule.ID),
		ProviderID:    types.StringValue(providerID),
		Name:          types.StringValue(rule.Name),
		Condition:     types.StringValue(rule.Condition),
		TargetService: types.StringValue(rule.TargetService),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
