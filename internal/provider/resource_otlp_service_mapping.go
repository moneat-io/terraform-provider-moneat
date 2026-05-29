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

var _ resource.Resource = &OtlpServiceMappingResource{}
var _ resource.ResourceWithImportState = &OtlpServiceMappingResource{}

type OtlpServiceMappingResource struct {
	client *apiclient.Client
}

type OtlpServiceMappingResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ServiceNamespace  types.String `tfsdk:"service_namespace"`
	ServiceName       types.String `tfsdk:"service_name"`
	ProjectID         types.Int64  `tfsdk:"project_id"`
	ProjectResourceID types.String `tfsdk:"project_resource_id"`
	ProjectName       types.String `tfsdk:"project_name"`
}

func NewOtlpServiceMappingResource() resource.Resource {
	return &OtlpServiceMappingResource{}
}

func (r *OtlpServiceMappingResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_otlp_service_mapping"
}

func (r *OtlpServiceMappingResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat telemetry service-to-project mapping.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the OTLP service mapping.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_namespace": schema.StringAttribute{
				Description: "The OpenTelemetry service.namespace value.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_name": schema.StringAttribute{
				Description: "The OpenTelemetry service.name value.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.Int64Attribute{
				Description: "The numeric Moneat project ID to route this service to.",
				Optional:    true,
				Computed:    true,
			},
			"project_resource_id": schema.StringAttribute{
				Description: "The stable Moneat project resource ID to route this service to.",
				Optional:    true,
				Computed:    true,
			},
			"project_name": schema.StringAttribute{
				Description: "The routed project display name.",
				Computed:    true,
			},
		},
	}
}

func (r *OtlpServiceMappingResource) Configure(
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

func (r *OtlpServiceMappingResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan OtlpServiceMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := r.client.UpsertOtlpServiceMapping(otlpServiceMappingRequest(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating OTLP service mapping", err.Error())
		return
	}

	mapOtlpServiceMappingToState(&plan, mapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OtlpServiceMappingResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state OtlpServiceMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "OTLP service mapping ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading OTLP service mapping", err.Error())
		return
	}
	mapping, err := r.client.FindOtlpServiceMappingFromObserved(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
		resp.Diagnostics.AddError("Error reading OTLP service mapping", err.Error())
		return
	}

	mapOtlpServiceMappingToState(&state, mapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OtlpServiceMappingResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan OtlpServiceMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := r.client.UpsertOtlpServiceMapping(otlpServiceMappingRequest(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating OTLP service mapping", err.Error())
		return
	}

	mapOtlpServiceMappingToState(&plan, mapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OtlpServiceMappingResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state OtlpServiceMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "OTLP service mapping ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting OTLP service mapping", err.Error())
		return
	}
	err = r.client.DeleteOtlpServiceMapping(id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting OTLP service mapping", err.Error())
		return
	}
}

func (r *OtlpServiceMappingResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := parseTerraformID(types.StringValue(req.ID), "OTLP service mapping ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing OTLP service mapping", err.Error())
		return
	}
	mapping, err := r.client.FindOtlpServiceMappingFromObserved(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing OTLP service mapping", err.Error())
		return
	}

	state := OtlpServiceMappingResourceModel{}
	mapOtlpServiceMappingToState(&state, mapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func otlpServiceMappingRequest(model OtlpServiceMappingResourceModel) apiclient.CreateOtlpServiceMappingRequest {
	req := apiclient.CreateOtlpServiceMappingRequest{
		ServiceName: model.ServiceName.ValueString(),
	}
	if !model.ServiceNamespace.IsNull() && !model.ServiceNamespace.IsUnknown() {
		req.ServiceNamespace = model.ServiceNamespace.ValueString()
	}
	if !model.ProjectID.IsNull() && !model.ProjectID.IsUnknown() {
		req.ProjectID = model.ProjectID.ValueInt64()
	}
	if !model.ProjectResourceID.IsNull() && !model.ProjectResourceID.IsUnknown() {
		req.ProjectResourceID = model.ProjectResourceID.ValueString()
	}
	return req
}

func mapOtlpServiceMappingToState(
	model *OtlpServiceMappingResourceModel,
	mapping *apiclient.OtlpServiceMapping,
) {
	model.ID = terraformID(mapping.ID)
	model.ServiceNamespace = types.StringValue(mapping.ServiceNamespace)
	model.ServiceName = types.StringValue(mapping.ServiceName)
	model.ProjectID = types.Int64Value(mapping.ProjectID)
	model.ProjectResourceID = types.StringValue(mapping.ProjectResourceID)
	model.ProjectName = types.StringValue(mapping.ProjectName)
}
