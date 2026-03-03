package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

var _ resource.Resource = &BusinessHoursResource{}

type BusinessHoursResource struct {
	client *apiclient.Client
}

type BusinessHoursSlotModel struct {
	Day       types.String `tfsdk:"day"`
	StartTime types.String `tfsdk:"start_time"`
	EndTime   types.String `tfsdk:"end_time"`
}

type BusinessHoursResourceModel struct {
	Timezone types.String             `tfsdk:"timezone"`
	Enabled  types.Bool               `tfsdk:"enabled"`
	Windows  []BusinessHoursSlotModel `tfsdk:"windows"`
}

func NewBusinessHoursResource() resource.Resource {
	return &BusinessHoursResource{}
}

func (r *BusinessHoursResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_business_hours"
}

func (r *BusinessHoursResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Moneat business hours configuration. This is a singleton resource (Enterprise only).",
		Attributes: map[string]schema.Attribute{
			"timezone": schema.StringAttribute{
				Description: "The timezone for business hours (e.g., America/New_York).",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether business hours are enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"windows": schema.ListNestedAttribute{
				Description: "Business hours time windows.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"day": schema.StringAttribute{
							Description: "Day of the week (e.g., monday, tuesday).",
							Required:    true,
						},
						"start_time": schema.StringAttribute{
							Description: "Start time (HH:MM format).",
							Required:    true,
						},
						"end_time": schema.StringAttribute{
							Description: "End time (HH:MM format).",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *BusinessHoursResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BusinessHoursResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BusinessHoursResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := r.buildAPIRequest(plan)
	result, err := r.client.UpdateBusinessHours(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error setting business hours", err.Error())
		return
	}

	r.mapToState(&plan, result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BusinessHoursResource) Read(ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse) {
	result, err := r.client.GetBusinessHours()
	if err != nil {
		resp.Diagnostics.AddError("Error reading business hours", err.Error())
		return
	}

	var state BusinessHoursResourceModel
	r.mapToState(&state, result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BusinessHoursResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BusinessHoursResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := r.buildAPIRequest(plan)
	result, err := r.client.UpdateBusinessHours(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating business hours", err.Error())
		return
	}

	r.mapToState(&plan, result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BusinessHoursResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	apiReq := apiclient.UpdateBusinessHoursRequest{
		Timezone: "UTC",
		Enabled:  false,
		Windows:  []apiclient.BusinessHoursSlot{},
	}
	_, err := r.client.UpdateBusinessHours(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error resetting business hours", err.Error())
		return
	}
}

func (r *BusinessHoursResource) buildAPIRequest(
	model BusinessHoursResourceModel,
) apiclient.UpdateBusinessHoursRequest {
	windows := make([]apiclient.BusinessHoursSlot, len(model.Windows))
	for i, w := range model.Windows {
		windows[i] = apiclient.BusinessHoursSlot{
			Day:       w.Day.ValueString(),
			StartTime: w.StartTime.ValueString(),
			EndTime:   w.EndTime.ValueString(),
		}
	}
	return apiclient.UpdateBusinessHoursRequest{
		Timezone: model.Timezone.ValueString(),
		Enabled:  model.Enabled.ValueBool(),
		Windows:  windows,
	}
}

func (r *BusinessHoursResource) mapToState(
	model *BusinessHoursResourceModel, result *apiclient.BusinessHours,
) {
	model.Timezone = types.StringValue(result.Timezone)
	model.Enabled = types.BoolValue(result.Enabled)
	model.Windows = make([]BusinessHoursSlotModel, len(result.Windows))
	for i, w := range result.Windows {
		model.Windows[i] = BusinessHoursSlotModel{
			Day:       types.StringValue(w.Day),
			StartTime: types.StringValue(w.StartTime),
			EndTime:   types.StringValue(w.EndTime),
		}
	}
}
