package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &ConnectorInstallationResource{}
var _ resource.ResourceWithImportState = &ConnectorInstallationResource{}

type ConnectorInstallationResource struct {
	client *apiclient.Client
}

type ConnectorInstallationResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	ProviderID          types.String `tfsdk:"provider_id"`
	AuthProfileID       types.String `tfsdk:"auth_profile_id"`
	UseID               types.String `tfsdk:"use_id"`
	Name                types.String `tfsdk:"name"`
	ExternalAccountJSON types.String `tfsdk:"external_account_json"`
	Secret              types.String `tfsdk:"secret"`
	IdentifierTags      types.Map    `tfsdk:"identifier_tags"`
	CredentialType      types.String `tfsdk:"credential_type"`
	Status              types.String `tfsdk:"status"`
	StatusReason        types.String `tfsdk:"status_reason"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	APISecretLastFour   types.String `tfsdk:"api_secret_last_four"`
	SecretVersion       types.Int64  `tfsdk:"secret_version"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
}

func NewConnectorInstallationResource() resource.Resource {
	return &ConnectorInstallationResource{}
}

func (r *ConnectorInstallationResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_connector_installation"
}

func (r *ConnectorInstallationResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat connector installation. Secrets are accepted only by Terraform and stored sensitively.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "UUID resource ID of the connector installation.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"provider_id": schema.StringAttribute{
				Description: "Connector provider catalog ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auth_profile_id": schema.StringAttribute{
				Description: "Provider authentication profile catalog ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"use_id": schema.StringAttribute{
				Description: "Connector use served by this installation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Display name of the installation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"external_account_json": schema.StringAttribute{
				Description: "Provider-side account selector as JSON.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secret": schema.StringAttribute{
				Description: "Provider credential. Terraform state is sensitive; the API never returns plaintext.",
				Optional:    true,
				Sensitive:   true,
			},
			"identifier_tags": schema.MapAttribute{
				Description: "Tags used by connection groups for deterministic selection.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"credential_type":      schema.StringAttribute{Computed: true},
			"status":               schema.StringAttribute{Computed: true},
			"status_reason":        schema.StringAttribute{Computed: true},
			"enabled":              schema.BoolAttribute{Computed: true},
			"api_secret_last_four": schema.StringAttribute{Computed: true},
			"secret_version":       schema.Int64Attribute{Computed: true},
			"created_at":           schema.StringAttribute{Computed: true},
			"updated_at":           schema.StringAttribute{Computed: true},
		},
	}
}

func (r *ConnectorInstallationResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
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

func (r *ConnectorInstallationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan ConnectorInstallationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	secret := plan.Secret.ValueString()
	if secret == "" {
		resp.Diagnostics.AddError("Missing connector secret", "secret is required when creating a connector installation")
		return
	}
	externalAccountJSON, err := optionalJSON(plan.ExternalAccountJSON, "external_account_json")
	if err != nil {
		resp.Diagnostics.AddError("Invalid external account JSON", err.Error())
		return
	}
	var externalAccount apiclient.ConnectorExternalAccount
	if len(externalAccountJSON) > 0 {
		if err := json.Unmarshal(externalAccountJSON, &externalAccount); err != nil {
			resp.Diagnostics.AddError("Invalid external account JSON", err.Error())
			return
		}
	}
	tags, ok := connectorTags(ctx, &resp.Diagnostics, plan.IdentifierTags)
	if !ok {
		return
	}
	installation, err := r.client.CreateConnectorInstallation(apiclient.CreateConnectorInstallationRequest{
		ProviderID:      plan.ProviderID.ValueString(),
		AuthProfileID:   plan.AuthProfileID.ValueString(),
		Name:            plan.Name.ValueString(),
		ExternalAccount: externalAccount,
		Secret:          secret,
		UseID:           plan.UseID.ValueString(),
		IdentifierTags:  tags,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating connector installation", err.Error())
		return
	}
	mapConnectorInstallationToState(ctx, &resp.Diagnostics, &plan, installation, plan.Secret, knownExternalAccount(plan.ExternalAccountJSON))
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConnectorInstallationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state ConnectorInstallationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := parseTerraformUUID(state.ID, "connector installation ID")
	if err != nil {
		resp.Diagnostics.AddError("Error reading connector installation", err.Error())
		return
	}
	installation, err := r.client.GetConnectorInstallation(id)
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading connector installation", err.Error())
		return
	}
	mapConnectorInstallationToState(ctx, &resp.Diagnostics, &state, installation, state.Secret, state.ExternalAccountJSON)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ConnectorInstallationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan ConnectorInstallationResourceModel
	var state ConnectorInstallationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Secret.IsNull() || plan.Secret.IsUnknown() || plan.Secret.ValueString() == state.Secret.ValueString() {
		if !plan.Secret.IsUnknown() {
			state.Secret = plan.Secret
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}
	id, err := parseTerraformUUID(plan.ID, "connector installation ID")
	if err != nil {
		resp.Diagnostics.AddError("Error rotating connector credential", err.Error())
		return
	}
	installation, err := r.client.RotateConnectorCredential(id, apiclient.RotateConnectorCredentialRequest{Secret: plan.Secret.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error rotating connector credential", err.Error())
		return
	}
	mapConnectorInstallationToState(ctx, &resp.Diagnostics, &plan, installation, plan.Secret, knownExternalAccount(plan.ExternalAccountJSON))
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConnectorInstallationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state ConnectorInstallationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := parseTerraformUUID(state.ID, "connector installation ID")
	if err != nil {
		resp.Diagnostics.AddError("Error deleting connector installation", err.Error())
		return
	}
	if err := r.client.DeleteConnectorInstallation(id); err != nil && !apiclient.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting connector installation", err.Error())
	}
}

func (r *ConnectorInstallationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := parseTerraformUUID(types.StringValue(req.ID), "connector installation ID")
	if err != nil {
		resp.Diagnostics.AddError("Error importing connector installation", err.Error())
		return
	}
	installation, err := r.client.GetConnectorInstallation(id)
	if err != nil {
		resp.Diagnostics.AddError("Error importing connector installation", err.Error())
		return
	}
	state := ConnectorInstallationResourceModel{Secret: types.StringNull()}
	mapConnectorInstallationToState(
		ctx,
		&resp.Diagnostics,
		&state,
		installation,
		types.StringNull(),
		importedExternalAccount(installation),
	)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// importedExternalAccount reconstructs the non-secret selector that the API
// exposes as externalProjectId. The selector is kept in state so importing an
// installation does not immediately plan a replacement for the same account.
// Google Ads names this selector customerId; other currently supported
// project-based connectors use projectId. Provider-specific fields that are
// not persisted by the API remain unknown and can be supplied explicitly.
func importedExternalAccount(installation *apiclient.ConnectorInstallation) types.String {
	projectID := installation.ExternalProjectID
	if projectID == "" {
		return types.StringNull()
	}
	field := "projectId"
	if installation.ProviderID == "google_ads" {
		field = "customerId"
	}
	raw, err := json.Marshal(map[string]string{field: projectID})
	if err != nil {
		return types.StringNull()
	}
	return types.StringValue(string(raw))
}

func optionalJSON(value types.String, name string) (json.RawMessage, error) {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return nil, nil
	}
	normalized, err := normalizeJSONString(value.ValueString())
	if err != nil {
		return nil, fmt.Errorf("%s must be valid JSON: %w", name, err)
	}
	return json.RawMessage(normalized), nil
}

func connectorTags(ctx context.Context, diags *diag.Diagnostics, value types.Map) (map[string]string, bool) {
	tags := map[string]string{}
	if value.IsNull() || value.IsUnknown() {
		return tags, true
	}
	diags.Append(value.ElementsAs(ctx, &tags, false)...)
	return tags, !diags.HasError()
}

func mapConnectorInstallationToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *ConnectorInstallationResourceModel,
	installation *apiclient.ConnectorInstallation,
	existingSecret types.String,
	existingExternalAccount types.String,
) {
	model.ID = types.StringValue(installation.ID)
	model.ProviderID = types.StringValue(installation.ProviderID)
	model.AuthProfileID = types.StringValue(installation.AuthProfileID)
	model.UseID = types.StringValue(installation.UseID)
	model.Name = types.StringValue(installation.Name)
	model.Secret = existingSecret
	model.ExternalAccountJSON = existingExternalAccount
	model.CredentialType = types.StringValue(installation.CredentialType)
	model.Status = types.StringValue(installation.Status)
	model.StatusReason = optionalString(installation.StatusReason)
	model.Enabled = types.BoolValue(installation.Enabled)
	model.APISecretLastFour = optionalString(installation.APISecretLastFour)
	model.SecretVersion = types.Int64Value(int64(installation.SecretVersion))
	model.CreatedAt = optionalString(installation.CreatedAt)
	model.UpdatedAt = optionalString(installation.UpdatedAt)
	tags, tagDiags := types.MapValueFrom(ctx, types.StringType, installation.IdentifierTags)
	diags.Append(tagDiags...)
	model.IdentifierTags = tags
}

func knownExternalAccount(value types.String) types.String {
	if value.IsUnknown() {
		return types.StringNull()
	}
	return value
}
