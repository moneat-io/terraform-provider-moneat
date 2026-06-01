package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &SecuritySignalsDataSource{}
var _ datasource.DataSource = &SecurityDetectionCoverageDataSource{}
var _ datasource.DataSource = &SecurityVulnerabilitySummaryDataSource{}
var _ datasource.DataSource = &SecurityVulnerabilityFindingsDataSource{}
var _ datasource.DataSource = &SecurityVulnerabilityInventoryDataSource{}

type SecurityJSONDataSourceModel struct {
	Status   types.String `tfsdk:"status"`
	Severity types.String `tfsdk:"severity"`
	Source   types.String `tfsdk:"source"`
	Package  types.String `tfsdk:"package"`
	Search   types.String `tfsdk:"search"`
	Target   types.String `tfsdk:"target"`
	Limit    types.Int64  `tfsdk:"limit"`
	Offset   types.Int64  `tfsdk:"offset"`
	JSON     types.String `tfsdk:"json"`
}

type RawJSONDataSourceModel struct {
	JSON types.String `tfsdk:"json"`
}

type SecuritySignalsDataSource struct {
	client *apiclient.Client
}

func NewSecuritySignalsDataSource() datasource.DataSource {
	return &SecuritySignalsDataSource{}
}

func (d *SecuritySignalsDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_security_signals"
}

func (d *SecuritySignalsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Reads Moneat security signals as JSON.",
		Attributes:  signalFilterAttributes(),
	}
}

func (d *SecuritySignalsDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	configureSecurityDataSource(req.ProviderData, &d.client, resp)
}

func (d *SecuritySignalsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config SecurityJSONDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListSecuritySignals(securityQueryParams(config))
	if err != nil {
		resp.Diagnostics.AddError("Error reading security signals", err.Error())
		return
	}
	config.JSON = types.StringValue(rawMessageString(result))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

type SecurityDetectionCoverageDataSource struct {
	client *apiclient.Client
}

func NewSecurityDetectionCoverageDataSource() datasource.DataSource {
	return &SecurityDetectionCoverageDataSource{}
}

func (d *SecurityDetectionCoverageDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_security_detection_coverage"
}

func (d *SecurityDetectionCoverageDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Reads Moneat security detection coverage as JSON.",
		Attributes:  jsonOnlyAttributes(),
	}
}

func (d *SecurityDetectionCoverageDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	configureSecurityDataSource(req.ProviderData, &d.client, resp)
}

func (d *SecurityDetectionCoverageDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	result, err := d.client.GetSecurityDetectionCoverage()
	if err != nil {
		resp.Diagnostics.AddError("Error reading security detection coverage", err.Error())
		return
	}
	state := RawJSONDataSourceModel{JSON: types.StringValue(rawMessageString(result))}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type SecurityVulnerabilitySummaryDataSource struct {
	client *apiclient.Client
}

func NewSecurityVulnerabilitySummaryDataSource() datasource.DataSource {
	return &SecurityVulnerabilitySummaryDataSource{}
}

func (d *SecurityVulnerabilitySummaryDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_security_vulnerability_summary"
}

func (d *SecurityVulnerabilitySummaryDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Reads Moneat vulnerability summary as JSON.",
		Attributes:  jsonOnlyAttributes(),
	}
}

func (d *SecurityVulnerabilitySummaryDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	configureSecurityDataSource(req.ProviderData, &d.client, resp)
}

func (d *SecurityVulnerabilitySummaryDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	result, err := d.client.GetVulnerabilitySummary()
	if err != nil {
		resp.Diagnostics.AddError("Error reading vulnerability summary", err.Error())
		return
	}
	state := RawJSONDataSourceModel{JSON: types.StringValue(rawMessageString(result))}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type SecurityVulnerabilityFindingsDataSource struct {
	client *apiclient.Client
}

func NewSecurityVulnerabilityFindingsDataSource() datasource.DataSource {
	return &SecurityVulnerabilityFindingsDataSource{}
}

func (d *SecurityVulnerabilityFindingsDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_security_vulnerability_findings"
}

func (d *SecurityVulnerabilityFindingsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Reads Moneat vulnerability findings as JSON.",
		Attributes:  vulnerabilityFilterAttributes(),
	}
}

func (d *SecurityVulnerabilityFindingsDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	configureSecurityDataSource(req.ProviderData, &d.client, resp)
}

func (d *SecurityVulnerabilityFindingsDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config SecurityJSONDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListVulnerabilityFindings(securityQueryParams(config))
	if err != nil {
		resp.Diagnostics.AddError("Error reading vulnerability findings", err.Error())
		return
	}
	config.JSON = types.StringValue(rawMessageString(result))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

type SecurityVulnerabilityInventoryDataSource struct {
	client *apiclient.Client
}

func NewSecurityVulnerabilityInventoryDataSource() datasource.DataSource {
	return &SecurityVulnerabilityInventoryDataSource{}
}

func (d *SecurityVulnerabilityInventoryDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_security_vulnerability_inventory"
}

func (d *SecurityVulnerabilityInventoryDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Reads Moneat vulnerability package inventory as JSON.",
		Attributes:  vulnerabilityInventoryFilterAttributes(),
	}
}

func (d *SecurityVulnerabilityInventoryDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	configureSecurityDataSource(req.ProviderData, &d.client, resp)
}

func (d *SecurityVulnerabilityInventoryDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config SecurityJSONDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListVulnerabilityInventory(securityQueryParams(config))
	if err != nil {
		resp.Diagnostics.AddError("Error reading vulnerability inventory", err.Error())
		return
	}
	config.JSON = types.StringValue(rawMessageString(result))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func configureSecurityDataSource(
	providerData interface{},
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

func jsonOnlyAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"json": schema.StringAttribute{
			Description: "Raw API response JSON.",
			Computed:    true,
		},
	}
}

func signalFilterAttributes() map[string]schema.Attribute {
	attrs := jsonOnlyAttributes()
	attrs["status"] = optionalFilter("Signal status filter.")
	attrs["severity"] = optionalFilter("Signal severity filter.")
	attrs["source"] = optionalFilter("Signal source filter.")
	attrs["limit"] = optionalIntFilter("Maximum rows to return.")
	attrs["offset"] = optionalIntFilter("Rows to skip.")
	return attrs
}

func vulnerabilityFilterAttributes() map[string]schema.Attribute {
	attrs := jsonOnlyAttributes()
	attrs["status"] = optionalFilter("Finding status filter.")
	attrs["severity"] = optionalFilter("Finding severity filter.")
	attrs["package"] = optionalFilter("Package name filter.")
	attrs["search"] = optionalFilter("Search filter.")
	attrs["target"] = optionalFilter("Target filter.")
	attrs["limit"] = optionalIntFilter("Maximum rows to return.")
	attrs["offset"] = optionalIntFilter("Rows to skip.")
	return attrs
}

func vulnerabilityInventoryFilterAttributes() map[string]schema.Attribute {
	attrs := jsonOnlyAttributes()
	attrs["package"] = optionalFilter("Package name filter.")
	attrs["search"] = optionalFilter("Search filter.")
	attrs["target"] = optionalFilter("Target filter.")
	attrs["limit"] = optionalIntFilter("Maximum rows to return.")
	attrs["offset"] = optionalIntFilter("Rows to skip.")
	return attrs
}

func optionalFilter(description string) schema.StringAttribute {
	return schema.StringAttribute{
		Description: description,
		Optional:    true,
	}
}

func optionalIntFilter(description string) schema.Int64Attribute {
	return schema.Int64Attribute{
		Description: description,
		Optional:    true,
	}
}

func securityQueryParams(config SecurityJSONDataSourceModel) map[string]string {
	params := map[string]string{}
	addStringParam(params, "status", config.Status)
	addStringParam(params, "severity", config.Severity)
	addStringParam(params, "source", config.Source)
	addStringParam(params, "package", config.Package)
	addStringParam(params, "search", config.Search)
	addStringParam(params, "target", config.Target)
	addIntParam(params, "limit", config.Limit)
	addIntParam(params, "offset", config.Offset)
	return params
}

func addStringParam(params map[string]string, key string, value types.String) {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return
	}
	params[key] = value.ValueString()
}

func addIntParam(params map[string]string, key string, value types.Int64) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	params[key] = fmt.Sprintf("%d", value.ValueInt64())
}
