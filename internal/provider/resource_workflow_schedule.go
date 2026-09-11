package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &WorkflowScheduleResource{}
var _ resource.ResourceWithImportState = &WorkflowScheduleResource{}

type WorkflowScheduleResource struct {
	client *apiclient.Client
}

type WorkflowScheduleResourceModel struct {
	ID              types.String `tfsdk:"id"`
	WorkflowID      types.String `tfsdk:"workflow_id"`
	ScheduleID      types.String `tfsdk:"schedule_id"`
	Cron            types.String `tfsdk:"cron"`
	Timezone        types.String `tfsdk:"timezone"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	MissedRunPolicy types.String `tfsdk:"missed_run_policy"`
	OverlapPolicy   types.String `tfsdk:"overlap_policy"`
	MaxConcurrency  types.Int64  `tfsdk:"max_concurrency"`
	WorkflowVersion types.Int64  `tfsdk:"workflow_version"`
}

type workflowScheduleState struct {
	ID              string `json:"id"`
	Cron            string `json:"cron"`
	Timezone        string `json:"timezone"`
	Enabled         bool   `json:"enabled"`
	MissedRunPolicy string `json:"missed_run_policy"`
	OverlapPolicy   string `json:"overlap_policy"`
	MaxConcurrency  int    `json:"max_concurrency"`
}

func NewWorkflowScheduleResource() resource.Resource {
	return &WorkflowScheduleResource{}
}

func (r *WorkflowScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_schedule"
}

func (r *WorkflowScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages one durable schedule attached to a UUID-addressed workflow. Use depends_on to serialize schedules when explicit ordering is required.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"workflow_id": schema.StringAttribute{
				Description:   "UUID workflow resource ID.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"schedule_id": schema.StringAttribute{
				Description:   "Stable schedule identifier within the workflow.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"cron":              schema.StringAttribute{Required: true},
			"timezone":          schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("UTC")},
			"enabled":           schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true)},
			"missed_run_policy": schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("skip")},
			"overlap_policy":    schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("reject")},
			"max_concurrency":   schema.Int64Attribute{Optional: true, Computed: true, Default: int64default.StaticInt64(1)},
			"workflow_version":  schema.Int64Attribute{Computed: true},
		},
	}
}

func (r *WorkflowScheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *apiclient.Client, got: %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *WorkflowScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkflowScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(plan.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	workflow, err := r.client.GetWorkflow(workflowID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading workflow", err.Error())
		return
	}
	desired := scheduleFromModel(plan)
	existing, found, err := findSchedule(workflow.Schedules, desired.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow schedules JSON", err.Error())
		return
	}
	if found {
		if existing != desired {
			resp.Diagnostics.AddError(
				"Workflow schedule already exists",
				fmt.Sprintf("A different schedule with ID %q is already attached to workflow %s.", desired.ID, workflowID),
			)
			return
		}
		mapScheduleToState(&plan, workflowID, existing, workflow.Version)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}
	updated, err := r.replaceSchedule(workflowID, desired)
	if err != nil {
		resp.Diagnostics.AddError("Error saving workflow schedule", err.Error())
		return
	}
	mapScheduleToState(&plan, workflowID, desired, updated.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(state.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	workflow, err := r.client.GetWorkflow(workflowID)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading workflow schedule", err.Error())
		return
	}
	schedule, found, err := findSchedule(workflow.Schedules, state.ScheduleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow schedules JSON", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	mapScheduleToState(&state, workflowID, schedule, workflow.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkflowScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(plan.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	updated, err := r.replaceSchedule(workflowID, scheduleFromModel(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error saving workflow schedule", err.Error())
		return
	}
	mapScheduleToState(&plan, workflowID, scheduleFromModel(plan), updated.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkflowScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflowID, err := parseTerraformUUID(state.WorkflowID, "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	_, err = updateWorkflowWithRetry(r.client, workflowID, func(workflow *apiclient.Workflow) (apiclient.UpdateWorkflowRequest, error) {
		schedules, err := decodeSchedules(workflow.Schedules)
		if err != nil {
			return apiclient.UpdateWorkflowRequest{}, fmt.Errorf("invalid workflow schedules JSON: %w", err)
		}
		filtered := make([]workflowScheduleState, 0, len(schedules))
		for _, schedule := range schedules {
			if schedule.ID != state.ScheduleID.ValueString() {
				filtered = append(filtered, schedule)
			}
		}
		raw, err := json.Marshal(filtered)
		if err != nil {
			return apiclient.UpdateWorkflowRequest{}, err
		}
		request := workflowUpdateRequest(workflow)
		request.Schedules = raw
		return request, nil
	})
	if err != nil {
		if apiclient.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting workflow schedule", err.Error())
		return
	}
}

func (r *WorkflowScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid workflow schedule ID", "Import IDs must be <workflow UUID>/<schedule ID>")
		return
	}
	workflowID, err := parseTerraformUUID(types.StringValue(parts[0]), "workflow ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid workflow ID", err.Error())
		return
	}
	workflow, err := r.client.GetWorkflow(workflowID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow schedule", err.Error())
		return
	}
	schedule, found, err := findSchedule(workflow.Schedules, parts[1])
	if err != nil || !found {
		if err != nil {
			resp.Diagnostics.AddError("Invalid workflow schedules JSON", err.Error())
		} else {
			resp.Diagnostics.AddError("Workflow schedule not found", parts[1])
		}
		return
	}
	state := WorkflowScheduleResourceModel{}
	mapScheduleToState(&state, workflowID, schedule, workflow.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func scheduleFromModel(model WorkflowScheduleResourceModel) workflowScheduleState {
	return workflowScheduleState{
		ID:              model.ScheduleID.ValueString(),
		Cron:            model.Cron.ValueString(),
		Timezone:        model.Timezone.ValueString(),
		Enabled:         model.Enabled.ValueBool(),
		MissedRunPolicy: model.MissedRunPolicy.ValueString(),
		OverlapPolicy:   model.OverlapPolicy.ValueString(),
		MaxConcurrency:  int(model.MaxConcurrency.ValueInt64()),
	}
}

func (r *WorkflowScheduleResource) replaceSchedule(
	workflowID string,
	replacement workflowScheduleState,
) (*apiclient.Workflow, error) {
	return updateWorkflowWithRetry(r.client, workflowID, func(workflow *apiclient.Workflow) (apiclient.UpdateWorkflowRequest, error) {
		schedules, err := decodeSchedules(workflow.Schedules)
		if err != nil {
			return apiclient.UpdateWorkflowRequest{}, err
		}
		found := false
		for index := range schedules {
			if schedules[index].ID == replacement.ID {
				schedules[index] = replacement
				found = true
			}
		}
		if !found {
			schedules = append(schedules, replacement)
		}
		raw, err := json.Marshal(schedules)
		if err != nil {
			return apiclient.UpdateWorkflowRequest{}, err
		}
		request := workflowUpdateRequest(workflow)
		request.Schedules = raw
		return request, nil
	})
}

func decodeSchedules(raw json.RawMessage) ([]workflowScheduleState, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return []workflowScheduleState{}, nil
	}
	var schedules []workflowScheduleState
	if err := json.Unmarshal(raw, &schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}

func findSchedule(raw json.RawMessage, id string) (workflowScheduleState, bool, error) {
	schedules, err := decodeSchedules(raw)
	if err != nil {
		return workflowScheduleState{}, false, err
	}
	for _, schedule := range schedules {
		if schedule.ID == id {
			return schedule, true, nil
		}
	}
	return workflowScheduleState{}, false, nil
}

func mapScheduleToState(model *WorkflowScheduleResourceModel, workflowID string, schedule workflowScheduleState, version int) {
	model.ID = types.StringValue(workflowID + "/" + schedule.ID)
	model.WorkflowID = types.StringValue(workflowID)
	model.ScheduleID = types.StringValue(schedule.ID)
	model.Cron = types.StringValue(schedule.Cron)
	model.Timezone = types.StringValue(schedule.Timezone)
	model.Enabled = types.BoolValue(schedule.Enabled)
	model.MissedRunPolicy = types.StringValue(schedule.MissedRunPolicy)
	model.OverlapPolicy = types.StringValue(schedule.OverlapPolicy)
	model.MaxConcurrency = types.Int64Value(int64(schedule.MaxConcurrency))
	model.WorkflowVersion = types.Int64Value(int64(version))
}
