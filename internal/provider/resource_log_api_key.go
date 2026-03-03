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

var _ resource.Resource = &LogAPIKeyResource{}
var _ resource.ResourceWithImportState = &LogAPIKeyResource{}

type LogAPIKeyResource struct {
	client *apiclient.Client
}

type LogAPIKeyResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Key  types.String `tfsdk:"key"`
}

func NewLogAPIKeyResource() resource.Resource {
	return &LogAPIKeyResource{}
}

func (r *LogAPIKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_api_key"
}

func (r *LogAPIKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat log API key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the log API key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the log API key.",
				Required:    true,
			},
			"key": schema.StringAttribute{
				Description: "The generated API key value. Only available after creation.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *LogAPIKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LogAPIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LogAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateLogAPIKeyRequest{
		Name: plan.Name.ValueString(),
	}

	key, err := r.client.CreateLogAPIKey(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating log API key", err.Error())
		return
	}

	plan.ID = types.StringValue(key.ID)
	plan.Name = types.StringValue(key.Name)
	plan.Key = types.StringValue(key.Key)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LogAPIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LogAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GetLogAPIKey(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading log API key", err.Error())
		return
	}

	state.Name = types.StringValue(key.Name)
	if key.Key != "" {
		state.Key = types.StringValue(key.Key)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LogAPIKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LogAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state LogAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateLogAPIKeyRequest{
		Name: plan.Name.ValueString(),
	}

	key, err := r.client.UpdateLogAPIKey(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating log API key", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Key = state.Key // Preserve key from state
	plan.Name = types.StringValue(key.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LogAPIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LogAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLogAPIKey(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting log API key", err.Error())
		return
	}
}

func (r *LogAPIKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	key, err := r.client.GetLogAPIKey(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing log API key", err.Error())
		return
	}

	state := LogAPIKeyResourceModel{
		ID:   types.StringValue(key.ID),
		Name: types.StringValue(key.Name),
		Key:  types.StringValue(key.Key),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
