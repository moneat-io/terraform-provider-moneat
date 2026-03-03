package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &OrgMembersDataSource{}

type OrgMembersDataSource struct {
	client *apiclient.Client
}

type OrgMemberDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

type OrgMembersDataSourceModel struct {
	Members []OrgMemberDataSourceModel `tfsdk:"members"`
}

func NewOrgMembersDataSource() datasource.DataSource {
	return &OrgMembersDataSource{}
}

func (d *OrgMembersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_members"
}

func (d *OrgMembersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list all Moneat organization members.",
		Attributes: map[string]schema.Attribute{
			"members": schema.ListNestedAttribute{
				Description: "List of organization members.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the member.",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "The email address of the member.",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "The role of the member.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *OrgMembersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrgMembersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	members, err := d.client.ListOrgMembers()
	if err != nil {
		resp.Diagnostics.AddError("Error listing organization members", err.Error())
		return
	}

	var state OrgMembersDataSourceModel
	for _, m := range members {
		state.Members = append(state.Members, OrgMemberDataSourceModel{
			ID:    types.StringValue(m.ID),
			Email: types.StringValue(m.Email),
			Role:  types.StringValue(m.Role),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
