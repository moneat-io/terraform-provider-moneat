package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &ProjectNotificationPreferencesResource{}

type ProjectNotificationPreferencesResource struct {
	client *apiclient.Client
}

type ProjectNotificationPreferencesResourceModel struct {
	ProjectID             types.String `tfsdk:"project_id"`
	IssueAlerts           types.Bool   `tfsdk:"issue_alerts"`
	ErrorAlerts           types.Bool   `tfsdk:"error_alerts"`
	WeeklySummary         types.Bool   `tfsdk:"weekly_summary"`
	AlertFrequencyMinutes types.Int64  `tfsdk:"alert_frequency_minutes"`
}

func NewProjectNotificationPreferencesResource() resource.Resource {
	return &ProjectNotificationPreferencesResource{}
}

func (r *ProjectNotificationPreferencesResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_project_notification_preferences"
}

func (r *ProjectNotificationPreferencesResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages per-project notification preference overrides in Moneat.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Description: "The project ID to configure notification preferences for.",
				Required:    true,
			},
			"issue_alerts": schema.BoolAttribute{
				Description: "Whether to receive issue alert notifications for this project.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"error_alerts": schema.BoolAttribute{
				Description: "Whether to receive error alert notifications for this project.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"weekly_summary": schema.BoolAttribute{
				Description: "Whether to receive weekly summary emails for this project.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"alert_frequency_minutes": schema.Int64Attribute{
				Description: "Minimum minutes between alert notifications for this project.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(15),
			},
		},
	}
}

func (r *ProjectNotificationPreferencesResource) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
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

func (r *ProjectNotificationPreferencesResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan ProjectNotificationPreferencesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateProjectNotificationPreferencesRequest{
		IssueAlerts:           plan.IssueAlerts.ValueBool(),
		ErrorAlerts:           plan.ErrorAlerts.ValueBool(),
		WeeklySummary:         plan.WeeklySummary.ValueBool(),
		AlertFrequencyMinutes: plan.AlertFrequencyMinutes.ValueInt64(),
	}

	prefs, err := r.client.UpdateProjectNotificationPreferences(plan.ProjectID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error setting project notification preferences", err.Error())
		return
	}

	plan.ProjectID = types.StringValue(prefs.ProjectID)
	plan.IssueAlerts = types.BoolValue(prefs.IssueAlerts)
	plan.ErrorAlerts = types.BoolValue(prefs.ErrorAlerts)
	plan.WeeklySummary = types.BoolValue(prefs.WeeklySummary)
	plan.AlertFrequencyMinutes = types.Int64Value(prefs.AlertFrequencyMinutes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectNotificationPreferencesResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state ProjectNotificationPreferencesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	prefs, err := r.client.GetProjectNotificationPreferences(state.ProjectID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading project notification preferences", err.Error())
		return
	}

	state.IssueAlerts = types.BoolValue(prefs.IssueAlerts)
	state.ErrorAlerts = types.BoolValue(prefs.ErrorAlerts)
	state.WeeklySummary = types.BoolValue(prefs.WeeklySummary)
	state.AlertFrequencyMinutes = types.Int64Value(prefs.AlertFrequencyMinutes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProjectNotificationPreferencesResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan ProjectNotificationPreferencesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateProjectNotificationPreferencesRequest{
		IssueAlerts:           plan.IssueAlerts.ValueBool(),
		ErrorAlerts:           plan.ErrorAlerts.ValueBool(),
		WeeklySummary:         plan.WeeklySummary.ValueBool(),
		AlertFrequencyMinutes: plan.AlertFrequencyMinutes.ValueInt64(),
	}

	prefs, err := r.client.UpdateProjectNotificationPreferences(plan.ProjectID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project notification preferences", err.Error())
		return
	}

	plan.ProjectID = types.StringValue(prefs.ProjectID)
	plan.IssueAlerts = types.BoolValue(prefs.IssueAlerts)
	plan.ErrorAlerts = types.BoolValue(prefs.ErrorAlerts)
	plan.WeeklySummary = types.BoolValue(prefs.WeeklySummary)
	plan.AlertFrequencyMinutes = types.Int64Value(prefs.AlertFrequencyMinutes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectNotificationPreferencesResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state ProjectNotificationPreferencesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProjectNotificationPreferences(state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting project notification preferences", err.Error())
		return
	}
}
