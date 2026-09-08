package provider

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

const responseConfigurationResourceID = "response-configuration"

var _ resource.Resource = &ResponseConfigurationResource{}
var _ resource.ResourceWithImportState = &ResponseConfigurationResource{}

type ResponseConfigurationResource struct {
	client *apiclient.Client
}

type ResponseConfigurationResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ConfigurationJSON types.String `tfsdk:"configuration_json"`
	AllowDestructive  types.Bool   `tfsdk:"allow_destructive"`
	RemoteID          types.String `tfsdk:"remote_id"`
	Revision          types.Int64  `tfsdk:"revision"`
	Checksum          types.String `tfsdk:"checksum"`
	Source            types.String `tfsdk:"source"`
}

func NewResponseConfigurationResource() resource.Resource {
	return &ResponseConfigurationResource{}
}

func (r *ResponseConfigurationResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_response_configuration"
}

func (r *ResponseConfigurationResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages the durable Moneat response configuration package. " +
			"Live incidents, alerts, and telemetry are not Terraform resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Stable Terraform identity for the organization response configuration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"configuration_json": schema.StringAttribute{
				Description: "Canonical response configuration package as JSON. Secrets must use secret:// references.",
				Required:    true,
			},
			"allow_destructive": schema.BoolAttribute{
				Description: "Allow deletions during apply and destroy. Keep false unless the change is intentional.",
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"remote_id": schema.StringAttribute{
				Description: "The current server-side configuration revision identifier.",
				Computed:    true,
			},
			"revision": schema.Int64Attribute{
				Description: "The current server-side configuration revision.",
				Computed:    true,
			},
			"checksum": schema.StringAttribute{
				Description: "The server checksum of the normalized configuration package.",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "The source of the latest server revision.",
				Computed:    true,
			},
		},
	}
}

func (r *ResponseConfigurationResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
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
	r.client = client.ResponseConfigurationClient()
}

func (r *ResponseConfigurationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan ResponseConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.apply(&plan, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ResponseConfigurationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state ResponseConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.client.GetResponseConfigurationState()
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading response configuration", err.Error())
		return
	}
	previousChecksum := state.Checksum.ValueString()
	previousConfiguration := state.ConfigurationJSON
	setResponseConfigurationState(&state, remote)
	if previousChecksum == remote.Checksum && !previousConfiguration.IsNull() && !previousConfiguration.IsUnknown() {
		state.ConfigurationJSON = previousConfiguration
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ResponseConfigurationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan ResponseConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ResponseConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.apply(&plan, intPointer(int(state.Revision.ValueInt64())), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ResponseConfigurationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state ResponseConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !state.AllowDestructive.ValueBool() {
		resp.Diagnostics.AddError(
			"Destructive response configuration deletion is disabled",
			"Set allow_destructive = true before destroying this resource so Terraform can explicitly approve "+
				"deleting the managed configuration.",
		)
		return
	}
	empty, err := json.Marshal(map[string]interface{}{
		"schemaVersion": 1,
		"source":        "TERRAFORM",
		"resources":     []interface{}{},
	})
	if err != nil {
		resp.Diagnostics.AddError("Error encoding empty response configuration", err.Error())
		return
	}
	plan, err := r.client.PlanResponseConfiguration(empty, intPointer(int(state.Revision.ValueInt64())), true)
	if err != nil {
		resp.Diagnostics.AddError("Error planning response configuration deletion", err.Error())
		return
	}
	if !plan.CanApply {
		addConfigurationIssues(&resp.Diagnostics, plan)
		return
	}
	if _, err = r.client.ApplyResponseConfiguration(plan, newIdempotencyKey(), true); err != nil {
		resp.Diagnostics.AddError("Error deleting response configuration", err.Error())
	}
}

func (r *ResponseConfigurationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	remote, err := r.client.GetResponseConfigurationState()
	if err != nil {
		resp.Diagnostics.AddError("Error importing response configuration", err.Error())
		return
	}
	state := ResponseConfigurationResourceModel{}
	setResponseConfigurationState(&state, remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ResponseConfigurationResource) apply(
	plan *ResponseConfigurationResourceModel,
	expectedRevision *int,
	diagnostics interface{ AddError(string, string) },
) {
	configuration, normalizedConfiguration, err := rawMessageFromJSONString(plan.ConfigurationJSON.ValueString())
	if err != nil {
		diagnostics.AddError("Invalid response configuration JSON", err.Error())
		return
	}
	remotePlan, err := r.client.PlanResponseConfiguration(
		configuration,
		expectedRevision,
		plan.AllowDestructive.ValueBool(),
	)
	if err != nil {
		diagnostics.AddError("Error planning response configuration", err.Error())
		return
	}
	if !remotePlan.CanApply {
		addConfigurationIssues(diagnostics, remotePlan)
		return
	}
	state, err := r.client.ApplyResponseConfiguration(
		remotePlan,
		newIdempotencyKey(),
		plan.AllowDestructive.ValueBool(),
	)
	if err != nil {
		diagnostics.AddError("Error applying response configuration", err.Error())
		return
	}
	setResponseConfigurationState(plan, state)
	plan.ConfigurationJSON = types.StringValue(normalizedConfiguration)
}

func setResponseConfigurationState(
	state *ResponseConfigurationResourceModel,
	remote *apiclient.ResponseConfigurationState,
) {
	state.ID = types.StringValue(responseConfigurationResourceID)
	state.RemoteID = types.StringValue(remote.ID)
	state.Revision = types.Int64Value(int64(remote.Revision))
	state.Checksum = types.StringValue(remote.Checksum)
	state.Source = types.StringValue(remote.Source)
	state.ConfigurationJSON = types.StringValue(responseConfigurationStateString(remote.Configuration))
	if state.AllowDestructive.IsNull() || state.AllowDestructive.IsUnknown() {
		state.AllowDestructive = types.BoolValue(false)
	}
}

func addConfigurationIssues(
	diagnostics interface{ AddError(string, string) },
	plan *apiclient.ResponseConfigurationPlan,
) {
	message := "The server rejected the response configuration plan."
	if len(plan.Issues) > 0 {
		message = "The server rejected the response configuration plan: "
		for index, issue := range plan.Issues {
			if index > 0 {
				message += "; "
			}
			message += fmt.Sprintf("%s (%s): %s", issue.Code, issue.Path, issue.Message)
		}
	}
	diagnostics.AddError("Response configuration plan is not applicable", message)
}

func intPointer(value int) *int {
	return &value
}

func newIdempotencyKey() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "terraform-response-configuration"
	}
	return "terraform-response-configuration-" + hex.EncodeToString(bytes)
}
