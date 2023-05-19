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
var _ datasource.DataSource = &ConnectionDataSource{}

// ConnectionDataSource defines the data source implementation.
type ConnectionDataSource struct {
	client *api.Client
}

func NewConnectionDataSource() datasource.DataSource {
	return &ConnectionDataSource{}
}

type ConnectionDataSourceModel struct {
	Name           types.String    `tfsdk:"name"`
	ConnectionId   types.String    `tfsdk:"id"`
  Type           types.String    `tfsdk:"type"`
  IsEnabled      types.Bool      `tfsdk:"is_enabled"`
  UpdatedAt      types.String    `tfsdk:"updated_at"`
  CreatedAt      types.String    `tfsdk:"created_at"`
  Status         types.String    `tfsdk:"status"`
}

func (d *ConnectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

func (d *ConnectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Connection data source",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Connection name",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Connection identifier",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of connection, row_count, etc",
				Computed:            true,
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Example identifier",
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
			"status": schema.StringAttribute{
				MarkdownDescription: "Connection status",
				Computed:            true,
			},
		},
	}
}

func (d *ConnectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ConnectionDataSourceModel

  // Parse the state from the request
  if err := req.Config.Get(ctx, &state); err != nil {
    resp.Diagnostics.AddError("Configuration Error", fmt.Sprintf("Unable to parse configuration: %s", err))
    return
  }

  name := state.Name.ValueString()

  connection, err := d.client.GetConnection(name)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read connection, got error: %s", err))
		return
	}

  state.Name         = types.StringValue(connection.Name)
  state.ConnectionId = types.StringValue(connection.ConnectionId)
  state.Type         = types.StringValue(connection.Type)
  state.IsEnabled    = types.BoolValue  (connection.IsEnabled)
  state.UpdatedAt    = types.StringValue(connection.UpdatedAt)
  state.CreatedAt    = types.StringValue(connection.CreatedAt)
  state.Status       = types.StringValue(connection.Status)

  // Set state
  diags := resp.State.Set(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return;
  }
}



