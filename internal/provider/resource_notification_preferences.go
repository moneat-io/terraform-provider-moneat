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

var _ resource.Resource = &NotificationPreferencesResource{}

// NotificationPreferencesResource defines the resource implementation.
type NotificationPreferencesResource struct {
	client *apiclient.Client
}

// NotificationPreferencesResourceModel describes the resource data model.
type NotificationPreferencesResourceModel struct {
	IssueAlerts           types.Bool  `tfsdk:"issue_alerts"`
	ErrorAlerts           types.Bool  `tfsdk:"error_alerts"`
	WeeklySummary         types.Bool  `tfsdk:"weekly_summary"`
	AlertFrequencyMinutes types.Int64 `tfsdk:"alert_frequency_minutes"`
}

// NewNotificationPreferencesResource returns a new resource factory function.
func NewNotificationPreferencesResource() resource.Resource {
	return &NotificationPreferencesResource{}
}

func (r *NotificationPreferencesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_preferences"
}

func (r *NotificationPreferencesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Moneat global notification preferences. This is a singleton resource.",
		Attributes: map[string]schema.Attribute{
			"issue_alerts": schema.BoolAttribute{
				Description: "Whether to receive issue alert notifications.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"error_alerts": schema.BoolAttribute{
				Description: "Whether to receive error alert notifications.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"weekly_summary": schema.BoolAttribute{
				Description: "Whether to receive weekly summary emails.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"alert_frequency_minutes": schema.Int64Attribute{
				Description: "Minimum minutes between alert notifications.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(15),
			},
		},
	}
}

func (r *NotificationPreferencesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationPreferencesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NotificationPreferencesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateNotificationPreferencesRequest{
		IssueAlerts:           plan.IssueAlerts.ValueBool(),
		ErrorAlerts:           plan.ErrorAlerts.ValueBool(),
		WeeklySummary:         plan.WeeklySummary.ValueBool(),
		AlertFrequencyMinutes: plan.AlertFrequencyMinutes.ValueInt64(),
	}

	prefs, err := r.client.UpdateNotificationPreferences(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error setting notification preferences", err.Error())
		return
	}

	plan.IssueAlerts = types.BoolValue(prefs.IssueAlerts)
	plan.ErrorAlerts = types.BoolValue(prefs.ErrorAlerts)
	plan.WeeklySummary = types.BoolValue(prefs.WeeklySummary)
	plan.AlertFrequencyMinutes = types.Int64Value(prefs.AlertFrequencyMinutes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NotificationPreferencesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	prefs, err := r.client.GetNotificationPreferences()
	if err != nil {
		resp.Diagnostics.AddError("Error reading notification preferences", err.Error())
		return
	}

	state := NotificationPreferencesResourceModel{
		IssueAlerts:           types.BoolValue(prefs.IssueAlerts),
		ErrorAlerts:           types.BoolValue(prefs.ErrorAlerts),
		WeeklySummary:         types.BoolValue(prefs.WeeklySummary),
		AlertFrequencyMinutes: types.Int64Value(prefs.AlertFrequencyMinutes),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NotificationPreferencesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NotificationPreferencesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateNotificationPreferencesRequest{
		IssueAlerts:           plan.IssueAlerts.ValueBool(),
		ErrorAlerts:           plan.ErrorAlerts.ValueBool(),
		WeeklySummary:         plan.WeeklySummary.ValueBool(),
		AlertFrequencyMinutes: plan.AlertFrequencyMinutes.ValueInt64(),
	}

	prefs, err := r.client.UpdateNotificationPreferences(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating notification preferences", err.Error())
		return
	}

	plan.IssueAlerts = types.BoolValue(prefs.IssueAlerts)
	plan.ErrorAlerts = types.BoolValue(prefs.ErrorAlerts)
	plan.WeeklySummary = types.BoolValue(prefs.WeeklySummary)
	plan.AlertFrequencyMinutes = types.Int64Value(prefs.AlertFrequencyMinutes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NotificationPreferencesResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton resource - reset to defaults on delete
	apiReq := apiclient.UpdateNotificationPreferencesRequest{
		IssueAlerts:           true,
		ErrorAlerts:           true,
		WeeklySummary:         true,
		AlertFrequencyMinutes: 15,
	}

	_, err := r.client.UpdateNotificationPreferences(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error resetting notification preferences", err.Error())
		return
	}
}
