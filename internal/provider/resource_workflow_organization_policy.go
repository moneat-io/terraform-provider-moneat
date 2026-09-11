package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &WorkflowOrganizationPolicyResource{}
var _ resource.ResourceWithImportState = &WorkflowOrganizationPolicyResource{}

type WorkflowOrganizationPolicyResource struct {
	client *apiclient.Client
}

type WorkflowOrganizationPolicyResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	OrganizationID         types.String `tfsdk:"organization_id"`
	DefaultRole            types.String `tfsdk:"default_role"`
	RequireApproval        types.Bool   `tfsdk:"require_approval"`
	AllowServicePrincipals types.Bool   `tfsdk:"allow_service_principals"`
	MaxConcurrentRuns      types.Int64  `tfsdk:"max_concurrent_runs"`
	AllowedTriggers        types.List   `tfsdk:"allowed_triggers"`
	CreatedAt              types.String `tfsdk:"created_at"`
	UpdatedAt              types.String `tfsdk:"updated_at"`
}

func NewWorkflowOrganizationPolicyResource() resource.Resource {
	return &WorkflowOrganizationPolicyResource{}
}

func (r *WorkflowOrganizationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_organization_policy"
}

func (r *WorkflowOrganizationPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages organization-wide workflow governance defaults.",
		Attributes: map[string]schema.Attribute{
			"id":                       schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"organization_id":          schema.StringAttribute{Computed: true},
			"default_role":             schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("member")},
			"require_approval":         schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(false)},
			"allow_service_principals": schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true)},
			"max_concurrent_runs":      schema.Int64Attribute{Optional: true, Computed: true, Default: int64default.StaticInt64(0)},
			"allowed_triggers": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"created_at": schema.StringAttribute{Computed: true},
			"updated_at": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *WorkflowOrganizationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *apiclient.Client, got: %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *WorkflowOrganizationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkflowOrganizationPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	policy, err := r.upsert(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating workflow organization policy", err.Error())
		return
	}
	mapWorkflowOrganizationPolicyToState(ctx, &resp.Diagnostics, &plan, policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowOrganizationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowOrganizationPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	policy, err := r.client.GetWorkflowOrganizationPolicy()
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading workflow organization policy", err.Error())
		return
	}
	mapWorkflowOrganizationPolicyToState(ctx, &resp.Diagnostics, &state, policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowOrganizationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkflowOrganizationPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	policy, err := r.upsert(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error updating workflow organization policy", err.Error())
		return
	}
	mapWorkflowOrganizationPolicyToState(ctx, &resp.Diagnostics, &plan, policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkflowOrganizationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// The policy endpoint is an upsert-only singleton. Removing Terraform
	// ownership resets it to the server defaults rather than deleting a row.
	defaults := WorkflowOrganizationPolicyResourceModel{
		DefaultRole:            types.StringValue("member"),
		RequireApproval:        types.BoolValue(false),
		AllowServicePrincipals: types.BoolValue(true),
		MaxConcurrentRuns:      types.Int64Value(0),
		AllowedTriggers:        types.ListValueMust(types.StringType, []attr.Value{}),
	}
	if _, err := r.upsert(ctx, defaults); err != nil {
		resp.Diagnostics.AddError("Error resetting workflow organization policy", err.Error())
	}
}

func (r *WorkflowOrganizationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	policy, err := r.client.GetWorkflowOrganizationPolicy()
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow organization policy", err.Error())
		return
	}
	state := WorkflowOrganizationPolicyResourceModel{}
	mapWorkflowOrganizationPolicyToState(ctx, &resp.Diagnostics, &state, policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowOrganizationPolicyResource) upsert(ctx context.Context, model WorkflowOrganizationPolicyResourceModel) (*apiclient.WorkflowOrganizationPolicy, error) {
	var triggers []string
	if !model.AllowedTriggers.IsNull() && !model.AllowedTriggers.IsUnknown() {
		if diags := model.AllowedTriggers.ElementsAs(ctx, &triggers, false); diags.HasError() {
			return nil, fmt.Errorf("allowed_triggers must contain strings")
		}
	}
	max := int(model.MaxConcurrentRuns.ValueInt64())
	var maxPointer *int
	if max > 0 {
		maxPointer = &max
	}
	return r.client.SetWorkflowOrganizationPolicy(apiclient.WorkflowOrganizationPolicyRequest{
		DefaultRole:            model.DefaultRole.ValueString(),
		RequireApproval:        model.RequireApproval.ValueBool(),
		AllowServicePrincipals: model.AllowServicePrincipals.ValueBool(),
		MaxConcurrentRuns:      maxPointer,
		AllowedTriggers:        triggers,
	})
}

func mapWorkflowOrganizationPolicyToState(
	ctx context.Context,
	diags *diag.Diagnostics,
	model *WorkflowOrganizationPolicyResourceModel,
	policy *apiclient.WorkflowOrganizationPolicy,
) {
	model.ID = types.StringValue(policy.ID)
	model.OrganizationID = types.StringValue(policy.OrganizationID)
	model.DefaultRole = types.StringValue(policy.DefaultRole)
	model.RequireApproval = types.BoolValue(policy.RequireApproval)
	model.AllowServicePrincipals = types.BoolValue(policy.AllowServicePrincipals)
	if policy.MaxConcurrentRuns == nil {
		model.MaxConcurrentRuns = types.Int64Value(0)
	} else {
		model.MaxConcurrentRuns = types.Int64Value(int64(*policy.MaxConcurrentRuns))
	}
	model.CreatedAt = optionalString(policy.CreatedAt)
	model.UpdatedAt = optionalString(policy.UpdatedAt)
	triggers, triggerDiags := types.ListValueFrom(ctx, types.StringType, policy.AllowedTriggers)
	diags.Append(triggerDiags...)
	model.AllowedTriggers = triggers
}
