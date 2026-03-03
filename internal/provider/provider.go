package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ provider.Provider = &MoneatProvider{}

// MoneatProvider defines the provider implementation.
type MoneatProvider struct {
	version string
}

// MoneatProviderModel describes the provider data model.
type MoneatProviderModel struct {
	BaseURL types.String `tfsdk:"base_url"`
	Token   types.String `tfsdk:"token"`
}

// New returns a new provider factory function.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MoneatProvider{
			version: version,
		}
	}
}

func (p *MoneatProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "moneat"
	resp.Version = p.version
}

func (p *MoneatProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Moneat observability platform.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "Moneat API base URL. Can also be set with the MONEAT_BASE_URL environment variable. " +
					"Defaults to https://api.moneat.io.",
				Optional: true,
			},
			"token": schema.StringAttribute{
				Description: "Moneat API authentication token. Can also be set with the MONEAT_AUTH_TOKEN " +
					"environment variable.",
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *MoneatProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config MoneatProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve base URL
	baseURL := "https://api.moneat.io"
	if !config.BaseURL.IsNull() && !config.BaseURL.IsUnknown() {
		baseURL = config.BaseURL.ValueString()
	} else if v := os.Getenv("MONEAT_BASE_URL"); v != "" {
		baseURL = v
	}

	// Resolve token
	var token string
	if !config.Token.IsNull() && !config.Token.IsUnknown() {
		token = config.Token.ValueString()
	} else if v := os.Getenv("MONEAT_AUTH_TOKEN"); v != "" {
		token = v
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing API Token",
			"The provider requires a Moneat API token. Set the 'token' attribute in the provider block "+
				"or the MONEAT_AUTH_TOKEN environment variable.",
		)
		return
	}

	client := apiclient.NewClient(baseURL, token)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MoneatProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewUptimeMonitorResource,
		NewSystemAlertResource,
		NewStatusPageResource,
		NewDashboardResource,
		NewNotificationPreferencesResource,
		NewAlertNotificationChannelsResource,
		// Phase 2
		NewCustomDataSourceResource,
		NewLogIndexResource,
		NewSilencePeriodResource,
		NewDashboardAlertResource,
		NewDashboardFolderResource,
		NewStatusPageMonitorResource,
		NewStatusPageIncidentResource,
		NewStatusPageCustomDomainResource,
		NewOnCallScheduleResource,
		NewOnCallOverrideResource,
		NewEscalationPolicyResource,
		NewIncidentProviderResource,
		NewIncidentRoutingRuleResource,
		NewSyntheticTestResource,
		// Phase 3
		NewOrgMemberResource,
		NewOrgInvitationResource,
		NewAuthTokenResource,
		NewLogAPIKeyResource,
		NewAgentAPIKeyResource,
		NewSSOConfigResource,
		NewDebuggerProbeResource,
		NewOnCallPrioritiesResource,
		NewBusinessHoursResource,
		NewOnCallScheduleSlackUsergroupResource,
		NewProjectNotificationPreferencesResource,
	}
}

func (p *MoneatProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewProjectsDataSource,
		NewHostDataSource,
		NewHostsDataSource,
		NewUptimeMonitorDataSource,
		NewUptimeMonitorsDataSource,
		NewStatusPageDataSource,
		NewStatusPagesDataSource,
		NewDashboardDataSource,
		NewDashboardsDataSource,
		NewOnCallScheduleDataSource,
		NewEscalationPolicyDataSource,
		NewOrgMembersDataSource,
	}
}
