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

var _ resource.Resource = &DebuggerProbeResource{}
var _ resource.ResourceWithImportState = &DebuggerProbeResource{}

type DebuggerProbeResource struct {
	client *apiclient.Client
}

type DebuggerProbeResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ServiceName types.String `tfsdk:"service_name"`
	FilePath    types.String `tfsdk:"file_path"`
	LineNumber  types.Int64  `tfsdk:"line_number"`
	Expression  types.String `tfsdk:"expression"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

func NewDebuggerProbeResource() resource.Resource {
	return &DebuggerProbeResource{}
}

func (r *DebuggerProbeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_debugger_probe"
}

func (r *DebuggerProbeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Moneat debugger probe.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the debugger probe.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_name": schema.StringAttribute{
				Description: "The name of the service to probe.",
				Required:    true,
			},
			"file_path": schema.StringAttribute{
				Description: "The file path for the probe.",
				Required:    true,
			},
			"line_number": schema.Int64Attribute{
				Description: "The line number for the probe.",
				Required:    true,
			},
			"expression": schema.StringAttribute{
				Description: "The expression to evaluate at the probe point.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the debugger probe is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *DebuggerProbeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DebuggerProbeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DebuggerProbeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateDebuggerProbeRequest{
		ServiceName: plan.ServiceName.ValueString(),
		FilePath:    plan.FilePath.ValueString(),
		LineNumber:  plan.LineNumber.ValueInt64(),
		Expression:  plan.Expression.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
	}

	probe, err := r.client.CreateDebuggerProbe(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating debugger probe", err.Error())
		return
	}

	plan.ID = types.StringValue(probe.ID)
	plan.ServiceName = types.StringValue(probe.ServiceName)
	plan.FilePath = types.StringValue(probe.FilePath)
	plan.LineNumber = types.Int64Value(probe.LineNumber)
	plan.Expression = types.StringValue(probe.Expression)
	plan.Enabled = types.BoolValue(probe.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DebuggerProbeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DebuggerProbeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	probe, err := r.client.GetDebuggerProbe(state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading debugger probe", err.Error())
		return
	}

	state.ServiceName = types.StringValue(probe.ServiceName)
	state.FilePath = types.StringValue(probe.FilePath)
	state.LineNumber = types.Int64Value(probe.LineNumber)
	state.Expression = types.StringValue(probe.Expression)
	state.Enabled = types.BoolValue(probe.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DebuggerProbeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DebuggerProbeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state DebuggerProbeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.UpdateDebuggerProbeRequest{
		ServiceName: plan.ServiceName.ValueString(),
		FilePath:    plan.FilePath.ValueString(),
		LineNumber:  plan.LineNumber.ValueInt64(),
		Expression:  plan.Expression.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
	}

	probe, err := r.client.UpdateDebuggerProbe(state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating debugger probe", err.Error())
		return
	}

	plan.ID = state.ID
	plan.ServiceName = types.StringValue(probe.ServiceName)
	plan.FilePath = types.StringValue(probe.FilePath)
	plan.LineNumber = types.Int64Value(probe.LineNumber)
	plan.Expression = types.StringValue(probe.Expression)
	plan.Enabled = types.BoolValue(probe.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DebuggerProbeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DebuggerProbeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDebuggerProbe(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting debugger probe", err.Error())
		return
	}
}

func (r *DebuggerProbeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	probe, err := r.client.GetDebuggerProbe(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing debugger probe", err.Error())
		return
	}

	state := DebuggerProbeResourceModel{
		ID:          types.StringValue(probe.ID),
		ServiceName: types.StringValue(probe.ServiceName),
		FilePath:    types.StringValue(probe.FilePath),
		LineNumber:  types.Int64Value(probe.LineNumber),
		Expression:  types.StringValue(probe.Expression),
		Enabled:     types.BoolValue(probe.Enabled),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
