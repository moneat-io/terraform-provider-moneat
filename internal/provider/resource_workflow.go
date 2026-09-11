package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &WorkflowResource{}
var _ resource.ResourceWithImportState = &WorkflowResource{}

type WorkflowResource struct {
	client *apiclient.Client
}

type WorkflowResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	TriggerName           types.String `tfsdk:"trigger_name"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	Published             types.Bool   `tfsdk:"published"`
	ConditionsJSON        types.String `tfsdk:"conditions_json"`
	StepsJSON             types.String `tfsdk:"steps_json"`
	GraphJSON             types.String `tfsdk:"graph_json"`
	OnceForTemplate       types.List   `tfsdk:"once_for_template"`
	RunOnceJSON           types.String `tfsdk:"run_once_json"`
	InputSchemaJSON       types.String `tfsdk:"input_schema_json"`
	TriggerNames          types.List   `tfsdk:"trigger_names"`
	SchedulesJSON         types.String `tfsdk:"schedules_json"`
	ExecutionIdentityJSON types.String `tfsdk:"execution_identity_json"`
	Version               types.Int64  `tfsdk:"version"`
}

func NewWorkflowResource() resource.Resource {
	return &WorkflowResource{}
}

func (r *WorkflowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (r *WorkflowResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat alert or automation workflow.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the workflow.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the workflow.",
				Required:    true,
			},
			"trigger_name": schema.StringAttribute{
				Description: "The workflow trigger name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the workflow is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"published": schema.BoolAttribute{
				Description: "Whether the latest workflow version is published and eligible for event triggers.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"conditions_json": schema.StringAttribute{
				Description: "Workflow conditions as JSON.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
			},
			"steps_json": schema.StringAttribute{
				Description: "Workflow steps as JSON.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
			},
			"graph_json": schema.StringAttribute{
				Description: "Workflow graph as JSON. Use this for branch, retry, control, approval, and connector nodes.",
				Optional:    true,
				Computed:    true,
			},
			"once_for_template": schema.ListAttribute{
				Description: "Scope fields used to deduplicate workflow runs.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default: listdefault.StaticValue(
					types.ListValueMust(types.StringType, []attr.Value{}),
				),
			},
			"run_once_json": schema.StringAttribute{
				Description: "Optional run-once/deduplication configuration as JSON.",
				Optional:    true,
				Computed:    true,
			},
			"input_schema_json": schema.StringAttribute{
				Description: "Workflow input schema as JSON.",
				Optional:    true,
				Computed:    true,
			},
			"trigger_names": schema.ListAttribute{
				Description: "Additional trigger names accepted by the workflow.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default: listdefault.StaticValue(
					types.ListValueMust(types.StringType, []attr.Value{}),
				),
			},
			"schedules_json": schema.StringAttribute{
				Description: "Durable workflow schedules as JSON.",
				Optional:    true,
				Computed:    true,
			},
			"execution_identity_json": schema.StringAttribute{
				Description: "Workflow execution identity as JSON.",
				Optional:    true,
				Computed:    true,
			},
			"version": schema.Int64Attribute{
				Description: "The current workflow version.",
				Computed:    true,
			},
		},
	}
}

func (r *WorkflowResource) Configure(
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

func (r *WorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkflowResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, normalizedConditions, normalizedSteps, normalizedGraph, ok := workflowPayload(ctx, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	apiReq := apiclient.CreateWorkflowRequest{
		Name:              plan.Name.ValueString(),
		TriggerName:       plan.TriggerName.ValueString(),
		Enabled:           plan.Enabled.ValueBool(),
		Conditions:        payload.Conditions,
		Steps:             payload.Steps,
		Graph:             payload.Graph,
		OnceForTemplate:   payload.OnceForTemplate,
		RunOnce:           payload.RunOnce,
		InputSchema:       payload.InputSchema,
		TriggerNames:      payload.TriggerNames,
		Schedules:         payload.Schedules,
		ExecutionIdentity: payload.ExecutionIdentity,
	}
	workflow, err := r.client.CreateWorkflow(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating workflow", err.Error())
		return
	}
	workflow, err = r.setWorkflowPublished(workflow.ID, plan.Published.ValueBool(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Error publishing workflow", err.Error())
		return
	}

	mapWorkflowToState(ctx, &resp.Diagnostics, &plan, workflow, normalizedConditions, normalizedSteps, normalizedGraph)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformUUID(state.ID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading workflow", err.Error())
		return
	}
	workflow, err := r.client.GetWorkflow(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading workflow", err.Error())
		return
	}

	mapWorkflowToState(ctx, &resp.Diagnostics, &state, workflow, "", "", "")
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkflowResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state WorkflowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	idValue := plan.ID
	if idValue.IsNull() || idValue.IsUnknown() || idValue.ValueString() == "" {
		idValue = state.ID
	}
	id, err := parseTerraformUUID(idValue, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Error updating workflow", err.Error())
		return
	}
	payload, normalizedConditions, normalizedSteps, normalizedGraph, ok := workflowPayload(ctx, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	enabled := plan.Enabled.ValueBool()
	apiReq := apiclient.UpdateWorkflowRequest{
		Name:              plan.Name.ValueString(),
		Enabled:           &enabled,
		Conditions:        payload.Conditions,
		Steps:             payload.Steps,
		Graph:             payload.Graph,
		OnceForTemplate:   payload.OnceForTemplate,
		ExpectedVersion:   workflowVersionPointer(state.Version),
		RunOnce:           payload.RunOnce,
		InputSchema:       payload.InputSchema,
		TriggerNames:      payload.TriggerNames,
		Schedules:         payload.Schedules,
		ExecutionIdentity: payload.ExecutionIdentity,
	}
	updated, err := r.client.UpdateWorkflow(id, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating workflow", err.Error())
		return
	}
	// Updating the declarative definition creates a new optimistic version. Use
	// that returned version for the subsequent publication transition; the
	// prior state version would always be stale after a successful update.
	updatedVersion := workflowVersionPointer(types.Int64Value(int64(updated.Version)))
	workflow, err := r.setWorkflowPublished(id, plan.Published.ValueBool(), updatedVersion)
	if err != nil {
		resp.Diagnostics.AddError("Error publishing workflow", err.Error())
		return
	}

	mapWorkflowToState(ctx, &resp.Diagnostics, &plan, workflow, normalizedConditions, normalizedSteps, normalizedGraph)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkflowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformUUID(state.ID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workflow", err.Error())
		return
	}
	err = r.client.DeleteWorkflow(id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workflow", err.Error())
		return
	}
}

func (r *WorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := parseTerraformUUID(types.StringValue(req.ID), "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow", err.Error())
		return
	}
	workflow, err := r.client.GetWorkflow(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow", err.Error())
		return
	}

	state := WorkflowResourceModel{}
	mapWorkflowToState(ctx, &resp.Diagnostics, &state, workflow, "", "", "")
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type workflowJSONPayload struct {
	Conditions        json.RawMessage
	Steps             json.RawMessage
	Graph             json.RawMessage
	OnceForTemplate   []string
	RunOnce           json.RawMessage
	InputSchema       json.RawMessage
	TriggerNames      []string
	Schedules         json.RawMessage
	ExecutionIdentity json.RawMessage
}

func workflowPayload(
	ctx context.Context,
	model WorkflowResourceModel,
	diags *diag.Diagnostics,
) (workflowJSONPayload, string, string, string, bool) {
	conditions, normalizedConditions, err := rawMessageFromJSONString(model.ConditionsJSON.ValueString())
	if err != nil {
		diags.AddError("Invalid workflow conditions JSON", err.Error())
		return workflowJSONPayload{}, "", "", "", false
	}
	steps, normalizedSteps, err := rawMessageFromJSONString(model.StepsJSON.ValueString())
	if err != nil {
		diags.AddError("Invalid workflow steps JSON", err.Error())
		return workflowJSONPayload{}, "", "", "", false
	}
	var graph json.RawMessage
	normalizedGraph := ""
	if !model.GraphJSON.IsNull() && !model.GraphJSON.IsUnknown() && model.GraphJSON.ValueString() != "" {
		graph, normalizedGraph, err = rawMessageFromJSONString(model.GraphJSON.ValueString())
		if err != nil {
			diags.AddError("Invalid workflow graph JSON", err.Error())
			return workflowJSONPayload{}, "", "", "", false
		}
	}
	var onceForTemplate []string
	diags.Append(model.OnceForTemplate.ElementsAs(ctx, &onceForTemplate, false)...)
	if diags.HasError() {
		return workflowJSONPayload{}, "", "", "", false
	}
	triggerNames := []string{}
	diags.Append(model.TriggerNames.ElementsAs(ctx, &triggerNames, false)...)
	if diags.HasError() {
		return workflowJSONPayload{}, "", "", "", false
	}
	runOnce, err := optionalRawMessage(model.RunOnceJSON)
	if err != nil {
		diags.AddError("Invalid workflow run_once JSON", err.Error())
		return workflowJSONPayload{}, "", "", "", false
	}
	inputSchema, err := optionalRawMessage(model.InputSchemaJSON)
	if err != nil {
		diags.AddError("Invalid workflow input schema JSON", err.Error())
		return workflowJSONPayload{}, "", "", "", false
	}
	schedules, err := optionalRawMessage(model.SchedulesJSON)
	if err != nil {
		diags.AddError("Invalid workflow schedules JSON", err.Error())
		return workflowJSONPayload{}, "", "", "", false
	}
	executionIdentity, err := optionalRawMessage(model.ExecutionIdentityJSON)
	if err != nil {
		diags.AddError("Invalid workflow execution identity JSON", err.Error())
		return workflowJSONPayload{}, "", "", "", false
	}
	return workflowJSONPayload{
		Conditions:        conditions,
		Steps:             steps,
		Graph:             graph,
		OnceForTemplate:   onceForTemplate,
		RunOnce:           runOnce,
		InputSchema:       inputSchema,
		TriggerNames:      triggerNames,
		Schedules:         schedules,
		ExecutionIdentity: executionIdentity,
	}, normalizedConditions, normalizedSteps, normalizedGraph, true
}

func mapWorkflowToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *WorkflowResourceModel,
	workflow *apiclient.Workflow,
	conditionsJSON string,
	stepsJSON string,
	graphJSON string,
) {
	model.ID = types.StringValue(workflow.ID)
	model.Name = types.StringValue(workflow.Name)
	model.TriggerName = types.StringValue(workflow.TriggerName)
	model.Enabled = types.BoolValue(workflow.Enabled)
	model.Published = types.BoolValue(workflow.Published)
	if conditionsJSON == "" {
		conditionsJSON = rawMessageString(workflow.Conditions)
	}
	if stepsJSON == "" {
		stepsJSON = rawMessageString(workflow.Steps)
	}
	if graphJSON == "" {
		graphJSON = rawMessageString(workflow.Graph)
	}
	model.ConditionsJSON = types.StringValue(conditionsJSON)
	model.StepsJSON = types.StringValue(stepsJSON)
	model.GraphJSON = types.StringValue(graphJSON)
	model.Version = types.Int64Value(int64(workflow.Version))

	onceForTemplate, listDiags := types.ListValueFrom(ctx, types.StringType, workflow.OnceForTemplate)
	diags.Append(listDiags...)
	if diags.HasError() {
		return
	}
	model.OnceForTemplate = onceForTemplate
	model.RunOnceJSON = optionalRawMessageState(workflow.RunOnce)
	model.InputSchemaJSON = optionalRawMessageState(workflow.InputSchema)
	triggerNames, triggerDiags := types.ListValueFrom(ctx, types.StringType, workflow.TriggerNames)
	diags.Append(triggerDiags...)
	model.TriggerNames = triggerNames
	model.SchedulesJSON = optionalRawMessageState(workflow.Schedules)
	model.ExecutionIdentityJSON = optionalRawMessageState(workflow.ExecutionIdentity)
}

func (r *WorkflowResource) setWorkflowPublished(id string, published bool, expectedVersion *int) (*apiclient.Workflow, error) {
	if published {
		return r.client.PublishWorkflow(id, expectedVersion)
	}
	return r.client.UnpublishWorkflow(id, expectedVersion)
}

func workflowVersionPointer(value types.Int64) *int {
	if value.IsNull() || value.IsUnknown() || value.ValueInt64() <= 0 {
		return nil
	}
	version := int(value.ValueInt64())
	return &version
}

func optionalRawMessage(value types.String) (json.RawMessage, error) {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return nil, nil
	}
	normalized, err := normalizeJSONString(value.ValueString())
	if err != nil {
		return nil, err
	}
	return json.RawMessage(normalized), nil
}

func optionalRawMessageState(value json.RawMessage) types.String {
	if len(value) == 0 {
		return types.StringNull()
	}
	return types.StringValue(rawMessageString(value))
}
