package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &OnCallScheduleSlackUsergroupResource{}

type OnCallScheduleSlackUsergroupResource struct {
	client *apiclient.Client
}

type OnCallScheduleSlackUsergroupResourceModel struct {
	ScheduleID      types.String `tfsdk:"schedule_id"`
	UsergroupID     types.String `tfsdk:"usergroup_id"`
	UsergroupHandle types.String `tfsdk:"usergroup_handle"`
}

func NewOnCallScheduleSlackUsergroupResource() resource.Resource {
	return &OnCallScheduleSlackUsergroupResource{}
}

func (r *OnCallScheduleSlackUsergroupResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_on_call_schedule_slack_usergroup"
}

func (r *OnCallScheduleSlackUsergroupResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Maps a Moneat on-call schedule to a Slack usergroup for automatic rotation updates " +
			"(Enterprise only).",
		Attributes: map[string]schema.Attribute{
			"schedule_id": schema.StringAttribute{
				Description: "The ID of the on-call schedule.",
				Required:    true,
			},
			"usergroup_id": schema.StringAttribute{
				Description: "The Slack usergroup ID (e.g., S0123456789).",
				Required:    true,
			},
			"usergroup_handle": schema.StringAttribute{
				Description: "The Slack usergroup handle (e.g., oncall-primary).",
				Required:    true,
			},
		},
	}
}

func (r *OnCallScheduleSlackUsergroupResource) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
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

func (r *OnCallScheduleSlackUsergroupResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan OnCallScheduleSlackUsergroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.SetOnCallScheduleSlackUsergroupRequest{
		UsergroupID:     plan.UsergroupID.ValueString(),
		UsergroupHandle: plan.UsergroupHandle.ValueString(),
	}

	mapping, err := r.client.SetOnCallScheduleSlackUsergroup(plan.ScheduleID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error setting Slack usergroup mapping", err.Error())
		return
	}

	plan.ScheduleID = types.StringValue(mapping.ScheduleID)
	plan.UsergroupID = types.StringValue(mapping.UsergroupID)
	plan.UsergroupHandle = types.StringValue(mapping.UsergroupHandle)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OnCallScheduleSlackUsergroupResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state OnCallScheduleSlackUsergroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := r.client.GetOnCallScheduleSlackUsergroup(state.ScheduleID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Slack usergroup mapping", err.Error())
		return
	}

	state.UsergroupID = types.StringValue(mapping.UsergroupID)
	state.UsergroupHandle = types.StringValue(mapping.UsergroupHandle)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OnCallScheduleSlackUsergroupResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan OnCallScheduleSlackUsergroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.SetOnCallScheduleSlackUsergroupRequest{
		UsergroupID:     plan.UsergroupID.ValueString(),
		UsergroupHandle: plan.UsergroupHandle.ValueString(),
	}

	mapping, err := r.client.SetOnCallScheduleSlackUsergroup(plan.ScheduleID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Slack usergroup mapping", err.Error())
		return
	}

	plan.ScheduleID = types.StringValue(mapping.ScheduleID)
	plan.UsergroupID = types.StringValue(mapping.UsergroupID)
	plan.UsergroupHandle = types.StringValue(mapping.UsergroupHandle)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OnCallScheduleSlackUsergroupResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state OnCallScheduleSlackUsergroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOnCallScheduleSlackUsergroup(state.ScheduleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Slack usergroup mapping", err.Error())
		return
	}
}
