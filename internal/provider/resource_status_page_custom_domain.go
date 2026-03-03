package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &StatusPageCustomDomainResource{}
var _ resource.ResourceWithImportState = &StatusPageCustomDomainResource{}

type StatusPageCustomDomainResource struct {
	client *apiclient.Client
}

type StatusPageCustomDomainResourceModel struct {
	ID       types.String `tfsdk:"id"`
	PageID   types.String `tfsdk:"page_id"`
	Domain   types.String `tfsdk:"domain"`
	Verified types.Bool   `tfsdk:"verified"`
}

func NewStatusPageCustomDomainResource() resource.Resource {
	return &StatusPageCustomDomainResource{}
}

func (r *StatusPageCustomDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page_custom_domain"
}

func (r *StatusPageCustomDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a custom domain on a Moneat status page.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the custom domain.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"page_id": schema.StringAttribute{
				Description: "The ID of the status page.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				Description: "The custom domain name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"verified": schema.BoolAttribute{
				Description: "Whether the domain has been verified.",
				Computed:    true,
			},
		},
	}
}

func (r *StatusPageCustomDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StatusPageCustomDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StatusPageCustomDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := apiclient.CreateStatusPageCustomDomainRequest{
		Domain: plan.Domain.ValueString(),
	}

	domain, err := r.client.CreateStatusPageCustomDomain(plan.PageID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating status page custom domain", err.Error())
		return
	}

	plan.ID = types.StringValue(domain.ID)
	plan.Domain = types.StringValue(domain.Domain)
	plan.Verified = types.BoolValue(domain.Verified)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageCustomDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StatusPageCustomDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := r.client.GetStatusPageCustomDomain(state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		if apiclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading status page custom domain", err.Error())
		return
	}

	state.Domain = types.StringValue(domain.Domain)
	state.Verified = types.BoolValue(domain.Verified)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StatusPageCustomDomainResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All mutable attributes use RequiresReplace, so Update is never called.
	resp.Diagnostics.AddError(
		"Update not supported",
		"Custom domains cannot be updated in place. Change the domain to trigger replacement.",
	)
}

func (r *StatusPageCustomDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StatusPageCustomDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStatusPageCustomDomain(state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting status page custom domain", err.Error())
		return
	}
}

func (r *StatusPageCustomDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: page_id/domain_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: page_id/domain_id",
		)
		return
	}

	pageID := parts[0]
	domainID := parts[1]

	domain, err := r.client.GetStatusPageCustomDomain(pageID, domainID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing status page custom domain", err.Error())
		return
	}

	state := StatusPageCustomDomainResourceModel{
		ID:       types.StringValue(domain.ID),
		PageID:   types.StringValue(pageID),
		Domain:   types.StringValue(domain.Domain),
		Verified: types.BoolValue(domain.Verified),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
