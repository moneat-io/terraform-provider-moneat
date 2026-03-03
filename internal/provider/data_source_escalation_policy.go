package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &EscalationPolicyDataSource{}

type EscalationPolicyDataSource struct {
	client *apiclient.Client
}

type EscalationPolicyDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func NewEscalationPolicyDataSource() datasource.DataSource {
	return &EscalationPolicyDataSource{}
}

func (d *EscalationPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_escalation_policy"
}

func (d *EscalationPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Moneat escalation policy by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the escalation policy.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the escalation policy.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the escalation policy.",
				Computed:    true,
			},
		},
	}
}

func (d *EscalationPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EscalationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config EscalationPolicyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := d.client.GetEscalationPolicy(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading escalation policy", err.Error())
		return
	}

	config.Name = types.StringValue(policy.Name)
	config.Description = types.StringValue(policy.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
