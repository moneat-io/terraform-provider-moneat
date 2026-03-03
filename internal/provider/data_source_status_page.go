package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &StatusPageDataSource{}

type StatusPageDataSource struct {
	client *apiclient.Client
}

func NewStatusPageDataSource() datasource.DataSource {
	return &StatusPageDataSource{}
}

func (d *StatusPageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (d *StatusPageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Moneat status page by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the status page.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the status page.",
				Computed:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The URL slug of the status page.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the status page.",
				Computed:    true,
			},
			"is_public": schema.BoolAttribute{
				Description: "Whether the status page is publicly accessible.",
				Computed:    true,
			},
		},
	}
}

func (d *StatusPageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config StatusPageDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	page, err := d.client.GetStatusPage(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading status page", err.Error())
		return
	}

	config.Name = types.StringValue(page.Name)
	config.Slug = types.StringValue(page.Slug)
	config.Description = types.StringValue(page.Description)
	config.IsPublic = types.BoolValue(page.IsPublic)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
