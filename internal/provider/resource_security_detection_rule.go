package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

var _ resource.Resource = &SecurityDetectionRuleResource{}
var _ resource.ResourceWithImportState = &SecurityDetectionRuleResource{}

type SecurityDetectionRuleResource struct {
	client *apiclient.Client
}

type SecurityDetectionRuleResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Source         types.String `tfsdk:"source"`
	Filter         types.String `tfsdk:"filter"`
	GroupBy        types.List   `tfsdk:"group_by"`
	WindowSeconds  types.Int64  `tfsdk:"window_seconds"`
	Type           types.String `tfsdk:"type"`
	ThresholdCount types.Int64  `tfsdk:"threshold_count"`
	Severity       types.String `tfsdk:"severity"`
	SignalTitle    types.String `tfsdk:"signal_title"`
	SignalMessage  types.String `tfsdk:"signal_message"`
	Suppressions   types.List   `tfsdk:"suppressions"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	Tags           types.List   `tfsdk:"tags"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func NewSecurityDetectionRuleResource() resource.Resource {
	return &SecurityDetectionRuleResource{}
}

func (r *SecurityDetectionRuleResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_security_detection_rule"
}

func (r *SecurityDetectionRuleResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat security detection rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the detection rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The detection rule name.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Human-readable rule description.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"source": schema.StringAttribute{
				Description: "Telemetry source. Currently logs.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("logs"),
			},
			"filter": schema.StringAttribute{
				Description: "Moneat log query filter expression.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"group_by": schema.ListAttribute{
				Description: "Fields used to group detection matches.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"window_seconds": schema.Int64Attribute{
				Description: "Evaluation window in seconds.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(300),
			},
			"type": schema.StringAttribute{
				Description: "Detection rule type: threshold, new_value, or rate_anomaly.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("threshold"),
			},
			"threshold_count": schema.Int64Attribute{
				Description: "Threshold count for threshold rules.",
				Optional:    true,
			},
			"severity": schema.StringAttribute{
				Description: "Signal severity: info, low, medium, high, or critical.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("medium"),
			},
			"signal_title": schema.StringAttribute{
				Description: "Rendered signal title.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"signal_message": schema.StringAttribute{
				Description: "Rendered signal message.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"suppressions": schema.ListAttribute{
				Description: "Suppression keys for the rule.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the rule is evaluated by the scheduler.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"tags": schema.ListAttribute{
				Description: "Rule tags.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "Rule creation timestamp.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Rule update timestamp.",
				Computed:    true,
			},
		},
	}
}

func (r *SecurityDetectionRuleResource) Configure(
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

func (r *SecurityDetectionRuleResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan SecurityDetectionRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq, ok := detectionRuleRequest(ctx, &resp.Diagnostics, plan)
	if !ok {
		return
	}
	rule, err := r.client.CreateDetectionRule(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating security detection rule", err.Error())
		return
	}

	mapDetectionRuleToState(ctx, &resp.Diagnostics, &plan, rule)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SecurityDetectionRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SecurityDetectionRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "security detection rule ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading security detection rule", err.Error())
		return
	}
	rule, err := r.client.GetDetectionRule(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading security detection rule", err.Error())
		return
	}

	mapDetectionRuleToState(ctx, &resp.Diagnostics, &state, rule)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SecurityDetectionRuleResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan SecurityDetectionRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(plan.ID, "security detection rule ID")
	if err != nil {
		resp.Diagnostics.AddError("Error updating security detection rule", err.Error())
		return
	}
	apiReq, ok := detectionRuleRequest(ctx, &resp.Diagnostics, plan)
	if !ok {
		return
	}
	rule, err := r.client.UpdateDetectionRule(id, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating security detection rule", err.Error())
		return
	}

	mapDetectionRuleToState(ctx, &resp.Diagnostics, &plan, rule)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SecurityDetectionRuleResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state SecurityDetectionRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseTerraformID(state.ID, "security detection rule ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting security detection rule", err.Error())
		return
	}
	err = r.client.DeleteDetectionRule(id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting security detection rule", err.Error())
		return
	}
}

func (r *SecurityDetectionRuleResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := parseTerraformID(types.StringValue(req.ID), "security detection rule ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing security detection rule", err.Error())
		return
	}
	rule, err := r.client.GetDetectionRule(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing security detection rule", err.Error())
		return
	}

	state := SecurityDetectionRuleResourceModel{}
	mapDetectionRuleToState(ctx, &resp.Diagnostics, &state, rule)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func detectionRuleRequest(
	ctx context.Context,
	diags *diag.Diagnostics,
	model SecurityDetectionRuleResourceModel,
) (apiclient.DetectionRuleRequest, bool) {
	groupBy, ok := stringList(ctx, diags, model.GroupBy)
	if !ok {
		return apiclient.DetectionRuleRequest{}, false
	}
	suppressions, ok := stringList(ctx, diags, model.Suppressions)
	if !ok {
		return apiclient.DetectionRuleRequest{}, false
	}
	tags, ok := stringList(ctx, diags, model.Tags)
	if !ok {
		return apiclient.DetectionRuleRequest{}, false
	}
	enabled := model.Enabled.ValueBool()
	return apiclient.DetectionRuleRequest{
		Name:           model.Name.ValueString(),
		Description:    model.Description.ValueString(),
		Source:         model.Source.ValueString(),
		Filter:         model.Filter.ValueString(),
		GroupBy:        groupBy,
		WindowSeconds:  int(model.WindowSeconds.ValueInt64()),
		Type:           model.Type.ValueString(),
		ThresholdCount: intPointerFromInt64(model.ThresholdCount),
		Severity:       model.Severity.ValueString(),
		SignalTitle:    model.SignalTitle.ValueString(),
		SignalMessage:  model.SignalMessage.ValueString(),
		Suppressions:   suppressions,
		Enabled:        &enabled,
		Tags:           tags,
	}, true
}

func mapDetectionRuleToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *SecurityDetectionRuleResourceModel,
	rule *apiclient.DetectionRule,
) {
	model.ID = terraformID(rule.ID)
	model.Name = types.StringValue(rule.Name)
	model.Description = types.StringValue(rule.Description)
	model.Source = types.StringValue(rule.Source)
	model.Filter = types.StringValue(rule.Filter)
	model.WindowSeconds = types.Int64Value(int64(rule.WindowSeconds))
	model.Type = types.StringValue(rule.Type)
	if rule.ThresholdCount == nil {
		model.ThresholdCount = types.Int64Null()
	} else {
		model.ThresholdCount = types.Int64Value(int64(*rule.ThresholdCount))
	}
	model.Severity = types.StringValue(rule.Severity)
	model.SignalTitle = types.StringValue(rule.SignalTitle)
	model.SignalMessage = types.StringValue(rule.SignalMessage)
	model.Enabled = types.BoolValue(rule.Enabled)
	model.CreatedAt = optionalString(rule.CreatedAt)
	model.UpdatedAt = optionalString(rule.UpdatedAt)

	groupBy, groupDiags := types.ListValueFrom(ctx, types.StringType, rule.GroupBy)
	diags.Append(groupDiags...)
	suppressions, suppressionDiags := types.ListValueFrom(ctx, types.StringType, rule.Suppressions)
	diags.Append(suppressionDiags...)
	tags, tagDiags := types.ListValueFrom(ctx, types.StringType, rule.Tags)
	diags.Append(tagDiags...)
	model.GroupBy = groupBy
	model.Suppressions = suppressions
	model.Tags = tags
}

func stringList(ctx context.Context, diags *diag.Diagnostics, value types.List) ([]string, bool) {
	var result []string
	if value.IsNull() || value.IsUnknown() {
		return []string{}, true
	}
	diags.Append(value.ElementsAs(ctx, &result, false)...)
	return result, !diags.HasError()
}
