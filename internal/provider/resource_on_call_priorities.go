package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &OnCallPrioritiesResource{}

type OnCallPrioritiesResource struct {
	client *apiclient.Client
}

type OnCallPriorityModel struct {
	Level    types.String `tfsdk:"level"`
	Label    types.String `tfsdk:"label"`
	Pageable types.Bool   `tfsdk:"pageable"`
}

type OnCallPrioritiesResourceModel struct {
	Priorities []OnCallPriorityModel `tfsdk:"priorities"`
}

func NewOnCallPrioritiesResource() resource.Resource {
	return &OnCallPrioritiesResource{}
}

func (r *OnCallPrioritiesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_on_call_priorities"
}

func (r *OnCallPrioritiesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Moneat on-call priority levels. This is a singleton resource (Enterprise only).",
		Attributes: map[string]schema.Attribute{
			"priorities": schema.ListNestedAttribute{
				Description: "List of priority levels.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"level": schema.StringAttribute{
							Description: "The priority level identifier (e.g., P1, P2).",
							Required:    true,
						},
						"label": schema.StringAttribute{
							Description: "Display label for the priority level.",
							Required:    true,
						},
						"pageable": schema.BoolAttribute{
							Description: "Whether this priority level triggers paging.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *OnCallPrioritiesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OnCallPrioritiesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OnCallPrioritiesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := r.buildAPIRequest(plan)
	result, err := r.client.UpdateOnCallPriorities(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error setting on-call priorities", err.Error())
		return
	}

	r.mapToState(&plan, result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OnCallPrioritiesResource) Read(ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse) {
	result, err := r.client.GetOnCallPriorities()
	if err != nil {
		resp.Diagnostics.AddError("Error reading on-call priorities", err.Error())
		return
	}

	var state OnCallPrioritiesResourceModel
	r.mapToState(&state, result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OnCallPrioritiesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OnCallPrioritiesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := r.buildAPIRequest(plan)
	result, err := r.client.UpdateOnCallPriorities(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating on-call priorities", err.Error())
		return
	}

	r.mapToState(&plan, result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OnCallPrioritiesResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton resource - reset to empty on delete
	apiReq := apiclient.UpdateOnCallPrioritiesRequest{
		Priorities: []apiclient.OnCallPriority{},
	}
	_, err := r.client.UpdateOnCallPriorities(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error resetting on-call priorities", err.Error())
		return
	}
}

func (r *OnCallPrioritiesResource) buildAPIRequest(
	model OnCallPrioritiesResourceModel,
) apiclient.UpdateOnCallPrioritiesRequest {
	priorities := make([]apiclient.OnCallPriority, len(model.Priorities))
	for i, p := range model.Priorities {
		priorities[i] = apiclient.OnCallPriority{
			Level:    p.Level.ValueString(),
			Label:    p.Label.ValueString(),
			Pageable: p.Pageable.ValueBool(),
		}
	}
	return apiclient.UpdateOnCallPrioritiesRequest{Priorities: priorities}
}

func (r *OnCallPrioritiesResource) mapToState(
	model *OnCallPrioritiesResourceModel, result *apiclient.OnCallPriorities,
) {
	model.Priorities = make([]OnCallPriorityModel, len(result.Priorities))
	for i, p := range result.Priorities {
		model.Priorities[i] = OnCallPriorityModel{
			Level:    types.StringValue(p.Level),
			Label:    types.StringValue(p.Label),
			Pageable: types.BoolValue(p.Pageable),
		}
	}
}
