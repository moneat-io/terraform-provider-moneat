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

var _ resource.Resource = &SSOConfigResource{}

type SSOConfigResource struct {
	client *apiclient.Client
}

type SSOConfigResourceModel struct {
	Provider    types.String `tfsdk:"provider_name"`
	EntityID    types.String `tfsdk:"entity_id"`
	SSOURL      types.String `tfsdk:"sso_url"`
	Certificate types.String `tfsdk:"certificate"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

func NewSSOConfigResource() resource.Resource {
	return &SSOConfigResource{}
}

func (r *SSOConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_config"
}

func (r *SSOConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the Moneat SSO configuration. This is a singleton resource.",
		Attributes: map[string]schema.Attribute{
			"provider_name": schema.StringAttribute{
				Description: "The SSO provider (e.g., saml, oidc).",
				Required:    true,
			},
			"entity_id": schema.StringAttribute{
				Description: "The entity ID for the SSO provider.",
				Required:    true,
			},
			"sso_url": schema.StringAttribute{
				Description: "The SSO login URL.",
				Required:    true,
			},
			"certificate": schema.StringAttribute{
				Description: "The SSO certificate.",
				Required:    true,
				Sensitive:   true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether SSO is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *SSOConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SSOConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSOConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateSSOConfigRequest{
		Provider:    plan.Provider.ValueString(),
		EntityID:    plan.EntityID.ValueString(),
		SSOURL:      plan.SSOURL.ValueString(),
		Certificate: plan.Certificate.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
	}

	config, err := r.client.UpdateSSOConfig(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error setting SSO configuration", err.Error())
		return
	}

	plan.Provider = types.StringValue(config.Provider)
	plan.EntityID = types.StringValue(config.EntityID)
	plan.SSOURL = types.StringValue(config.SSOURL)
	plan.Certificate = types.StringValue(config.Certificate)
	plan.Enabled = types.BoolValue(config.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSOConfigResource) Read(ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse) {
	config, err := r.client.GetSSOConfig()
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSO configuration", err.Error())
		return
	}

	state := SSOConfigResourceModel{
		Provider:    types.StringValue(config.Provider),
		EntityID:    types.StringValue(config.EntityID),
		SSOURL:      types.StringValue(config.SSOURL),
		Certificate: types.StringValue(config.Certificate),
		Enabled:     types.BoolValue(config.Enabled),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SSOConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SSOConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateSSOConfigRequest{
		Provider:    plan.Provider.ValueString(),
		EntityID:    plan.EntityID.ValueString(),
		SSOURL:      plan.SSOURL.ValueString(),
		Certificate: plan.Certificate.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
	}

	config, err := r.client.UpdateSSOConfig(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating SSO configuration", err.Error())
		return
	}

	plan.Provider = types.StringValue(config.Provider)
	plan.EntityID = types.StringValue(config.EntityID)
	plan.SSOURL = types.StringValue(config.SSOURL)
	plan.Certificate = types.StringValue(config.Certificate)
	plan.Enabled = types.BoolValue(config.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSOConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton resource - disable SSO on delete
	apiReq := apiclient.UpdateSSOConfigRequest{
		Provider:    "",
		EntityID:    "",
		SSOURL:      "",
		Certificate: "",
		Enabled:     false,
	}

	_, err := r.client.UpdateSSOConfig(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error resetting SSO configuration", err.Error())
		return
	}
}
