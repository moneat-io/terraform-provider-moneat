package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ datasource.DataSource = &ProjectDataSource{}

// ProjectDataSource defines the data source implementation.
type ProjectDataSource struct {
	client *apiclient.Client
}

// ProjectDataSourceModel describes the data source data model.
type ProjectDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Platform  types.String `tfsdk:"platform"`
	Framework types.String `tfsdk:"framework"`
}

// NewProjectDataSource returns a new data source factory function.
func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

func (d *ProjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Moneat project by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the project.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the project.",
				Computed:    true,
			},
			"platform": schema.StringAttribute{
				Description: "The platform of the project.",
				Computed:    true,
			},
			"framework": schema.StringAttribute{
				Description: "The framework used by the project.",
				Computed:    true,
			},
		},
	}
}

func (d *ProjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ProjectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := d.client.GetProject(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	config.Name = types.StringValue(project.Name)
	config.Platform = types.StringValue(project.Platform)
	if project.Framework != "" {
		config.Framework = types.StringValue(project.Framework)
	} else {
		config.Framework = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// --- Projects (list) data source ---

var _ datasource.DataSource = &ProjectsDataSource{}

// ProjectsDataSource defines the data source for listing all projects.
type ProjectsDataSource struct {
	client *apiclient.Client
}

// ProjectsDataSourceModel describes the data source data model.
type ProjectsDataSourceModel struct {
	Projects []ProjectDataSourceModel `tfsdk:"projects"`
}

// NewProjectsDataSource returns a new data source factory function.
func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

func (d *ProjectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list all Moneat projects.",
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Description: "List of projects.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the project.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the project.",
							Computed:    true,
						},
						"platform": schema.StringAttribute{
							Description: "The platform of the project.",
							Computed:    true,
						},
						"framework": schema.StringAttribute{
							Description: "The framework used by the project.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ProjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProjectsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	projects, err := d.client.ListProjects()
	if err != nil {
		resp.Diagnostics.AddError("Error listing projects", err.Error())
		return
	}

	var state ProjectsDataSourceModel
	for _, p := range projects {
		model := ProjectDataSourceModel{
			ID:       types.StringValue(p.ID),
			Name:     types.StringValue(p.Name),
			Platform: types.StringValue(p.Platform),
		}
		if p.Framework != "" {
			model.Framework = types.StringValue(p.Framework)
		} else {
			model.Framework = types.StringValue("")
		}
		state.Projects = append(state.Projects, model)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
