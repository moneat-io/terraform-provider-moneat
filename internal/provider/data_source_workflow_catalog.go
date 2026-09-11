package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &WorkflowCatalogDataSource{}
var _ datasource.DataSource = &WorkflowBlueprintDataSource{}
var _ datasource.DataSource = &ConnectorProvidersDataSource{}
var _ datasource.DataSource = &ConnectorInstallationsDataSource{}
var _ datasource.DataSource = &ConnectorGroupsDataSource{}
var _ datasource.DataSource = &ConnectorResourcesDataSource{}

type rawCatalogDataSource struct {
	client *apiclient.Client
}

type rawCatalogModel struct {
	JSON types.String `tfsdk:"json"`
}

type WorkflowCatalogDataSource struct{ rawCatalogDataSource }

func NewWorkflowCatalogDataSource() datasource.DataSource { return &WorkflowCatalogDataSource{} }

func (d *WorkflowCatalogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_catalog"
}

func (d *WorkflowCatalogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Reads the Moneat workflow action catalog as JSON.", Attributes: rawJSONAttributes()}
}

func (d *WorkflowCatalogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureRawCatalogClient(req.ProviderData, &d.client, resp)
}

func (d *WorkflowCatalogDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !workflowClientConfigured(d.client, &resp.Diagnostics, "WorkflowCatalogDataSource.Read") {
		return
	}
	value, err := d.client.GetWorkflowCatalog()
	if err != nil {
		resp.Diagnostics.AddError("Error reading workflow catalog", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, rawCatalogModel{JSON: types.StringValue(rawMessageString(value))})...)
}

type WorkflowBlueprintDataSource struct{ rawCatalogDataSource }

type WorkflowBlueprintDataSourceModel struct {
	Key  types.String `tfsdk:"key"`
	JSON types.String `tfsdk:"json"`
}

func NewWorkflowBlueprintDataSource() datasource.DataSource { return &WorkflowBlueprintDataSource{} }

func (d *WorkflowBlueprintDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_blueprint"
}

func (d *WorkflowBlueprintDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Reads one portable Moneat workflow blueprint as JSON.", Attributes: map[string]schema.Attribute{
		"key":  schema.StringAttribute{Required: true},
		"json": schema.StringAttribute{Computed: true},
	}}
}

func (d *WorkflowBlueprintDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureRawCatalogClient(req.ProviderData, &d.client, resp)
}

func (d *WorkflowBlueprintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !workflowClientConfigured(d.client, &resp.Diagnostics, "WorkflowBlueprintDataSource.Read") {
		return
	}
	var config WorkflowBlueprintDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	value, err := d.client.GetWorkflowBlueprint(config.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading workflow blueprint", err.Error())
		return
	}
	config.JSON = types.StringValue(rawMessageString(value))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

type ConnectorProvidersDataSource struct{ rawCatalogDataSource }

func NewConnectorProvidersDataSource() datasource.DataSource { return &ConnectorProvidersDataSource{} }

func (d *ConnectorProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector_providers"
}

func (d *ConnectorProvidersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Reads connector providers and uses without credentials.", Attributes: rawJSONAttributes()}
}

func (d *ConnectorProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureRawCatalogClient(req.ProviderData, &d.client, resp)
}

func (d *ConnectorProvidersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !workflowClientConfigured(d.client, &resp.Diagnostics, "ConnectorProvidersDataSource.Read") {
		return
	}
	value, err := d.client.ListConnectorProviders()
	if err != nil {
		resp.Diagnostics.AddError("Error reading connector providers", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, rawCatalogModel{JSON: types.StringValue(rawMessageString(value))})...)
}

type ConnectorInstallationsDataSource struct{ rawCatalogDataSource }

type ConnectorInstallationsDataSourceModel struct {
	ProviderID    types.String                           `tfsdk:"provider_id"`
	UseID         types.String                           `tfsdk:"use_id"`
	JSON          types.String                           `tfsdk:"json"`
	Installations []ConnectorInstallationDataSourceModel `tfsdk:"installations"`
}

type ConnectorInstallationDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	ProviderID        types.String `tfsdk:"provider_id"`
	UseID             types.String `tfsdk:"use_id"`
	Name              types.String `tfsdk:"name"`
	AuthProfileID     types.String `tfsdk:"auth_profile_id"`
	Status            types.String `tfsdk:"status"`
	CredentialType    types.String `tfsdk:"credential_type"`
	IdentifierTags    types.Map    `tfsdk:"identifier_tags"`
	APISecretLastFour types.String `tfsdk:"api_secret_last_four"`
	SecretVersion     types.Int64  `tfsdk:"secret_version"`
}

func NewConnectorInstallationsDataSource() datasource.DataSource {
	return &ConnectorInstallationsDataSource{}
}

func (d *ConnectorInstallationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector_installations"
}

func (d *ConnectorInstallationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Reads redacted connector installations; no credential material is returned.", Attributes: map[string]schema.Attribute{
		"provider_id": schema.StringAttribute{Optional: true},
		"use_id":      schema.StringAttribute{Optional: true},
		"json":        schema.StringAttribute{Computed: true},
		"installations": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true}, "provider_id": schema.StringAttribute{Computed: true}, "use_id": schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Computed: true}, "auth_profile_id": schema.StringAttribute{Computed: true}, "status": schema.StringAttribute{Computed: true},
			"credential_type": schema.StringAttribute{Computed: true}, "identifier_tags": schema.MapAttribute{Computed: true, ElementType: types.StringType},
			"api_secret_last_four": schema.StringAttribute{Computed: true}, "secret_version": schema.Int64Attribute{Computed: true},
		}}},
	}}
}

func (d *ConnectorInstallationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureRawCatalogClient(req.ProviderData, &d.client, resp)
}

func (d *ConnectorInstallationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !workflowClientConfigured(d.client, &resp.Diagnostics, "ConnectorInstallationsDataSource.Read") {
		return
	}
	var config ConnectorInstallationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	installations, err := d.client.ListConnectorInstallations(config.ProviderID.ValueString(), config.UseID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading connector installations", err.Error())
		return
	}
	config.JSON = types.StringValue(rawConnectorInstallationsJSON(installations))
	config.Installations = make([]ConnectorInstallationDataSourceModel, 0, len(installations))
	for _, installation := range installations {
		tags, tagDiags := types.MapValueFrom(ctx, types.StringType, installation.IdentifierTags)
		resp.Diagnostics.Append(tagDiags...)
		config.Installations = append(config.Installations, ConnectorInstallationDataSourceModel{
			ID: types.StringValue(installation.ID), ProviderID: types.StringValue(installation.ProviderID), UseID: types.StringValue(installation.UseID),
			Name: types.StringValue(installation.Name), AuthProfileID: types.StringValue(installation.AuthProfileID), Status: types.StringValue(installation.Status),
			CredentialType: types.StringValue(installation.CredentialType), IdentifierTags: tags,
			APISecretLastFour: optionalString(installation.APISecretLastFour), SecretVersion: types.Int64Value(int64(installation.SecretVersion)),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

type ConnectorGroupsDataSource struct{ rawCatalogDataSource }

type ConnectorGroupsDataSourceModel struct {
	ProviderID types.String                    `tfsdk:"provider_id"`
	UseID      types.String                    `tfsdk:"use_id"`
	JSON       types.String                    `tfsdk:"json"`
	Groups     []ConnectorGroupDataSourceModel `tfsdk:"groups"`
}

type ConnectorGroupDataSourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	ProviderID            types.String `tfsdk:"provider_id"`
	UseID                 types.String `tfsdk:"use_id"`
	MemberInstallationIDs types.List   `tfsdk:"member_installation_ids"`
	SelectionStrategy     types.String `tfsdk:"selection_strategy"`
}

func NewConnectorGroupsDataSource() datasource.DataSource { return &ConnectorGroupsDataSource{} }

func (d *ConnectorGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector_groups"
}

func (d *ConnectorGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Reads connector connection groups and UUID installation membership.", Attributes: map[string]schema.Attribute{
		"provider_id": schema.StringAttribute{Optional: true}, "use_id": schema.StringAttribute{Optional: true}, "json": schema.StringAttribute{Computed: true},
		"groups": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true}, "name": schema.StringAttribute{Computed: true}, "provider_id": schema.StringAttribute{Computed: true}, "use_id": schema.StringAttribute{Computed: true},
			"member_installation_ids": schema.ListAttribute{Computed: true, ElementType: types.StringType}, "selection_strategy": schema.StringAttribute{Computed: true},
		}}},
	}}
}

func (d *ConnectorGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureRawCatalogClient(req.ProviderData, &d.client, resp)
}

func (d *ConnectorGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !workflowClientConfigured(d.client, &resp.Diagnostics, "ConnectorGroupsDataSource.Read") {
		return
	}
	var config ConnectorGroupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	groups, err := d.client.ListConnectorGroups(config.ProviderID.ValueString(), config.UseID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading connector groups", err.Error())
		return
	}
	config.JSON = types.StringValue(rawConnectorGroupsJSON(groups))
	config.Groups = make([]ConnectorGroupDataSourceModel, 0, len(groups))
	for _, group := range groups {
		members, memberDiags := types.ListValueFrom(ctx, types.StringType, group.MemberInstallationIDs)
		resp.Diagnostics.Append(memberDiags...)
		config.Groups = append(config.Groups, ConnectorGroupDataSourceModel{ID: types.StringValue(group.ID), Name: types.StringValue(group.Name), ProviderID: types.StringValue(group.ProviderID), UseID: types.StringValue(group.UseID), MemberInstallationIDs: members, SelectionStrategy: types.StringValue(group.SelectionStrategy)})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

type ConnectorResourcesDataSource struct{ rawCatalogDataSource }

type ConnectorResourcesDataSourceModel struct {
	InstallationID types.String                       `tfsdk:"installation_id"`
	JSON           types.String                       `tfsdk:"json"`
	Resources      []ConnectorResourceDataSourceModel `tfsdk:"resources"`
}

type ConnectorResourceDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	InstallationID       types.String `tfsdk:"installation_id"`
	ExternalResourceType types.String `tfsdk:"external_resource_type"`
	ExternalResourceID   types.String `tfsdk:"external_resource_id"`
	DisplayName          types.String `tfsdk:"display_name"`
	ProviderMetadata     types.Map    `tfsdk:"provider_metadata"`
}

func NewConnectorResourcesDataSource() datasource.DataSource { return &ConnectorResourcesDataSource{} }

func (d *ConnectorResourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector_resources"
}

func (d *ConnectorResourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Reads selectable external resources for a connector installation.", Attributes: map[string]schema.Attribute{
		"installation_id": schema.StringAttribute{Required: true}, "json": schema.StringAttribute{Computed: true},
		"resources": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true}, "installation_id": schema.StringAttribute{Computed: true}, "external_resource_type": schema.StringAttribute{Computed: true},
			"external_resource_id": schema.StringAttribute{Computed: true}, "display_name": schema.StringAttribute{Computed: true}, "provider_metadata": schema.MapAttribute{Computed: true, ElementType: types.StringType},
		}}},
	}}
}

func (d *ConnectorResourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureRawCatalogClient(req.ProviderData, &d.client, resp)
}

func (d *ConnectorResourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !workflowClientConfigured(d.client, &resp.Diagnostics, "ConnectorResourcesDataSource.Read") {
		return
	}
	var config ConnectorResourcesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := parseTerraformUUID(config.InstallationID, "connector installation ID")
	if err != nil {
		resp.Diagnostics.AddError("Invalid connector installation ID", err.Error())
		return
	}
	resources, err := d.client.ListConnectorResources(id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading connector resources", err.Error())
		return
	}
	config.JSON = types.StringValue(rawConnectorResourcesJSON(resources))
	config.Resources = make([]ConnectorResourceDataSourceModel, 0, len(resources))
	for _, resource := range resources {
		metadata, metadataDiags := types.MapValueFrom(ctx, types.StringType, resource.ProviderMetadata)
		resp.Diagnostics.Append(metadataDiags...)
		config.Resources = append(config.Resources, ConnectorResourceDataSourceModel{ID: types.StringValue(resource.ID), InstallationID: types.StringValue(resource.InstallationID), ExternalResourceType: types.StringValue(resource.ExternalResourceType), ExternalResourceID: types.StringValue(resource.ExternalResourceID), DisplayName: optionalString(resource.DisplayName), ProviderMetadata: metadata})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func rawJSONAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{"json": schema.StringAttribute{Description: "Raw JSON returned by the Moneat API.", Computed: true}}
}

func configureRawCatalogClient(providerData any, target **apiclient.Client, resp *datasource.ConfigureResponse) {
	if providerData == nil {
		return
	}
	client, ok := providerData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *apiclient.Client, got: %T", providerData))
		return
	}
	*target = client
}

func rawConnectorInstallationsJSON(value []apiclient.ConnectorInstallation) string {
	return rawMessageString(mustJSON(value))
}

func rawConnectorGroupsJSON(value []apiclient.ConnectorGroup) string {
	return rawMessageString(mustJSON(value))
}
func rawConnectorResourcesJSON(value []apiclient.ConnectorExternalResource) string {
	return rawMessageString(mustJSON(value))
}

func mustJSON(value any) json.RawMessage {
	encoded, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage("null")
	}
	return encoded
}
