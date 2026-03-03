package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &AlertNotificationChannelsResource{}

// AlertNotificationChannelsResource defines the resource implementation.
type AlertNotificationChannelsResource struct {
	client *apiclient.Client
}

// AlertNotificationChannelsResourceModel describes the resource data model.
type AlertNotificationChannelsResourceModel struct {
	AlertSource    types.String `tfsdk:"alert_source"`
	EmailEnabled   types.Bool   `tfsdk:"email_enabled"`
	SlackEnabled   types.Bool   `tfsdk:"slack_enabled"`
	DiscordEnabled types.Bool   `tfsdk:"discord_enabled"`
}

// NewAlertNotificationChannelsResource returns a new resource factory function.
func NewAlertNotificationChannelsResource() resource.Resource {
	return &AlertNotificationChannelsResource{}
}

func (r *AlertNotificationChannelsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_notification_channels"
}

func (r *AlertNotificationChannelsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Moneat alert notification channel preferences per alert source.",
		Attributes: map[string]schema.Attribute{
			"alert_source": schema.StringAttribute{
				Description: "The alert source type (e.g., system_alert, uptime_monitor, error_alert).",
				Required:    true,
			},
			"email_enabled": schema.BoolAttribute{
				Description: "Whether email notifications are enabled for this source.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"slack_enabled": schema.BoolAttribute{
				Description: "Whether Slack notifications are enabled for this source.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"discord_enabled": schema.BoolAttribute{
				Description: "Whether Discord notifications are enabled for this source.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *AlertNotificationChannelsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlertNotificationChannelsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlertNotificationChannelsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateAlertNotificationChannelsRequest{
		AlertSource:    plan.AlertSource.ValueString(),
		EmailEnabled:   plan.EmailEnabled.ValueBool(),
		SlackEnabled:   plan.SlackEnabled.ValueBool(),
		DiscordEnabled: plan.DiscordEnabled.ValueBool(),
	}

	channels, err := r.client.UpdateAlertNotificationChannels(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error setting alert notification channels", err.Error())
		return
	}

	plan.AlertSource = types.StringValue(channels.AlertSource)
	plan.EmailEnabled = types.BoolValue(channels.EmailEnabled)
	plan.SlackEnabled = types.BoolValue(channels.SlackEnabled)
	plan.DiscordEnabled = types.BoolValue(channels.DiscordEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AlertNotificationChannelsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlertNotificationChannelsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	allChannels, err := r.client.GetAlertNotificationChannels()
	if err != nil {
		resp.Diagnostics.AddError("Error reading alert notification channels", err.Error())
		return
	}

	alertSource := state.AlertSource.ValueString()
	for _, ch := range allChannels {
		if ch.AlertSource == alertSource {
			state.EmailEnabled = types.BoolValue(ch.EmailEnabled)
			state.SlackEnabled = types.BoolValue(ch.SlackEnabled)
			state.DiscordEnabled = types.BoolValue(ch.DiscordEnabled)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// Not found - remove from state
	resp.State.RemoveResource(ctx)
}

func (r *AlertNotificationChannelsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AlertNotificationChannelsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateAlertNotificationChannelsRequest{
		AlertSource:    plan.AlertSource.ValueString(),
		EmailEnabled:   plan.EmailEnabled.ValueBool(),
		SlackEnabled:   plan.SlackEnabled.ValueBool(),
		DiscordEnabled: plan.DiscordEnabled.ValueBool(),
	}

	channels, err := r.client.UpdateAlertNotificationChannels(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating alert notification channels", err.Error())
		return
	}

	plan.AlertSource = types.StringValue(channels.AlertSource)
	plan.EmailEnabled = types.BoolValue(channels.EmailEnabled)
	plan.SlackEnabled = types.BoolValue(channels.SlackEnabled)
	plan.DiscordEnabled = types.BoolValue(channels.DiscordEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AlertNotificationChannelsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlertNotificationChannelsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to defaults on delete
	apiReq := apiclient.UpdateAlertNotificationChannelsRequest{
		AlertSource:    state.AlertSource.ValueString(),
		EmailEnabled:   true,
		SlackEnabled:   false,
		DiscordEnabled: false,
	}

	_, err := r.client.UpdateAlertNotificationChannels(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error resetting alert notification channels", err.Error())
		return
	}
}
