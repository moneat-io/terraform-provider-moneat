package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &HostsDataSource{}

type HostsDataSource struct {
	client *apiclient.Client
}

type HostsDataSourceModel struct {
	Hosts []HostDataSourceModel `tfsdk:"hosts"`
}

func NewHostsDataSource() datasource.DataSource {
	return &HostsDataSource{}
}

func (d *HostsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hosts"
}

func (d *HostsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list all Moneat monitored hosts.",
		Attributes: map[string]schema.Attribute{
			"hosts": schema.ListNestedAttribute{
				Description: "List of monitored hosts.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the host.",
							Computed:    true,
						},
						"hostname": schema.StringAttribute{
							Description: "The hostname.",
							Computed:    true,
						},
						"os": schema.StringAttribute{
							Description: "The operating system.",
							Computed:    true,
						},
						"platform": schema.StringAttribute{
							Description: "The platform.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The host status.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *HostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HostsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	hosts, err := d.client.ListHosts()
	if err != nil {
		resp.Diagnostics.AddError("Error listing hosts", err.Error())
		return
	}

	var state HostsDataSourceModel
	for _, h := range hosts {
		state.Hosts = append(state.Hosts, HostDataSourceModel{
			ID:       types.StringValue(h.ID),
			Hostname: types.StringValue(h.Hostname),
			OS:       types.StringValue(h.OS),
			Platform: types.StringValue(h.Platform),
			Status:   types.StringValue(h.Status),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
