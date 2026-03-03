package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &StatusPagesDataSource{}

type StatusPagesDataSource struct {
	client *apiclient.Client
}

type StatusPageDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
	IsPublic    types.Bool   `tfsdk:"is_public"`
}

type StatusPagesDataSourceModel struct {
	StatusPages []StatusPageDataSourceModel `tfsdk:"status_pages"`
}

func NewStatusPagesDataSource() datasource.DataSource {
	return &StatusPagesDataSource{}
}

func (d *StatusPagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_pages"
}

func (d *StatusPagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list all Moneat status pages.",
		Attributes: map[string]schema.Attribute{
			"status_pages": schema.ListNestedAttribute{
				Description: "List of status pages.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the status page.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *StatusPagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPagesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	pages, err := d.client.ListStatusPages()
	if err != nil {
		resp.Diagnostics.AddError("Error listing status pages", err.Error())
		return
	}

	var state StatusPagesDataSourceModel
	for _, p := range pages {
		state.StatusPages = append(state.StatusPages, StatusPageDataSourceModel{
			ID:          types.StringValue(p.ID),
			Name:        types.StringValue(p.Name),
			Slug:        types.StringValue(p.Slug),
			Description: types.StringValue(p.Description),
			IsPublic:    types.BoolValue(p.IsPublic),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
