package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &McpToolCatalogDataSource{}

type McpToolCatalogDataSource struct {
	client *apiclient.Client
}

type McpToolCatalogDataSourceModel struct {
	JSON types.String `tfsdk:"json"`
}

func NewMcpToolCatalogDataSource() datasource.DataSource {
	return &McpToolCatalogDataSource{}
}

func (d *McpToolCatalogDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_mcp_tool_catalog"
}

func (d *McpToolCatalogDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Reads the Moneat MCP tool and resource catalog as JSON.",
		Attributes: map[string]schema.Attribute{
			"json": schema.StringAttribute{
				Description: "Raw MCP tool catalog JSON.",
				Computed:    true,
			},
		},
	}
}

func (d *McpToolCatalogDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *apiclient.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *McpToolCatalogDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	catalog, err := d.client.GetMcpToolCatalog()
	if err != nil {
		resp.Diagnostics.AddError("Error reading MCP tool catalog", err.Error())
		return
	}
	state := McpToolCatalogDataSourceModel{
		JSON: types.StringValue(rawMessageString(catalog)),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
