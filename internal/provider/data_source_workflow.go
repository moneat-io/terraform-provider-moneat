package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &WorkflowDataSource{}
var _ datasource.DataSource = &WorkflowsDataSource{}

type WorkflowDataSource struct {
	client *apiclient.Client
}

type WorkflowsDataSource struct {
	client *apiclient.Client
}

type WorkflowDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	TriggerName     types.String `tfsdk:"trigger_name"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Published       types.Bool   `tfsdk:"published"`
	ConditionsJSON  types.String `tfsdk:"conditions_json"`
	StepsJSON       types.String `tfsdk:"steps_json"`
	GraphJSON       types.String `tfsdk:"graph_json"`
	OnceForTemplate types.List   `tfsdk:"once_for_template"`
	Version         types.Int64  `tfsdk:"version"`
	SystemKey       types.String `tfsdk:"system_key"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	LastRunAt       types.String `tfsdk:"last_run_at"`
	RunCount        types.Int64  `tfsdk:"run_count"`
}

type WorkflowsDataSourceModel struct {
	Workflows []WorkflowDataSourceModel `tfsdk:"workflows"`
}

func NewWorkflowDataSource() datasource.DataSource {
	return &WorkflowDataSource{}
}

func NewWorkflowsDataSource() datasource.DataSource {
	return &WorkflowsDataSource{}
}

func (d *WorkflowDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (d *WorkflowsDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_workflows"
}

func (d *WorkflowDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Moneat workflow by ID.",
		Attributes:  workflowDataSourceAttributes(true),
	}
}

func (d *WorkflowsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list Moneat workflows.",
		Attributes: map[string]schema.Attribute{
			"workflows": schema.ListNestedAttribute{
				Description: "List of workflows.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: workflowDataSourceAttributes(false),
				},
			},
		},
	}
}

func (d *WorkflowDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	configureWorkflowDataSource(req.ProviderData, &d.client, resp)
}

func (d *WorkflowsDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	configureWorkflowDataSource(req.ProviderData, &d.client, resp)
}

func (d *WorkflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config WorkflowDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(config.ID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading workflow", err.Error())
		return
	}
	workflow, err := d.client.GetWorkflow(id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading workflow", err.Error())
		return
	}

	state := workflowDataSourceModel(ctx, &resp.Diagnostics, workflow)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *WorkflowsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	workflows, err := d.client.ListWorkflows()
	if err != nil {
		resp.Diagnostics.AddError("Error listing workflows", err.Error())
		return
	}

	state := WorkflowsDataSourceModel{
		Workflows: make([]WorkflowDataSourceModel, 0, len(workflows)),
	}
	for _, workflow := range workflows {
		model := workflowDataSourceModel(ctx, &resp.Diagnostics, &workflow)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Workflows = append(state.Workflows, model)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func configureWorkflowDataSource(
	providerData any,
	client **apiclient.Client,
	resp *datasource.ConfigureResponse,
) {
	if providerData == nil {
		return
	}
	configuredClient, ok := providerData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *apiclient.Client, got: %T", providerData),
		)
		return
	}
	*client = configuredClient
}

func workflowDataSourceAttributes(requireID bool) map[string]schema.Attribute {
	idAttribute := schema.StringAttribute{
		Description: "The unique identifier of the workflow.",
		Computed:    true,
	}
	if requireID {
		idAttribute = schema.StringAttribute{
			Description: "The unique identifier of the workflow.",
			Required:    true,
		}
	}
	return map[string]schema.Attribute{
		"id": idAttribute,
		"name": schema.StringAttribute{
			Description: "The workflow display name.",
			Computed:    true,
		},
		"trigger_name": schema.StringAttribute{
			Description: "The workflow trigger name.",
			Computed:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Whether the workflow is enabled.",
			Computed:    true,
		},
		"published": schema.BoolAttribute{
			Description: "Whether the latest workflow version is published.",
			Computed:    true,
		},
		"conditions_json": schema.StringAttribute{
			Description: "Workflow conditions as normalized JSON.",
			Computed:    true,
		},
		"steps_json": schema.StringAttribute{
			Description: "Workflow steps as normalized JSON.",
			Computed:    true,
		},
		"graph_json": schema.StringAttribute{
			Description: "Workflow graph as normalized JSON.",
			Computed:    true,
		},
		"once_for_template": schema.ListAttribute{
			Description: "Scope fields used to deduplicate workflow runs.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"version": schema.Int64Attribute{
			Description: "The current workflow version.",
			Computed:    true,
		},
		"system_key": schema.StringAttribute{
			Description: "System workflow key, when the workflow is system-managed.",
			Computed:    true,
		},
		"created_at": schema.StringAttribute{
			Description: "Workflow creation timestamp.",
			Computed:    true,
		},
		"updated_at": schema.StringAttribute{
			Description: "Workflow update timestamp.",
			Computed:    true,
		},
		"last_run_at": schema.StringAttribute{
			Description: "Timestamp of the most recent workflow run.",
			Computed:    true,
		},
		"run_count": schema.Int64Attribute{
			Description: "Number of workflow runs reported by the API.",
			Computed:    true,
		},
	}
}

func workflowDataSourceModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	workflow *apiclient.Workflow,
) WorkflowDataSourceModel {
	onceForTemplate, listDiags := types.ListValueFrom(ctx, types.StringType, workflow.OnceForTemplate)
	diags.Append(listDiags...)
	return WorkflowDataSourceModel{
		ID:              terraformID(workflow.ID),
		Name:            types.StringValue(workflow.Name),
		TriggerName:     types.StringValue(workflow.TriggerName),
		Enabled:         types.BoolValue(workflow.Enabled),
		Published:       types.BoolValue(workflow.Published),
		ConditionsJSON:  types.StringValue(rawMessageString(workflow.Conditions)),
		StepsJSON:       types.StringValue(rawMessageString(workflow.Steps)),
		GraphJSON:       types.StringValue(rawMessageString(workflow.Graph)),
		OnceForTemplate: onceForTemplate,
		Version:         types.Int64Value(int64(workflow.Version)),
		SystemKey:       types.StringValue(workflow.SystemKey),
		CreatedAt:       types.StringValue(workflow.CreatedAt),
		UpdatedAt:       types.StringValue(workflow.UpdatedAt),
		LastRunAt:       types.StringValue(workflow.LastRunAt),
		RunCount:        types.Int64Value(workflow.RunCount),
	}
}
