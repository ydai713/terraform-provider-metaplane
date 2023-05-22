package provider

import (
	"context"
  "fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/klaviyo/terraform-provider-metaplane/internal/api"
  "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
  "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure the implementation satisfies the expected interfaces.
var (
  _ resource.Resource = &MonitorResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewMonitorResource() resource.Resource {
    return &MonitorResource{}
}

// MonitorResource is the resource implementation.
type MonitorResource struct{
	client *api.Client
}

type MonitorResourceModel struct {
	ConnectionId   types.String    `tfsdk:"connection_id"`
	MonitorId      types.String    `tfsdk:"monitor_id"`
  Type           types.String    `tfsdk:"type"`
  CronTab        types.String    `tfsdk:"cron_tab"`
  AbsolutePath   types.String    `tfsdk:"absolute_path"`
  EntityType     types.String    `tfsdk:"entity_type"`
  UpdatedAt      types.String    `tfsdk:"updated_at"`
  CreatedAt      types.String    `tfsdk:"created_at"`
}

// Metadata returns the resource type name.
func (r *MonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_monitor"
}

// Configure adds the provider configured client to the resource.
func (r *MonitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *MonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Monitor resource",
		Attributes: map[string]schema.Attribute{
			"monitor_id": schema.StringAttribute{
				MarkdownDescription: "Monitor identifier",
				Computed:            true,
        PlanModifiers: []planmodifier.String{
          stringplanmodifier.UseStateForUnknown(),
        },
			},
			"connection_id": schema.StringAttribute{
				MarkdownDescription: "Connection identifier",
				Required:            true,
			},
			"entity_type": schema.StringAttribute{
        MarkdownDescription: "Entity type: table or column",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of monitor, row_count, etc",
				Required:            true,
			},
			"cron_tab": schema.StringAttribute{
				MarkdownDescription: "cron job schedule in * * * * * format",
				Required:            true,
			},
			"absolute_path": schema.StringAttribute{
				MarkdownDescription: "{database}.{schema}.{table}.{column}",
				Required:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "datetime updated",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "datetime created",
				Computed:            true,
        PlanModifiers: []planmodifier.String{
          stringplanmodifier.UseStateForUnknown(),
        },
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
  // Retrieve values from plan
  var plan MonitorResourceModel
  diags := req.Plan.Get(ctx, &plan)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
      return
  }

  // Generate API request body from plan
  newMonitor := api.NewMonitor{
    ConnectionId: plan.ConnectionId.ValueString(),
    Type:         plan.Type.ValueString(),
    EntityType:   plan.EntityType.ValueString(),
    CronTab:      plan.CronTab.ValueString(),
    AbsolutePath: plan.AbsolutePath.ValueString(),
  }

  // Create new monitor
  monitor, err := r.client.CreateMonitor(newMonitor)
  if err != nil {
      resp.Diagnostics.AddError(
          "Error creating monitor",
          "Could not create monitor, unexpected error: "+err.Error(),
      )
      return
  }

  // Map response body to schema and populate Computed attribute values
  plan.MonitorId = types.StringValue(monitor.MonitorId)
  plan.UpdatedAt = types.StringValue(monitor.UpdatedAt)
  plan.CreatedAt = types.StringValue(monitor.CreatedAt)

  // Set state to fully populated data
  diags = resp.State.Set(ctx, plan)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
      return
  }
}

// Read refreshes the Terraform state with the latest data.
func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
  var state MonitorResourceModel
  diags := req.State.Get(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
      return
  }
  // Get refreshed monitor value from API
  connectionId := state.ConnectionId.ValueString()
  monitorId := state.MonitorId.ValueString()

  monitor, err := r.client.GetMonitor(connectionId, monitorId)
  if err != nil {
      resp.Diagnostics.AddError(
          "Error Reading Metaplane Monitor",
          "Could not read Metaplane Monitor ID "+state.MonitorId.ValueString()+": "+err.Error(),
      )
      return
  }
 
  // Overwrite items with refreshed state
  state.Type         = types.StringValue(monitor.Type)
  state.CronTab      = types.StringValue(monitor.CronTab)
  state.AbsolutePath = types.StringValue(monitor.AbsolutePath)
  state.UpdatedAt    = types.StringValue(monitor.UpdatedAt)
  state.CreatedAt    = types.StringValue(monitor.CreatedAt)

  // Set refreshed state
  diags = resp.State.Set(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
      return
  }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
  // Retrieve values from plan
  var plan MonitorResourceModel
  diags := req.Plan.Get(ctx, &plan)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
      return
  }

  // Generate API request body from plan
  updateMonitor := api.UpdateMonitor{
    CronTab:      plan.CronTab.ValueString(),
    MonitorId:    plan.MonitorId.ValueString(),
  }

  // Update existing monitor
  monitor, err := r.client.UpdateMonitor(updateMonitor)
  if err != nil {
      resp.Diagnostics.AddError(
          "Error updating monitor",
          "Could not update monitor, unexpected error: "+err.Error(),
      )
      return
  }

  // Update resource state with updated items and timestamp
  plan.CronTab      = types.StringValue(monitor.CronTab)
  plan.UpdatedAt    = types.StringValue(monitor.UpdatedAt)

  // Set state to fully populated data
  diags = resp.State.Set(ctx, plan)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
      return
  }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
  var state MonitorResourceModel
  diags := req.State.Get(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
      return
  }

  // Delete existing order
  updateMonitor := api.UpdateMonitor {
    MonitorId: state.MonitorId.ValueString(),
    IsEnabled: false,
  }

  _, err := r.client.UpdateMonitor(updateMonitor)

  if err != nil {
      resp.Diagnostics.AddError(
          "Error Deleting Metaplane Monitor",
          "Could not delete monitor, unexpected error: "+err.Error(),
      )
      return
  }
}
