package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &SyntheticTestResource{}
var _ resource.ResourceWithImportState = &SyntheticTestResource{}

type SyntheticTestResource struct {
	client *apiclient.Client
}

type SyntheticTestResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	URL             types.String `tfsdk:"url"`
	IntervalSeconds types.Int64  `tfsdk:"interval_seconds"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

func NewSyntheticTestResource() resource.Resource {
	return &SyntheticTestResource{}
}

func (r *SyntheticTestResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synthetic_test"
}

func (r *SyntheticTestResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat synthetic test.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the synthetic test.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the synthetic test.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of synthetic test (http, browser, api).",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL to test.",
				Required:    true,
			},
			"interval_seconds": schema.Int64Attribute{
				Description: "The test interval in seconds.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the synthetic test is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *SyntheticTestResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SyntheticTestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SyntheticTestResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateSyntheticTestRequest{
		Name:            plan.Name.ValueString(),
		Type:            plan.Type.ValueString(),
		URL:             plan.URL.ValueString(),
		IntervalSeconds: plan.IntervalSeconds.ValueInt64(),
		Enabled:         plan.Enabled.ValueBool(),
	}

	test, err := r.client.CreateSyntheticTest(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating synthetic test", err.Error())
		return
	}

	plan.ID = types.StringValue(test.ID)
	plan.Name = types.StringValue(test.Name)
	plan.Type = types.StringValue(test.Type)
	plan.URL = types.StringValue(test.URL)
	plan.IntervalSeconds = types.Int64Value(test.IntervalSeconds)
	plan.Enabled = types.BoolValue(test.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SyntheticTestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SyntheticTestResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	test, err := r.client.GetSyntheticTest(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading synthetic test", err.Error())
		return
	}

	state.Name = types.StringValue(test.Name)
	state.Type = types.StringValue(test.Type)
	state.URL = types.StringValue(test.URL)
	state.IntervalSeconds = types.Int64Value(test.IntervalSeconds)
	state.Enabled = types.BoolValue(test.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SyntheticTestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SyntheticTestResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SyntheticTestResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateSyntheticTestRequest{
		Name:            plan.Name.ValueString(),
		Type:            plan.Type.ValueString(),
		URL:             plan.URL.ValueString(),
		IntervalSeconds: plan.IntervalSeconds.ValueInt64(),
		Enabled:         plan.Enabled.ValueBool(),
	}

	test, err := r.client.UpdateSyntheticTest(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating synthetic test", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(test.Name)
	plan.Type = types.StringValue(test.Type)
	plan.URL = types.StringValue(test.URL)
	plan.IntervalSeconds = types.Int64Value(test.IntervalSeconds)
	plan.Enabled = types.BoolValue(test.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SyntheticTestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SyntheticTestResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSyntheticTest(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting synthetic test", err.Error())
		return
	}
}

func (r *SyntheticTestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	test, err := r.client.GetSyntheticTest(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing synthetic test", err.Error())
		return
	}

	state := SyntheticTestResourceModel{
		ID:              types.StringValue(test.ID),
		Name:            types.StringValue(test.Name),
		Type:            types.StringValue(test.Type),
		URL:             types.StringValue(test.URL),
		IntervalSeconds: types.Int64Value(test.IntervalSeconds),
		Enabled:         types.BoolValue(test.Enabled),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
