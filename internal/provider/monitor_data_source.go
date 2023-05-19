package provider

import (
	"context"
	"fmt"

  "github.com/klaviyo/terraform-provider-metaplane/internal/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MonitorDataSource{}

func NewMonitorDataSource() datasource.DataSource {
	return &MonitorDataSource{}
}

// MonitorDataSource defines the data source implementation.
type MonitorDataSource struct {
	client *api.Client
}

type MonitorDataSourceModel struct {
	ConnectionId   types.String    `tfsdk:"connection_id"`
	MonitorId      types.String    `tfsdk:"monitor_id"`
  Type           types.String    `tfsdk:"type"`
  CronTab        types.String    `tfsdk:"cron_tab"`
  IsEnabled      types.Bool      `tfsdk:"is_enabled"`
  AbsolutePath   types.String    `tfsdk:"absolute_path"`
  UpdatedAt      types.String    `tfsdk:"updated_at"`
  CreatedAt      types.String    `tfsdk:"created_at"`
}

func (d *MonitorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *MonitorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Monitor data source",
		Attributes: map[string]schema.Attribute{
			"monitor_id": schema.StringAttribute{
				MarkdownDescription: "Monitor identifier",
				Required:            true,
			},
			"connection_id": schema.StringAttribute{
				MarkdownDescription: "Connection identifier",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of monitor, row_count, etc",
				Computed:            true,
			},
			"cron_tab": schema.StringAttribute{
				MarkdownDescription: "cron job schedule in * * * * * format",
				Computed:            true,
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"absolute_path": schema.StringAttribute{
				MarkdownDescription: "{database}.{schema}.{table}.{column}",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "datetime updated",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "datetime created",
				Computed:            true,
			},
		},
	}
}

func (d *MonitorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *MonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state MonitorDataSourceModel

  // Parse the state from the request
  if err := req.Config.Get(ctx, &state); err != nil {
    resp.Diagnostics.AddError("Configuration Error", fmt.Sprintf("Unable to parse configuration: %s", err))
    return
  }

  connectionId := state.ConnectionId.ValueString()
  monitorId := state.MonitorId.ValueString()

  monitor, err := d.client.GetMonitor(connectionId, monitorId)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor, got error: %s", err))
		return
	}

  state.Type         = types.StringValue(monitor.Type)
  state.CronTab      = types.StringValue(monitor.CronTab)
  state.IsEnabled    = types.BoolValue  (monitor.IsEnabled)
  state.AbsolutePath = types.StringValue(monitor.AbsolutePath)
  state.UpdatedAt    = types.StringValue(monitor.UpdatedAt)
  state.CreatedAt    = types.StringValue(monitor.CreatedAt)

  // Set state
  diags := resp.State.Set(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return;
  }
}



