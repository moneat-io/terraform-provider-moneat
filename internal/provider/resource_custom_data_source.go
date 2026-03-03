package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &CustomDataSourceResource{}
var _ resource.ResourceWithImportState = &CustomDataSourceResource{}

type CustomDataSourceResource struct {
	client *apiclient.Client
}

type CustomDataSourceResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Type             types.String `tfsdk:"type"`
	ConnectionString types.String `tfsdk:"connection_string"`
}

func NewCustomDataSourceResource() resource.Resource {
	return &CustomDataSourceResource{}
}

func (r *CustomDataSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_data_source"
}

func (r *CustomDataSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat custom data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the custom data source.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the custom data source.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the custom data source.",
				Required:    true,
			},
			"connection_string": schema.StringAttribute{
				Description: "The connection string for the data source.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *CustomDataSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *apiclient.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *CustomDataSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CustomDataSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateCustomDataSourceRequest{
		Name:             plan.Name.ValueString(),
		Type:             plan.Type.ValueString(),
		ConnectionString: plan.ConnectionString.ValueString(),
	}

	ds, err := r.client.CreateCustomDataSource(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating custom data source", err.Error())
		return
	}

	plan.ID = types.StringValue(ds.ID)
	plan.Name = types.StringValue(ds.Name)
	plan.Type = types.StringValue(ds.Type)
	plan.ConnectionString = types.StringValue(ds.ConnectionString)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CustomDataSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CustomDataSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ds, err := r.client.GetCustomDataSource(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading custom data source", err.Error())
		return
	}

	state.Name = types.StringValue(ds.Name)
	state.Type = types.StringValue(ds.Type)
	state.ConnectionString = types.StringValue(ds.ConnectionString)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CustomDataSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CustomDataSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CustomDataSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateCustomDataSourceRequest{
		Name:             plan.Name.ValueString(),
		Type:             plan.Type.ValueString(),
		ConnectionString: plan.ConnectionString.ValueString(),
	}

	ds, err := r.client.UpdateCustomDataSource(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating custom data source", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(ds.Name)
	plan.Type = types.StringValue(ds.Type)
	plan.ConnectionString = types.StringValue(ds.ConnectionString)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CustomDataSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CustomDataSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCustomDataSource(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting custom data source", err.Error())
		return
	}
}

func (r *CustomDataSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ds, err := r.client.GetCustomDataSource(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing custom data source", err.Error())
		return
	}

	state := CustomDataSourceResourceModel{
		ID:               types.StringValue(ds.ID),
		Name:             types.StringValue(ds.Name),
		Type:             types.StringValue(ds.Type),
		ConnectionString: types.StringValue(ds.ConnectionString),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
