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

var _ resource.Resource = &LogIndexResource{}
var _ resource.ResourceWithImportState = &LogIndexResource{}

type LogIndexResource struct {
	client *apiclient.Client
}

type LogIndexResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	FilterQuery   types.String `tfsdk:"filter_query"`
	RetentionDays types.Int64  `tfsdk:"retention_days"`
}

func NewLogIndexResource() resource.Resource {
	return &LogIndexResource{}
}

func (r *LogIndexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_index"
}

func (r *LogIndexResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat log index.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the log index.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the log index.",
				Required:    true,
			},
			"filter_query": schema.StringAttribute{
				Description: "The filter query for the log index.",
				Required:    true,
			},
			"retention_days": schema.Int64Attribute{
				Description: "The retention period in days.",
				Required:    true,
			},
		},
	}
}

func (r *LogIndexResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LogIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LogIndexResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateLogIndexRequest{
		Name:          plan.Name.ValueString(),
		FilterQuery:   plan.FilterQuery.ValueString(),
		RetentionDays: plan.RetentionDays.ValueInt64(),
	}

	index, err := r.client.CreateLogIndex(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating log index", err.Error())
		return
	}

	plan.ID = types.StringValue(index.ID)
	plan.Name = types.StringValue(index.Name)
	plan.FilterQuery = types.StringValue(index.FilterQuery)
	plan.RetentionDays = types.Int64Value(index.RetentionDays)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LogIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LogIndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	index, err := r.client.GetLogIndex(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading log index", err.Error())
		return
	}

	state.Name = types.StringValue(index.Name)
	state.FilterQuery = types.StringValue(index.FilterQuery)
	state.RetentionDays = types.Int64Value(index.RetentionDays)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LogIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LogIndexResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state LogIndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateLogIndexRequest{
		Name:          plan.Name.ValueString(),
		FilterQuery:   plan.FilterQuery.ValueString(),
		RetentionDays: plan.RetentionDays.ValueInt64(),
	}

	index, err := r.client.UpdateLogIndex(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating log index", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(index.Name)
	plan.FilterQuery = types.StringValue(index.FilterQuery)
	plan.RetentionDays = types.Int64Value(index.RetentionDays)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LogIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LogIndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLogIndex(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting log index", err.Error())
		return
	}
}

func (r *LogIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	index, err := r.client.GetLogIndex(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing log index", err.Error())
		return
	}

	state := LogIndexResourceModel{
		ID:            types.StringValue(index.ID),
		Name:          types.StringValue(index.Name),
		FilterQuery:   types.StringValue(index.FilterQuery),
		RetentionDays: types.Int64Value(index.RetentionDays),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
